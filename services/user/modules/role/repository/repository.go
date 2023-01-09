package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/example-golang-projects/clean-architecture/services/user/entities"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) CreateRole(ctx context.Context, role entities.Role) error {
	fmt.Println("RoleRepository")
	return nil
}
