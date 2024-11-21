package repository

import "database/sql"

type AuthRepository struct {
	Db *sql.DB
}

func NewAuthRepository(Db *sql.DB) *AuthRepository {
	return &AuthRepository{Db: Db}
}
