package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type RoleRepositoryImpl struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewRoleRepository(db *pgxpool.Pool, logger *zap.Logger) *RoleRepositoryImpl {
	logger.Info("create")
	return &RoleRepositoryImpl{db: db, logger: logger}
}

type Role struct {
	ID   int    `db:"id" json:"ID"`
	Name string `db:"name" json:"name"`
}

type RoleRepository interface {
	GetRoleById(id int) (Role, error)
}

func (r RoleRepositoryImpl) GetRoleById(id int) (Role, error) {
	query := "SELECT id, name FROM roles WHERE id = $1"

	row := r.db.QueryRow(context.Background(), query, id)
	var role Role

	err := row.Scan(&role.ID, &role.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return Role{}, fmt.Errorf("role with id %d not found", id)
		}
		return Role{}, fmt.Errorf("error retrieving role: %w", err)
	}
	return role, nil
}
