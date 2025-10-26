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
	PolkaKey       string
}

type response struct {
	Error       *string `json:"error,omitempty"`
	CleanedBody *string `json:"cleaned_body,omitempty"`
}

type baseModel struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type userCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userData struct {
	baseModel
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type chirpData struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

type eventData struct {
	UserId string `json:"user_id"`
}

type loginRequest struct {
	userCredentials
}

type loginResponse struct {
	userData
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type refreshTokenResponse struct {
	Token string `json:"token"`
}

type createUserRequest struct {
	userCredentials
}

type updateUserRequest struct {
	userCredentials
}

type createUserResponse struct {
	userData
}

type createChirpRequest struct {
	chirpData
}

type createChirpResponse struct {
	baseModel
	chirpData
}

type webhookRequest struct {
	Event string    `json:"event"`
	Data  eventData `json:"data"`
}
