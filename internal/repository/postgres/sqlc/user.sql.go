// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package postgres

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (name, username, email, password, role, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, username, email, password, role, status, created_at, updated_at
`

type CreateUserParams struct {
	Name     string
	Username string
	Email    string
	Password string
	Role     string
	Status   string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Name,
		arg.Username,
		arg.Email,
		arg.Password,
		arg.Role,
		arg.Status,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getAllByRole = `-- name: GetAllByRole :many
SELECT id, name, username, email, password, role, status, created_at, updated_at FROM users
WHERE role = $1
`

func (q *Queries) GetAllByRole(ctx context.Context, role string) ([]User, error) {
	rows, err := q.db.Query(ctx, getAllByRole, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.Role,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, username, email, password, role, status, created_at, updated_at FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, username, email, password, role, status, created_at, updated_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET name = $2, username = $3, email = $4, password = $5, role = $6, status = $7, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
`

type UpdateUserParams struct {
	ID       int32
	Name     string
	Username string
	Email    string
	Password string
	Role     string
	Status   string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.ID,
		arg.Name,
		arg.Username,
		arg.Email,
		arg.Password,
		arg.Role,
		arg.Status,
	)
	return err
}
