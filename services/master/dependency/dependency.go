package dependency

import (
	"database/sql"
	"fmt"
	"github.com/example-golang-projects/clean-architecture/cmd/server/master/config"
	"github.com/example-golang-projects/clean-architecture/golibs/database"
	"github.com/example-golang-projects/clean-architecture/golibs/database/migration"
	"log"
)

type MasterDependency struct {
	Config config.Config
	DB     *sql.DB

	// List of client/third-party
}

func InitMasterDependency(cfg config.Config) MasterDependency {
	// Init database
	// ...
	// Init client/third-party
	// ...
	connStr := fmt.Sprintf("dbname=%v user=%v password=%v host=%v port=%v sslmode=%v", cfg.Database.MasterDB.Database, cfg.Database.MasterDB.Username, cfg.Database.MasterDB.Password, cfg.Database.MasterDB.Host, cfg.Database.MasterDB.Port, cfg.Database.MasterDB.SSLMode)

	db, err := database.NewDatabase(connStr)
	if err != nil {
		log.Panic(err)
	}

	err = migration.Up(db, "./services/master/migrations", cfg.ServiceName)
	if err != nil {
		log.Panic(err)
	}
	return MasterDependency{
		Config: cfg,
		DB:     db,
	}
}
