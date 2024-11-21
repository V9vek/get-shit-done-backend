package repository

import (
	"context"
	"database/sql"
	"fmt"
	"get-shit-done/model"
	"get-shit-done/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	Db *sql.DB
}

func NewAuthRepository(Db *sql.DB) *AuthRepository {
	return &AuthRepository{Db: Db}
}

func (r *AuthRepository) DoesUserIdExist(context context.Context, userId int) (bool, error) {
	var id int
	QUERY := "SELECT username FROM users WHERE id = $1"

	err := r.Db.QueryRowContext(context, QUERY, userId).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("could not check if userId exist or not: %w", err)
	}

	return id != 0, nil
}

func (r *AuthRepository) DoesUsernameExist(context context.Context, username string) (bool, error) {
	var usernameStr string
	QUERY := "SELECT username FROM users WHERE username = $1"

	err := r.Db.QueryRowContext(context, QUERY, username).Scan(&usernameStr)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("could not check if username exist or not: %w", err)
	}

	return usernameStr != "", nil
}

func (r *AuthRepository) SignUp(context context.Context, cred model.AuthCredentials) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cred.Password), bcrypt.DefaultCost)
	utils.PanicIfError(err)

	QUERY := "INESRT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING id"

	var userId int
	err = r.Db.QueryRowContext(context, QUERY, cred.Email, cred.Username, string(hashedPassword)).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("failed to create the user: %w", err)
	}

	return userId, nil
}

func (r *AuthRepository) SignIn(context context.Context, username string, password string) (int, error) {
	var dbPassword string
	var userId int

	QUERY := "SELECT (id, password) FROM users WHERE username = $1"

	err := r.Db.QueryRowContext(context, QUERY, username).Scan(&userId, &dbPassword)
	if err != nil {
		return 0, fmt.Errorf("failed to get the user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		return 0, fmt.Errorf("username or password is wrong: %w", err)
	}

	return userId, nil
}