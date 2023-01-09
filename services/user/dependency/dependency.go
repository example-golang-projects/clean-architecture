package dependency

import (
	"database/sql"
	"fmt"
	"github.com/example-golang-projects/clean-architecture/cmd/server/user/config"
	"github.com/example-golang-projects/clean-architecture/packages/database"
	"github.com/example-golang-projects/clean-architecture/packages/database/migration"
	"log"
)

type UserDependency struct {
	Config config.Config
	DB     *sql.DB

	// List of client/third-party
}

func InitUserDependency(cfg config.Config) UserDependency {
	connStr := fmt.Sprintf("dbname=%v user=%v password=%v host=%v port=%v sslmode=%v", cfg.Database.UserDB.Database, cfg.Database.UserDB.Username, cfg.Database.UserDB.Password, cfg.Database.UserDB.Host, cfg.Database.UserDB.Port, cfg.Database.UserDB.SSLMode)

	db, err := database.NewDatabase(connStr)
	if err != nil {
		log.Panic(err)
	}

	err = migration.Up(db, "./services/user/migrations", cfg.ServiceName)
	if err != nil {
		log.Panic(err)
	}

	return UserDependency{
		Config: cfg,
		DB:     db,
	}
}
