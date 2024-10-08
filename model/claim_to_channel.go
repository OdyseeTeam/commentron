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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// ClaimToChannel is an object representing the database table.
type ClaimToChannel struct {
	ClaimID   string    `boil:"claim_id" json:"claim_id" toml:"claim_id" yaml:"claim_id"`
	ChannelID string    `boil:"channel_id" json:"channel_id" toml:"channel_id" yaml:"channel_id"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *claimToChannelR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L claimToChannelL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ClaimToChannelColumns = struct {
	ClaimID   string
	ChannelID string
	CreatedAt string
	UpdatedAt string
}{
	ClaimID:   "claim_id",
	ChannelID: "channel_id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

var ClaimToChannelTableColumns = struct {
	ClaimID   string
	ChannelID string
	CreatedAt string
	UpdatedAt string
}{
	ClaimID:   "claim_to_channel.claim_id",
	ChannelID: "claim_to_channel.channel_id",
	CreatedAt: "claim_to_channel.created_at",
	UpdatedAt: "claim_to_channel.updated_at",
}

// Generated where

var ClaimToChannelWhere = struct {
	ClaimID   whereHelperstring
	ChannelID whereHelperstring
	CreatedAt whereHelpertime_Time
	UpdatedAt whereHelpertime_Time
}{
	ClaimID:   whereHelperstring{field: "`claim_to_channel`.`claim_id`"},
	ChannelID: whereHelperstring{field: "`claim_to_channel`.`channel_id`"},
	CreatedAt: whereHelpertime_Time{field: "`claim_to_channel`.`created_at`"},
	UpdatedAt: whereHelpertime_Time{field: "`claim_to_channel`.`updated_at`"},
}

// ClaimToChannelRels is where relationship names are stored.
var ClaimToChannelRels = struct {
}{}

// claimToChannelR is where relationships are stored.
type claimToChannelR struct {
}

// NewStruct creates a new relationship struct
func (*claimToChannelR) NewStruct() *claimToChannelR {
	return &claimToChannelR{}
}

// claimToChannelL is where Load methods for each relationship are stored.
type claimToChannelL struct{}

var (
	claimToChannelAllColumns            = []string{"claim_id", "channel_id", "created_at", "updated_at"}
	claimToChannelColumnsWithoutDefault = []string{"claim_id", "channel_id", "created_at", "updated_at"}
	claimToChannelColumnsWithDefault    = []string{}
	claimToChannelPrimaryKeyColumns     = []string{"claim_id"}
	claimToChannelGeneratedColumns      = []string{}
)

type (
	// ClaimToChannelSlice is an alias for a slice of pointers to ClaimToChannel.
	// This should almost always be used instead of []ClaimToChannel.
	ClaimToChannelSlice []*ClaimToChannel

	claimToChannelQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	claimToChannelType                 = reflect.TypeOf(&ClaimToChannel{})
	claimToChannelMapping              = queries.MakeStructMapping(claimToChannelType)
	claimToChannelPrimaryKeyMapping, _ = queries.BindMapping(claimToChannelType, claimToChannelMapping, claimToChannelPrimaryKeyColumns)
	claimToChannelInsertCacheMut       sync.RWMutex
	claimToChannelInsertCache          = make(map[string]insertCache)
	claimToChannelUpdateCacheMut       sync.RWMutex
	claimToChannelUpdateCache          = make(map[string]updateCache)
	claimToChannelUpsertCacheMut       sync.RWMutex
	claimToChannelUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single claimToChannel record from the query.
func (q claimToChannelQuery) One(exec boil.Executor) (*ClaimToChannel, error) {
	o := &ClaimToChannel{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for claim_to_channel")
	}

	return o, nil
}

// All returns all ClaimToChannel records from the query.
func (q claimToChannelQuery) All(exec boil.Executor) (ClaimToChannelSlice, error) {
	var o []*ClaimToChannel

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to ClaimToChannel slice")
	}

	return o, nil
}

// Count returns the count of all ClaimToChannel records in the query.
func (q claimToChannelQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count claim_to_channel rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q claimToChannelQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if claim_to_channel exists")
	}

	return count > 0, nil
}

// ClaimToChannels retrieves all the records using an executor.
func ClaimToChannels(mods ...qm.QueryMod) claimToChannelQuery {
	mods = append(mods, qm.From("`claim_to_channel`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`claim_to_channel`.*"})
	}

	return claimToChannelQuery{q}
}

// FindClaimToChannel retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindClaimToChannel(exec boil.Executor, claimID string, selectCols ...string) (*ClaimToChannel, error) {
	claimToChannelObj := &ClaimToChannel{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `claim_to_channel` where `claim_id`=?", sel,
	)

	q := queries.Raw(query, claimID)

	err := q.Bind(nil, exec, claimToChannelObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from claim_to_channel")
	}

	return claimToChannelObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ClaimToChannel) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no claim_to_channel provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(claimToChannelColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	claimToChannelInsertCacheMut.RLock()
	cache, cached := claimToChannelInsertCache[key]
	claimToChannelInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			claimToChannelAllColumns,
			claimToChannelColumnsWithDefault,
			claimToChannelColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(claimToChannelType, claimToChannelMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(claimToChannelType, claimToChannelMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `claim_to_channel` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `claim_to_channel` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `claim_to_channel` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, claimToChannelPrimaryKeyColumns))
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
		return errors.Wrap(err, "model: unable to insert into claim_to_channel")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ClaimID,
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, identifierCols...)
	}
	err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "model: unable to populate default values for claim_to_channel")
	}

CacheNoHooks:
	if !cached {
		claimToChannelInsertCacheMut.Lock()
		claimToChannelInsertCache[key] = cache
		claimToChannelInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the ClaimToChannel.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ClaimToChannel) Update(exec boil.Executor, columns boil.Columns) error {
	var err error
	key := makeCacheKey(columns, nil)
	claimToChannelUpdateCacheMut.RLock()
	cache, cached := claimToChannelUpdateCache[key]
	claimToChannelUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			claimToChannelAllColumns,
			claimToChannelPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return errors.New("model: unable to update claim_to_channel, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `claim_to_channel` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, claimToChannelPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(claimToChannelType, claimToChannelMapping, append(wl, claimToChannelPrimaryKeyColumns...))
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
		return errors.Wrap(err, "model: unable to update claim_to_channel row")
	}

	if !cached {
		claimToChannelUpdateCacheMut.Lock()
		claimToChannelUpdateCache[key] = cache
		claimToChannelUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAll updates all rows with the specified column values.
func (q claimToChannelQuery) UpdateAll(exec boil.Executor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "model: unable to update all for claim_to_channel")
	}

	return nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ClaimToChannelSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), claimToChannelPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `claim_to_channel` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, claimToChannelPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to update all in claimToChannel slice")
	}

	return nil
}

var mySQLClaimToChannelUniqueColumns = []string{
	"claim_id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ClaimToChannel) Upsert(exec boil.Executor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model: no claim_to_channel provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(claimToChannelColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLClaimToChannelUniqueColumns, o)

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

	claimToChannelUpsertCacheMut.RLock()
	cache, cached := claimToChannelUpsertCache[key]
	claimToChannelUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			claimToChannelAllColumns,
			claimToChannelColumnsWithDefault,
			claimToChannelColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			claimToChannelAllColumns,
			claimToChannelPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("model: unable to upsert claim_to_channel, could not build update column list")
		}

		ret := strmangle.SetComplement(claimToChannelAllColumns, strmangle.SetIntersect(insert, update))

		cache.query = buildUpsertQueryMySQL(dialect, "`claim_to_channel`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `claim_to_channel` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(claimToChannelType, claimToChannelMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(claimToChannelType, claimToChannelMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert for claim_to_channel")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(claimToChannelType, claimToChannelMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model: unable to retrieve unique values for claim_to_channel")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, nzUniqueCols...)
	}
	err = exec.QueryRow(cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model: unable to populate default values for claim_to_channel")
	}

CacheNoHooks:
	if !cached {
		claimToChannelUpsertCacheMut.Lock()
		claimToChannelUpsertCache[key] = cache
		claimToChannelUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single ClaimToChannel record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ClaimToChannel) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("model: no ClaimToChannel provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), claimToChannelPrimaryKeyMapping)
	sql := "DELETE FROM `claim_to_channel` WHERE `claim_id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete from claim_to_channel")
	}

	return nil
}

// DeleteAll deletes all matching rows.
func (q claimToChannelQuery) DeleteAll(exec boil.Executor) error {
	if q.Query == nil {
		return errors.New("model: no claimToChannelQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete all from claim_to_channel")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ClaimToChannelSlice) DeleteAll(exec boil.Executor) error {
	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), claimToChannelPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `claim_to_channel` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, claimToChannelPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete all from claimToChannel slice")
	}

	return nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ClaimToChannel) Reload(exec boil.Executor) error {
	ret, err := FindClaimToChannel(exec, o.ClaimID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ClaimToChannelSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ClaimToChannelSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), claimToChannelPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `claim_to_channel`.* FROM `claim_to_channel` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, claimToChannelPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in ClaimToChannelSlice")
	}

	*o = slice

	return nil
}

// ClaimToChannelExists checks if the ClaimToChannel row exists.
func ClaimToChannelExists(exec boil.Executor, claimID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `claim_to_channel` where `claim_id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, claimID)
	}
	row := exec.QueryRow(sql, claimID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if claim_to_channel exists")
	}

	return exists, nil
}

// Exists checks if the ClaimToChannel row exists.
func (o *ClaimToChannel) Exists(exec boil.Executor) (bool, error) {
	return ClaimToChannelExists(exec, o.ClaimID)
}
