package migration

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	schema_migration = "schema_migration"
)

type Migration struct {
	FileName      string
	Version       int
	SqlQuery      string
	MigrationPath string // path to .sql script or go file
}

type Migrations []Migration

func (m *Migration) Noti() {
	fmt.Println(fmt.Sprintf("%s Migrate database file %s successfully!", time.Now().Format("2006-01-02 03:04:05"), m.FileName))
}
func (ms Migrations) Len() int {
	return len(ms)
}

func (ms Migrations) Less(i, j int) bool {
	return ms[i].Version < ms[j].Version
}

func (ms Migrations) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

type Migrator struct {
	db         *sql.DB
	Migrations map[string]*Migration
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{
		db: db,
	}
}

func (m *Migrator) initSchemaMigration() error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
					id serial NOT NULL,
					service VARCHAR NOT NULL,
					version bigint NOT NULL,
					tstamp timestamp NULL default now(),
					PRIMARY KEY(id),
					UNIQUE(service, version)
		);

`, schema_migration)
	_, err := m.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (m *Migrator) insertVersion(tx *sql.Tx, version int, serviceName string) error {
	query := fmt.Sprintf(`
		insert into %s (version,service) 
		values (%d,'%v');
	`, schema_migration, version, serviceName)
	_, err := tx.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (m *Migrator) getLatestVersionByServiceName(service string) (int, error) {
	currentVersion := 0
	query := fmt.Sprintf(`
		SELECT sm.version 
		FROM %s sm 
		WHERE sm.service = '%s' 
		ORDER BY sm.id DESC 
		LIMIT 1;
	`, schema_migration, service)
	err := m.db.QueryRow(query).Scan(&currentVersion)
	if err != nil {
		return 0, err
	}
	return currentVersion, nil
}

func Up(db *sql.DB, dir, serviceName string) (err error) {
	m := NewMigrator(db)
	if err = m.initSchemaMigration(); err != nil {
		return err
	}

	latestVersion, err := m.getLatestVersionByServiceName(serviceName)
	if err != nil && !strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
		return err
	}

	migrations, err := collectNewMigrations(dir, latestVersion)
	if err != nil {
		return err
	}
	fmt.Println(migrations)
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	for _, migration := range migrations {
		_, err = tx.Exec(migration.SqlQuery)
		if err != nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("error when execute migration file %s: %v", migration.MigrationPath, err.Error()))
		}
		err = m.insertVersion(tx, migration.Version, serviceName)
		if err != nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("error when insert version of migration file %s: %v", migration.MigrationPath, err.Error()))
		}
		migration.Noti()
	}
	tx.Commit()
	return nil
}

func collectNewMigrations(dirPath string, currentVer int) (migrations []Migration, err error) {
	if _, err := os.Stat(filepath.FromSlash(dirPath)); errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%s directory does not exist", dirPath)
	}

	mapVersion := make(map[int]bool)
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".sql" {
			// Ex: 00001_create_order.sql => ['00001', 'create', 'order']
			fileNameDetails := strings.Split(info.Name(), "_")

			migrationFileBytes, err := os.ReadFile(path)
			if err != nil {
				return errors.New(fmt.Sprintf("error when read migration file %v", path))
			}
			version, err := strconv.Atoi(fileNameDetails[0])
			if err != nil {
				return errors.New(fmt.Sprintf("error when version format of migration file %v is invalid", path))
			}
			if _, exists := mapVersion[version]; exists {
				return errors.New(fmt.Sprintf("error when duplicate migration file %v", path))
			}
			mapVersion[version] = true
			if version > currentVer {
				migrations = append(migrations, Migration{
					FileName:      info.Name(),
					Version:       version,
					SqlQuery:      string(migrationFileBytes),
					MigrationPath: path,
				})
			}
		}
		return nil
	})
	if err != nil {
		return
	}
	sort.Sort(Migrations(migrations))
	return
}
