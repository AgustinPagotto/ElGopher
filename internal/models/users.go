package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword string
	Created        time.Time
}

type UserModelInterface interface {
	Insert(ctx context.Context, name, email, password string) error
	Authenticate(ctx context.Context, email, password string) (int, error)
	Exists(ctx context.Context, id int) (bool, error)
	Get(ctx context.Context, id int) (*User, error)
	PasswordUpdate(ctx context.Context, id int, currentPassword, newPassword string) error
}

type UserModel struct {
	POOL *pgxpool.Pool
}

func (um *UserModel) Insert(ctx context.Context, name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	sqlQuery := `INSERT INTO users (name, email, hashed_password, created) 
		VALUES ($1,$2,$3, UTC_TIMESTAMP())`
	_, err = um.POOL.Exec(ctx, sqlQuery, name, email, string(hashedPassword))

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" && pgErr.ConstraintName == "users_uc_email" {
			return ErrDuplicateEmail
		}
	}
	return nil
}

func (um *UserModel) Authenticate(ctx context.Context, email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	sqlQuery := "SELECT id, hashed_password FROM users WHERE email = $1;"
	err := um.POOL.QueryRow(ctx, sqlQuery, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}
func (um *UserModel) Exists(ctx context.Context, id int) (bool, error) {
	var exists bool
	sqlQuery := "SELECT EXISTS(SELECT true FROM users WHERE id = $1);"
	err := um.POOL.QueryRow(ctx, sqlQuery, id).Scan(&exists)
	return exists, err
}
func (um *UserModel) Get(ctx context.Context, id int) (*User, error) {
	var user User
	sqlQuery := "SELECT id, name, email, created FROM users WHERE id = $1;"
	err := um.POOL.QueryRow(ctx, sqlQuery, id).Scan(&user.ID, &user.Name, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		} else {
			return nil, err
		}
	}
	return &user, nil
}
func (um *UserModel) PasswordUpdate(ctx context.Context, id int, currentPassword, newPassword string) error {
	var hashedPassword []byte
	sqlQuery := "SELECT hashed_password FROM users WHERE id = $1;"
	err := um.POOL.QueryRow(ctx, sqlQuery, id).Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	sqlQuery = "UPDATE users SET hashed_password = $1 WHERE id = $2;"
	_, err = um.POOL.Exec(ctx, sqlQuery, string(newHashedPassword), id)
	return err
}
