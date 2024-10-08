// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package model

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// GorpMigration is an object representing the database table.
type GorpMigration struct {
	ID        string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	AppliedAt null.Time `boil:"applied_at" json:"applied_at,omitempty" toml:"applied_at" yaml:"applied_at,omitempty"`

	R *gorpMigrationR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L gorpMigrationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var GorpMigrationColumns = struct {
	ID        string
	AppliedAt string
}{
	ID:        "id",
	AppliedAt: "applied_at",
}

var GorpMigrationTableColumns = struct {
	ID        string
	AppliedAt string
}{
	ID:        "gorp_migrations.id",
	AppliedAt: "gorp_migrations.applied_at",
}

// Generated where

var GorpMigrationWhere = struct {
	ID        whereHelperstring
	AppliedAt whereHelpernull_Time
}{
	ID:        whereHelperstring{field: "`gorp_migrations`.`id`"},
	AppliedAt: whereHelpernull_Time{field: "`gorp_migrations`.`applied_at`"},
}

// GorpMigrationRels is where relationship names are stored.
var GorpMigrationRels = struct {
}{}

// gorpMigrationR is where relationships are stored.
type gorpMigrationR struct {
}

// NewStruct creates a new relationship struct
func (*gorpMigrationR) NewStruct() *gorpMigrationR {
	return &gorpMigrationR{}
}

// gorpMigrationL is where Load methods for each relationship are stored.
type gorpMigrationL struct{}

var (
	gorpMigrationAllColumns            = []string{"id", "applied_at"}
	gorpMigrationColumnsWithoutDefault = []string{"id", "applied_at"}
	gorpMigrationColumnsWithDefault    = []string{}
	gorpMigrationPrimaryKeyColumns     = []string{"id"}
	gorpMigrationGeneratedColumns      = []string{}
)

type (
	// GorpMigrationSlice is an alias for a slice of pointers to GorpMigration.
	// This should almost always be used instead of []GorpMigration.
	GorpMigrationSlice []*GorpMigration

	gorpMigrationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	gorpMigrationType                 = reflect.TypeOf(&GorpMigration{})
	gorpMigrationMapping              = queries.MakeStructMapping(gorpMigrationType)
	gorpMigrationPrimaryKeyMapping, _ = queries.BindMapping(gorpMigrationType, gorpMigrationMapping, gorpMigrationPrimaryKeyColumns)
	gorpMigrationInsertCacheMut       sync.RWMutex
	gorpMigrationInsertCache          = make(map[string]insertCache)
	gorpMigrationUpdateCacheMut       sync.RWMutex
	gorpMigrationUpdateCache          = make(map[string]updateCache)
	gorpMigrationUpsertCacheMut       sync.RWMutex
	gorpMigrationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single gorpMigration record from the query.
func (q gorpMigrationQuery) One(exec boil.Executor) (*GorpMigration, error) {
	o := &GorpMigration{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for gorp_migrations")
	}

	return o, nil
}

// All returns all GorpMigration records from the query.
func (q gorpMigrationQuery) All(exec boil.Executor) (GorpMigrationSlice, error) {
	var o []*GorpMigration

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to GorpMigration slice")
	}

	return o, nil
}

// Count returns the count of all GorpMigration records in the query.
func (q gorpMigrationQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count gorp_migrations rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q gorpMigrationQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if gorp_migrations exists")
	}

	return count > 0, nil
}

// GorpMigrations retrieves all the records using an executor.
func GorpMigrations(mods ...qm.QueryMod) gorpMigrationQuery {
	mods = append(mods, qm.From("`gorp_migrations`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`gorp_migrations`.*"})
	}

	return gorpMigrationQuery{q}
}

// FindGorpMigration retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindGorpMigration(exec boil.Executor, iD string, selectCols ...string) (*GorpMigration, error) {
	gorpMigrationObj := &GorpMigration{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `gorp_migrations` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, gorpMigrationObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from gorp_migrations")
	}

	return gorpMigrationObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *GorpMigration) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no gorp_migrations provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(gorpMigrationColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	gorpMigrationInsertCacheMut.RLock()
	cache, cached := gorpMigrationInsertCache[key]
	gorpMigrationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			gorpMigrationAllColumns,
			gorpMigrationColumnsWithDefault,
			gorpMigrationColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(gorpMigrationType, gorpMigrationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(gorpMigrationType, gorpMigrationMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `gorp_migrations` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `gorp_migrations` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `gorp_migrations` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, gorpMigrationPrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}
	_, err = exec.Exec(cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "model: unable to insert into gorp_migrations")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ID,
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, identifierCols...)
	}
	err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "model: unable to populate default values for gorp_migrations")
	}

CacheNoHooks:
	if !cached {
		gorpMigrationInsertCacheMut.Lock()
		gorpMigrationInsertCache[key] = cache
		gorpMigrationInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the GorpMigration.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *GorpMigration) Update(exec boil.Executor, columns boil.Columns) error {
	var err error
	key := makeCacheKey(columns, nil)
	gorpMigrationUpdateCacheMut.RLock()
	cache, cached := gorpMigrationUpdateCache[key]
	gorpMigrationUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			gorpMigrationAllColumns,
			gorpMigrationPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return errors.New("model: unable to update gorp_migrations, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `gorp_migrations` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, gorpMigrationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(gorpMigrationType, gorpMigrationMapping, append(wl, gorpMigrationPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "model: unable to update gorp_migrations row")
	}

	if !cached {
		gorpMigrationUpdateCacheMut.Lock()
		gorpMigrationUpdateCache[key] = cache
		gorpMigrationUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAll updates all rows with the specified column values.
func (q gorpMigrationQuery) UpdateAll(exec boil.Executor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "model: unable to update all for gorp_migrations")
	}

	return nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o GorpMigrationSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("model: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gorpMigrationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `gorp_migrations` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, gorpMigrationPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to update all in gorpMigration slice")
	}

	return nil
}

var mySQLGorpMigrationUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *GorpMigration) Upsert(exec boil.Executor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model: no gorp_migrations provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(gorpMigrationColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLGorpMigrationUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	gorpMigrationUpsertCacheMut.RLock()
	cache, cached := gorpMigrationUpsertCache[key]
	gorpMigrationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			gorpMigrationAllColumns,
			gorpMigrationColumnsWithDefault,
			gorpMigrationColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			gorpMigrationAllColumns,
			gorpMigrationPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("model: unable to upsert gorp_migrations, could not build update column list")
		}

		ret := strmangle.SetComplement(gorpMigrationAllColumns, strmangle.SetIntersect(insert, update))

		cache.query = buildUpsertQueryMySQL(dialect, "`gorp_migrations`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `gorp_migrations` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(gorpMigrationType, gorpMigrationMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(gorpMigrationType, gorpMigrationMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}
	_, err = exec.Exec(cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "model: unable to upsert for gorp_migrations")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(gorpMigrationType, gorpMigrationMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model: unable to retrieve unique values for gorp_migrations")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, nzUniqueCols...)
	}
	err = exec.QueryRow(cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model: unable to populate default values for gorp_migrations")
	}

CacheNoHooks:
	if !cached {
		gorpMigrationUpsertCacheMut.Lock()
		gorpMigrationUpsertCache[key] = cache
		gorpMigrationUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single GorpMigration record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *GorpMigration) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("model: no GorpMigration provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), gorpMigrationPrimaryKeyMapping)
	sql := "DELETE FROM `gorp_migrations` WHERE `id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete from gorp_migrations")
	}

	return nil
}

// DeleteAll deletes all matching rows.
func (q gorpMigrationQuery) DeleteAll(exec boil.Executor) error {
	if q.Query == nil {
		return errors.New("model: no gorpMigrationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete all from gorp_migrations")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o GorpMigrationSlice) DeleteAll(exec boil.Executor) error {
	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gorpMigrationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `gorp_migrations` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, gorpMigrationPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete all from gorpMigration slice")
	}

	return nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *GorpMigration) Reload(exec boil.Executor) error {
	ret, err := FindGorpMigration(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GorpMigrationSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := GorpMigrationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gorpMigrationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `gorp_migrations`.* FROM `gorp_migrations` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, gorpMigrationPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in GorpMigrationSlice")
	}

	*o = slice

	return nil
}

// GorpMigrationExists checks if the GorpMigration row exists.
func GorpMigrationExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `gorp_migrations` where `id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if gorp_migrations exists")
	}

	return exists, nil
}

// Exists checks if the GorpMigration row exists.
func (o *GorpMigration) Exists(exec boil.Executor) (bool, error) {
	return GorpMigrationExists(exec, o.ID)
}
