package dependency

import (
	"github.com/example-golang-projects/clean-architecture/cmd/server/user/config"
	"github.com/jackc/pgx/v5"
)

type UserDependency struct {
	Config config.Config
	DB     *pgx.Conn

	// List of client/third-party
}

func InitUserDependency(cfg config.Config) UserDependency {
	// Init database
	// ...
	// Init client/third-party
	// ...

	return UserDependency{
		Config: cfg,
		DB:     nil,
	}
}
