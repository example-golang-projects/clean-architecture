package database

import (
	"context"
	"fmt"
	"github.com/example-golang-projects/clean-architecture/core/config"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

type Database struct {
	Conn *pgx.Conn
}

func NewDatabase(d config.Config) *Database {
	c := d.Databases.PostgresDB
	connString := fmt.Sprintf("dbname=%v user=%v password=%v host=%v port=%v sslmode=%v", c.Database, c.Username, c.Password, c.Host, c.Port, c.SSLMode)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Panic(err, nil, nil)
	}
	return &Database{Conn: conn}
}

func newLogger() logger.Interface {
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	return logger
}
