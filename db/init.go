package db

import (
	"github.com/lbryio/commentron/migration"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	_ "github.com/go-sql-driver/mysql" // import mysql
	"github.com/jmoiron/sqlx"
	_ "github.com/jteeuwen/go-bindata" // so it's detected by `dep ensure`
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
)

var RW boil.Executor
var RO boil.Executor

// Init initializes a database connection based on the dsn provided. It also sets it as the global db connection.
func Init(dsnRO, dsnRW string, debug bool) error {
	dsnSuffix := "?parseTime=1&collation=utf8mb4_unicode_ci"
	dsnRO += dsnSuffix
	dbConnRO, err := sqlx.Connect("mysql", dsnRO)
	if err != nil {
		return errors.Err(err)
	}
	dbConnRO.SetMaxOpenConns(300)
	err = dbConnRO.Ping()
	if err != nil {
		return errors.Err(err)
	}

	dsnRW += dsnSuffix
	dbConnRW, err := sqlx.Connect("mysql", dsnRW)
	if err != nil {
		return errors.Err(err)
	}
	dbConnRW.SetMaxOpenConns(300)
	err = dbConnRW.Ping()
	if err != nil {
		return errors.Err(err)
	}

	logWrapperRO := &QueryLogger{DB: dbConnRO, Name: "Commentron-RO"}
	if debug {
		logWrapperRO.Logger = log.StandardLogger()
	}

	logWrapperRW := &QueryLogger{DB: dbConnRW, Name: "Commentron-RW"}
	if debug {
		logWrapperRW.Logger = log.StandardLogger()
	}

	RO = logWrapperRO
	RW = logWrapperRW

	migrations := &migrate.AssetMigrationSource{
		Asset:    migration.Asset,
		AssetDir: migration.AssetDir,
		Dir:      "migration",
	}
	n, migrationErr := migrate.Exec(dbConnRW.DB, "mysql", migrations, migrate.Up)
	if migrationErr != nil {
		return errors.Err(migrationErr)
	}
	log.Printf("Applied %d migrations", n)

	return nil
}
