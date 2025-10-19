package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TestPasswordHashing verifies that hashing a password and then checking it works correctly.
// It also checks that an incorrect password does not match the hash.
func TestPasswordHashing(t *testing.T) {
	password := "my-super-secret-password"

	// 1. Hash the password
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() returned an unexpected error: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword() returned an empty hash")
	}

	// 2. Check the correct password against the hash
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() returned an unexpected error for correct password: %v", err)
	}
	if !match {
		t.Error("expected password to match hash, but it didn't")
	}

	// 3. Check an incorrect password against the hash
	wrongPassword := "not-the-right-password"
	match, err = CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() returned an unexpected error for incorrect password: %v", err)
	}
	if match {
		t.Error("expected incorrect password to not match hash, but it did")
	}
}

// TestJWTFlow tests the complete lifecycle of creating and validating a JSON Web Token.
func TestJWTFlow(t *testing.T) {
	tokenSecret := "a-very-secure-secret-key"
	userId := uuid.New()
	expiresIn := time.Hour * 1 // Token is valid for 1 hour

	// 1. Create a new JWT
	tokenString, err := MakeJWT(userId, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() returned an unexpected error: %v", err)
	}
	if tokenString == "" {
		t.Fatal("MakeJWT() returned an empty token string")
	}

	// 2. Validate the JWT with the correct secret
	validatedUserId, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT() returned an unexpected error for a valid token: %v", err)
	}
	if validatedUserId != userId {
		t.Errorf("expected user ID %v, but got %v", userId, validatedUserId)
	}
}

// TestInvalidJWT tests that validation fails for tokens that are malformed,
// use the wrong secret, or are otherwise invalid.
func TestInvalidJWT(t *testing.T) {
	tokenSecret := "a-very-secure-secret-key"
	wrongSecret := "this-is-not-the-right-key"
	userId := uuid.New()

	// Create a valid token to tamper with
	validToken, err := MakeJWT(userId, tokenSecret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create a token for testing: %v", err)
	}

	testCases := []struct {
		name        string
		token       string
		secret      string
		expectError bool
	}{
		{
			name:        "Malformed Token",
			token:       "this.is.not.a.jwt",
			secret:      tokenSecret,
			expectError: true,
		},
		{
			name:        "Token with Wrong Secret",
			token:       validToken,
			secret:      wrongSecret,
			expectError: true,
		},
		{
			name:        "Empty Token",
			token:       "",
			secret:      tokenSecret,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ValidateJWT(tc.token, tc.secret)
			if tc.expectError && err == nil {
				t.Error("expected an error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("did not expect an error but got: %v", err)
			}
		})
	}
}

// TestExpiredJWT verifies that an expired token fails validation with the correct error.
func TestExpiredJWT(t *testing.T) {
	tokenSecret := "a-very-secure-secret-key"
	userId := uuid.New()
	
	// Create a token that expired 1 hour ago
	expiredToken, err := MakeJWT(userId, tokenSecret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() failed to create an expired token: %v", err)
	}

	_, err = ValidateJWT(expiredToken, tokenSecret)
	if err == nil {
		t.Fatal("expected an error for an expired token, but got none")
	}

	// Check if the error is specifically due to token expiration
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Errorf("expected jwt.ErrTokenExpired, but got a different error: %v", err)
	}
}
