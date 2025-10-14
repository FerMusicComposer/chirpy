package handlers

import (
	"sync/atomic"

	"github.com/FerMusicComposer/chirpy/internal/database"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	DbQueries      *database.Queries
	Environment    string
}

type chirp struct {
	Body string `json:"body"`
}

type response struct {
	Error       *string `json:"error,omitempty"`
	CleanedBody *string `json:"cleaned_body,omitempty"`
}

type createUserRequest struct {
	Email string `json:"email"`
}

type createUserResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
}
