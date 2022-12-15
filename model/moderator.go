// Code generated by SQLBoiler 4.14.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// Moderator is an object representing the database table.
type Moderator struct {
	ID           uint64      `boil:"id" json:"id" toml:"id" yaml:"id"`
	ModChannelID null.String `boil:"mod_channel_id" json:"mod_channel_id,omitempty" toml:"mod_channel_id" yaml:"mod_channel_id,omitempty"`
	ModLevel     int64       `boil:"mod_level" json:"mod_level" toml:"mod_level" yaml:"mod_level"`
	CreatedAt    time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt    time.Time   `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *moderatorR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L moderatorL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ModeratorColumns = struct {
	ID           string
	ModChannelID string
	ModLevel     string
	CreatedAt    string
	UpdatedAt    string
}{
	ID:           "id",
	ModChannelID: "mod_channel_id",
	ModLevel:     "mod_level",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

var ModeratorTableColumns = struct {
	ID           string
	ModChannelID string
	ModLevel     string
	CreatedAt    string
	UpdatedAt    string
}{
	ID:           "moderator.id",
	ModChannelID: "moderator.mod_channel_id",
	ModLevel:     "moderator.mod_level",
	CreatedAt:    "moderator.created_at",
	UpdatedAt:    "moderator.updated_at",
}

// Generated where

var ModeratorWhere = struct {
	ID           whereHelperuint64
	ModChannelID whereHelpernull_String
	ModLevel     whereHelperint64
	CreatedAt    whereHelpertime_Time
	UpdatedAt    whereHelpertime_Time
}{
	ID:           whereHelperuint64{field: "`moderator`.`id`"},
	ModChannelID: whereHelpernull_String{field: "`moderator`.`mod_channel_id`"},
	ModLevel:     whereHelperint64{field: "`moderator`.`mod_level`"},
	CreatedAt:    whereHelpertime_Time{field: "`moderator`.`created_at`"},
	UpdatedAt:    whereHelpertime_Time{field: "`moderator`.`updated_at`"},
}

// ModeratorRels is where relationship names are stored.
var ModeratorRels = struct {
	ModChannel string
}{
	ModChannel: "ModChannel",
}

// moderatorR is where relationships are stored.
type moderatorR struct {
	ModChannel *Channel `boil:"ModChannel" json:"ModChannel" toml:"ModChannel" yaml:"ModChannel"`
}

// NewStruct creates a new relationship struct
func (*moderatorR) NewStruct() *moderatorR {
	return &moderatorR{}
}

func (r *moderatorR) GetModChannel() *Channel {
	if r == nil {
		return nil
	}
	return r.ModChannel
}

// moderatorL is where Load methods for each relationship are stored.
type moderatorL struct{}

var (
	moderatorAllColumns            = []string{"id", "mod_channel_id", "mod_level", "created_at", "updated_at"}
	moderatorColumnsWithoutDefault = []string{"mod_channel_id"}
	moderatorColumnsWithDefault    = []string{"id", "mod_level", "created_at", "updated_at"}
	moderatorPrimaryKeyColumns     = []string{"id"}
	moderatorGeneratedColumns      = []string{}
)

type (
	// ModeratorSlice is an alias for a slice of pointers to Moderator.
	// This should almost always be used instead of []Moderator.
	ModeratorSlice []*Moderator

	moderatorQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	moderatorType                 = reflect.TypeOf(&Moderator{})
	moderatorMapping              = queries.MakeStructMapping(moderatorType)
	moderatorPrimaryKeyMapping, _ = queries.BindMapping(moderatorType, moderatorMapping, moderatorPrimaryKeyColumns)
	moderatorInsertCacheMut       sync.RWMutex
	moderatorInsertCache          = make(map[string]insertCache)
	moderatorUpdateCacheMut       sync.RWMutex
	moderatorUpdateCache          = make(map[string]updateCache)
	moderatorUpsertCacheMut       sync.RWMutex
	moderatorUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single moderator record from the query.
func (q moderatorQuery) One(exec boil.Executor) (*Moderator, error) {
	o := &Moderator{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for moderator")
	}

	return o, nil
}

// All returns all Moderator records from the query.
func (q moderatorQuery) All(exec boil.Executor) (ModeratorSlice, error) {
	var o []*Moderator

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to Moderator slice")
	}

	return o, nil
}

// Count returns the count of all Moderator records in the query.
func (q moderatorQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count moderator rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q moderatorQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if moderator exists")
	}

	return count > 0, nil
}

// ModChannel pointed to by the foreign key.
func (o *Moderator) ModChannel(mods ...qm.QueryMod) channelQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`claim_id` = ?", o.ModChannelID),
	}

	queryMods = append(queryMods, mods...)

	return Channels(queryMods...)
}

// LoadModChannel allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (moderatorL) LoadModChannel(e boil.Executor, singular bool, maybeModerator interface{}, mods queries.Applicator) error {
	var slice []*Moderator
	var object *Moderator

	if singular {
		var ok bool
		object, ok = maybeModerator.(*Moderator)
		if !ok {
			object = new(Moderator)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeModerator)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeModerator))
			}
		}
	} else {
		s, ok := maybeModerator.(*[]*Moderator)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeModerator)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeModerator))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &moderatorR{}
		}
		if !queries.IsNil(object.ModChannelID) {
			args = append(args, object.ModChannelID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &moderatorR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.ModChannelID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.ModChannelID) {
				args = append(args, obj.ModChannelID)
			}

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`channel`),
		qm.WhereIn(`channel.claim_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Channel")
	}

	var resultSlice []*Channel
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Channel")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for channel")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for channel")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.ModChannel = foreign
		if foreign.R == nil {
			foreign.R = &channelR{}
		}
		foreign.R.ModChannelModerators = append(foreign.R.ModChannelModerators, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.ModChannelID, foreign.ClaimID) {
				local.R.ModChannel = foreign
				if foreign.R == nil {
					foreign.R = &channelR{}
				}
				foreign.R.ModChannelModerators = append(foreign.R.ModChannelModerators, local)
				break
			}
		}
	}

	return nil
}

// SetModChannel of the moderator to the related item.
// Sets o.R.ModChannel to related.
// Adds o to related.R.ModChannelModerators.
func (o *Moderator) SetModChannel(exec boil.Executor, insert bool, related *Channel) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `moderator` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"mod_channel_id"}),
		strmangle.WhereClause("`", "`", 0, moderatorPrimaryKeyColumns),
	)
	values := []interface{}{related.ClaimID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.ModChannelID, related.ClaimID)
	if o.R == nil {
		o.R = &moderatorR{
			ModChannel: related,
		}
	} else {
		o.R.ModChannel = related
	}

	if related.R == nil {
		related.R = &channelR{
			ModChannelModerators: ModeratorSlice{o},
		}
	} else {
		related.R.ModChannelModerators = append(related.R.ModChannelModerators, o)
	}

	return nil
}

// RemoveModChannel relationship.
// Sets o.R.ModChannel to nil.
// Removes o from all passed in related items' relationships struct.
func (o *Moderator) RemoveModChannel(exec boil.Executor, related *Channel) error {
	var err error

	queries.SetScanner(&o.ModChannelID, nil)
	if err = o.Update(exec, boil.Whitelist("mod_channel_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.ModChannel = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.ModChannelModerators {
		if queries.Equal(o.ModChannelID, ri.ModChannelID) {
			continue
		}

		ln := len(related.R.ModChannelModerators)
		if ln > 1 && i < ln-1 {
			related.R.ModChannelModerators[i] = related.R.ModChannelModerators[ln-1]
		}
		related.R.ModChannelModerators = related.R.ModChannelModerators[:ln-1]
		break
	}
	return nil
}

// Moderators retrieves all the records using an executor.
func Moderators(mods ...qm.QueryMod) moderatorQuery {
	mods = append(mods, qm.From("`moderator`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`moderator`.*"})
	}

	return moderatorQuery{q}
}

// FindModerator retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindModerator(exec boil.Executor, iD uint64, selectCols ...string) (*Moderator, error) {
	moderatorObj := &Moderator{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `moderator` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, moderatorObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from moderator")
	}

	return moderatorObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Moderator) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no moderator provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(moderatorColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	moderatorInsertCacheMut.RLock()
	cache, cached := moderatorInsertCache[key]
	moderatorInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			moderatorAllColumns,
			moderatorColumnsWithDefault,
			moderatorColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(moderatorType, moderatorMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(moderatorType, moderatorMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `moderator` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `moderator` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `moderator` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, moderatorPrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}
	result, err := exec.Exec(cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "model: unable to insert into moderator")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = uint64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == moderatorMapping["id"] {
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
		return errors.Wrap(err, "model: unable to populate default values for moderator")
	}

CacheNoHooks:
	if !cached {
		moderatorInsertCacheMut.Lock()
		moderatorInsertCache[key] = cache
		moderatorInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Moderator.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Moderator) Update(exec boil.Executor, columns boil.Columns) error {
	var err error
	key := makeCacheKey(columns, nil)
	moderatorUpdateCacheMut.RLock()
	cache, cached := moderatorUpdateCache[key]
	moderatorUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			moderatorAllColumns,
			moderatorPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return errors.New("model: unable to update moderator, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `moderator` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, moderatorPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(moderatorType, moderatorMapping, append(wl, moderatorPrimaryKeyColumns...))
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
		return errors.Wrap(err, "model: unable to update moderator row")
	}

	if !cached {
		moderatorUpdateCacheMut.Lock()
		moderatorUpdateCache[key] = cache
		moderatorUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAll updates all rows with the specified column values.
func (q moderatorQuery) UpdateAll(exec boil.Executor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "model: unable to update all for moderator")
	}

	return nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ModeratorSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), moderatorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `moderator` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, moderatorPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to update all in moderator slice")
	}

	return nil
}

var mySQLModeratorUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Moderator) Upsert(exec boil.Executor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model: no moderator provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(moderatorColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLModeratorUniqueColumns, o)

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

	moderatorUpsertCacheMut.RLock()
	cache, cached := moderatorUpsertCache[key]
	moderatorUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			moderatorAllColumns,
			moderatorColumnsWithDefault,
			moderatorColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			moderatorAllColumns,
			moderatorPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("model: unable to upsert moderator, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`moderator`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `moderator` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(moderatorType, moderatorMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(moderatorType, moderatorMapping, ret)
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
	result, err := exec.Exec(cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "model: unable to upsert for moderator")
	}

	var lastID int64
	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = uint64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == moderatorMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(moderatorType, moderatorMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model: unable to retrieve unique values for moderator")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, nzUniqueCols...)
	}
	err = exec.QueryRow(cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model: unable to populate default values for moderator")
	}

CacheNoHooks:
	if !cached {
		moderatorUpsertCacheMut.Lock()
		moderatorUpsertCache[key] = cache
		moderatorUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Moderator record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Moderator) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("model: no Moderator provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), moderatorPrimaryKeyMapping)
	sql := "DELETE FROM `moderator` WHERE `id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete from moderator")
	}

	return nil
}

// DeleteAll deletes all matching rows.
func (q moderatorQuery) DeleteAll(exec boil.Executor) error {
	if q.Query == nil {
		return errors.New("model: no moderatorQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete all from moderator")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ModeratorSlice) DeleteAll(exec boil.Executor) error {
	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), moderatorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `moderator` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, moderatorPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete all from moderator slice")
	}

	return nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Moderator) Reload(exec boil.Executor) error {
	ret, err := FindModerator(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ModeratorSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ModeratorSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), moderatorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `moderator`.* FROM `moderator` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, moderatorPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in ModeratorSlice")
	}

	*o = slice

	return nil
}

// ModeratorExists checks if the Moderator row exists.
func ModeratorExists(exec boil.Executor, iD uint64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `moderator` where `id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if moderator exists")
	}

	return exists, nil
}

// Exists checks if the Moderator row exists.
func (o *Moderator) Exists(exec boil.Executor) (bool, error) {
	return ModeratorExists(exec, o.ID)
}
