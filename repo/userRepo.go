package repo

import (
	"context"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}
type User struct {
	ID       int
	Login    string
	Password string
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(login, HashPassword string, ctx context.Context) error {
	query := "INSERT INTO users (login, password) VALUES ($1, $2)"
	_, err := r.db.ExecContext(ctx, query, login, HashPassword)
	return err
}

func (r *UserRepository) GetUserByLogin(login string, ctx context.Context) (User, error) {

	var user User
	query := "SELECT id, login, password FROM users WHERE login = $1"
	err := r.db.QueryRowContext(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
