// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"
)

type Chirp struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"user_id"`
	Body      string    `json:"body"`
}

type RefreshToken struct {
	Token     string       `json:"token"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	ExpiresAt time.Time    `json:"expires_at"`
	RevokedAt sql.NullTime `json:"revoked_at"`
	UserID    string       `json:"user_id"`
}

type User struct {
	ID             string       `json:"id"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Email          string       `json:"email"`
	HashedPassword string       `json:"hashed_password"`
	IsChirpyRed    sql.NullBool `json:"is_chirpy_red"`
}
