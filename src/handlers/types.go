package handlers

import (
	"sync/atomic"

	"github.com/FerMusicComposer/chirpy/internal/database"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	DbQueries      *database.Queries
	Environment    string
	JWTSecret      string
}


type response struct {
	Error       *string `json:"error,omitempty"`
	CleanedBody *string `json:"cleaned_body,omitempty"`
}

type loginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	ExpiresInSeconds int `json:"expires_in_seconds,omitempty"`
}

type loginResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
	Token     string `json:"token"`
}

type createUserRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}


type createUserResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
}

type createChirpRequest struct {
	Body string `json:"body"`
	UserID string `json:"user_id"`
}

type createChirpResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserID    string `json:"user_id"`
}
