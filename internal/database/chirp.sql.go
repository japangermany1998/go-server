// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: chirp.sql

package database

import (
	"context"
)

const createChirp = `-- name: CreateChirp :one
INSERT INTO chirp (id, created_at, updated_at, body, user_id)
VALUES (
           gen_random_uuid(), now(), now(), $1, $2
    )
RETURNING id, created_at, updated_at, user_id, body
`

type CreateChirpParams struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

func (q *Queries) CreateChirp(ctx context.Context, arg CreateChirpParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, createChirp, arg.Body, arg.UserID)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Body,
	)
	return i, err
}

const deleteChirp = `-- name: DeleteChirp :one
DELETE FROM chirp
WHERE user_id = $1 AND id = $2
RETURNING id, created_at, updated_at, user_id, body
`

type DeleteChirpParams struct {
	UserID string `json:"user_id"`
	ID     string `json:"id"`
}

func (q *Queries) DeleteChirp(ctx context.Context, arg DeleteChirpParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, deleteChirp, arg.UserID, arg.ID)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Body,
	)
	return i, err
}

const getAllChirps = `-- name: GetAllChirps :many
SELECT id, created_at, updated_at, user_id, body FROM chirp
WHERE user_id = $1 OR $2
ORDER BY created_at
`

type GetAllChirpsParams struct {
	UserID  string      `json:"user_id"`
	Column2 interface{} `json:"column_2"`
}

func (q *Queries) GetAllChirps(ctx context.Context, arg GetAllChirpsParams) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getAllChirps, arg.UserID, arg.Column2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.Body,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getChirp = `-- name: GetChirp :one
SELECT id, created_at, updated_at, user_id, body FROM chirp
WHERE id = $1
`

func (q *Queries) GetChirp(ctx context.Context, id string) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, getChirp, id)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Body,
	)
	return i, err
}
