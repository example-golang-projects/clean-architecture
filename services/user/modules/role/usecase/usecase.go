package usecase

import (
	"context"
	"database/sql"
	"github.com/example-golang-projects/clean-architecture/services/user/entities"
	"github.com/example-golang-projects/clean-architecture/services/user/modules/role/repository"
)

type RoleRepoForRoleUseCase interface {
	CreateRole(context.Context, entities.Role) error
}

type RoleUseCase struct {
	db          *sql.DB
	roleUseCase RoleRepoForRoleUseCase
}

func NewRoleUseCase(db *sql.DB) RoleUseCase {
	return RoleUseCase{
		db:          db,
		roleUseCase: repository.NewRoleRepository(db),
	}
}

func (u *RoleUseCase) CreateRole(ctx context.Context, args int) (err error) {
	return u.roleUseCase.CreateRole(ctx, entities.Role{})
}
