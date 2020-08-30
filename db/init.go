package db

import (
	"github.com/lbryio/commentron/migration"
	"github.com/lbryio/lbry.go/extras/errors"

	_ "github.com/go-sql-driver/mysql" // import mysql
	"github.com/jmoiron/sqlx"
	_ "github.com/jteeuwen/go-bindata" // so it's detected by `dep ensure`
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
)

// Init initializes a database connection based on the dsn provided. It also sets it as the global db connection.
func Init(dsn string, debug bool) error {
	dsn += "?parseTime=1&collation=utf8mb4_unicode_ci"
	dbConn, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return errors.Err(err)
	}

	err = dbConn.Ping()
	if err != nil {
		return errors.Err(err)
	}

	if debug {
		boil.DebugMode = true
	}

	boil.SetDB(dbConn)

	migrations := &migrate.AssetMigrationSource{
		Asset:    migration.Asset,
		AssetDir: migration.AssetDir,
		Dir:      "migration",
	}
	n, migrationErr := migrate.Exec(dbConn.DB, "mysql", migrations, migrate.Up)
	if migrationErr != nil {
		return errors.Err(migrationErr)
	}
	log.Printf("Applied %d migrations", n)

	return nil
}
