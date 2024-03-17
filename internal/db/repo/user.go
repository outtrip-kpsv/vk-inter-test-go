package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepositoryImpl struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewUserRepository(db *pgxpool.Pool, logger *zap.Logger) *UserRepositoryImpl {
	logger.Info("create")
	return &UserRepositoryImpl{db: db, logger: logger}
}

type User struct {
	ID     int    `db:"id" json:"ID"`
	Login  string `db:"login" json:"login"`
	Pass   string `db:"pass" json:"pass"`
	RoleID int    `db:"role_id" json:"-"`
}

type UserRepository interface {
	GetUserByLogin(login string) (User, error)
	CreateUser(user User) error
}

func (u UserRepositoryImpl) GetUserByLogin(login string) (User, error) {
	var user User
	sql := "SELECT * FROM users WHERE login = $1 LIMIT 1"
	err := u.db.QueryRow(context.Background(), sql, login).Scan(&user.ID, &user.Login, &user.Pass, &user.RoleID)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u UserRepositoryImpl) CreateUser(user User) error {
	sql := "INSERT INTO users (login, pass) VALUES ($1, $2) RETURNING id"
	err := u.db.QueryRow(context.Background(), sql, user.Login, user.Pass).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}
