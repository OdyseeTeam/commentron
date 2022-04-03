package db

import (
	"github.com/lbryio/commentron/migration"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	_ "github.com/go-sql-driver/mysql" // import mysql
	"github.com/jmoiron/sqlx"
	_ "github.com/kevinburke/go-bindata" // so it's detected by `dep ensure`
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
)

// RW this db is used for read-write calls, it can be used for RO calls too but to load balance use the RO please.
var RW boil.Executor

// RO this db can only be used for read-only calls. Calls made that trigger a RW will break replication if successful
// but also should actually produce an error from mysql.
var RO boil.Executor

// Init initializes a database connection based on the dsn provided. It also sets it as the global db connection.
func Init(dsnRO, dsnRW string, debug bool) error {
	dsnSuffix := "?parseTime=1&collation=utf8mb4_unicode_ci"
	dsnRO += dsnSuffix
	dbConnRO, err := sqlx.Connect("mysql", dsnRO)
	if err != nil {
		return errors.Err(err)
	}
	dbConnRO.SetMaxOpenConns(500)
	err = dbConnRO.Ping()
	if err != nil {
		return errors.Err(err)
	}

	dsnRW += dsnSuffix
	dbConnRW, err := sqlx.Connect("mysql", dsnRW)
	if err != nil {
		return errors.Err(err)
	}
	dbConnRW.SetMaxOpenConns(500)
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
