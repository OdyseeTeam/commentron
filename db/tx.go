package db

import (
	"github.com/lbryio/lbry.go/v2/extras/errors"

	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// TxFunc is a function that can be wrapped in a transaction
type TxFunc func(tx boil.Transactor) error

// WithTx wraps a function in an sql transaction. the transaction is committed if there's no error, or rolled back if there is one.
// if `currentTx` is nil, a new transaction is started. otherwise currentTx is used as the transaction
func WithTx(exec boil.Executor, currentTx boil.Transactor, f TxFunc) (err error) {
	var tx boil.Transactor

	if currentTx != nil {
		tx = currentTx
	} else {
		creator, ok := exec.(*QueryLogger)
		if !ok {
			return errors.Err("database does not support transactions")
		}

		tx, err = creator.Begin()
		if err != nil {
			return errors.Err(err)
		}

		defer func() {
			if p := recover(); p != nil {
				rollbackErr := tx.Rollback()
				if rollbackErr != nil {
					log.Errorln(errors.Prefix("rollback failed after tx error", rollbackErr))
				}
				panic(p)
			} else if err != nil {
				rollbackErr := tx.Rollback()
				if rollbackErr != nil {
					log.Errorln(errors.Prefix("rollback failed after tx error", rollbackErr))
				}
			} else {
				err = errors.Err(tx.Commit())
			}
		}()
	}

	return f(tx)
}
