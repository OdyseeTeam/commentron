// Code generated by SQLBoiler 4.16.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// ReactionType is an object representing the database table.
type ReactionType struct {
	ID        uint64    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name      string    `boil:"name" json:"name" toml:"name" yaml:"name"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *reactionTypeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L reactionTypeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ReactionTypeColumns = struct {
	ID        string
	Name      string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	Name:      "name",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

var ReactionTypeTableColumns = struct {
	ID        string
	Name      string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "reaction_type.id",
	Name:      "reaction_type.name",
	CreatedAt: "reaction_type.created_at",
	UpdatedAt: "reaction_type.updated_at",
}

// Generated where

var ReactionTypeWhere = struct {
	ID        whereHelperuint64
	Name      whereHelperstring
	CreatedAt whereHelpertime_Time
	UpdatedAt whereHelpertime_Time
}{
	ID:        whereHelperuint64{field: "`reaction_type`.`id`"},
	Name:      whereHelperstring{field: "`reaction_type`.`name`"},
	CreatedAt: whereHelpertime_Time{field: "`reaction_type`.`created_at`"},
	UpdatedAt: whereHelpertime_Time{field: "`reaction_type`.`updated_at`"},
}

// ReactionTypeRels is where relationship names are stored.
var ReactionTypeRels = struct {
	Reactions string
}{
	Reactions: "Reactions",
}

// reactionTypeR is where relationships are stored.
type reactionTypeR struct {
	Reactions ReactionSlice `boil:"Reactions" json:"Reactions" toml:"Reactions" yaml:"Reactions"`
}

// NewStruct creates a new relationship struct
func (*reactionTypeR) NewStruct() *reactionTypeR {
	return &reactionTypeR{}
}

func (r *reactionTypeR) GetReactions() ReactionSlice {
	if r == nil {
		return nil
	}
	return r.Reactions
}

// reactionTypeL is where Load methods for each relationship are stored.
type reactionTypeL struct{}

var (
	reactionTypeAllColumns            = []string{"id", "name", "created_at", "updated_at"}
	reactionTypeColumnsWithoutDefault = []string{"name"}
	reactionTypeColumnsWithDefault    = []string{"id", "created_at", "updated_at"}
	reactionTypePrimaryKeyColumns     = []string{"id"}
	reactionTypeGeneratedColumns      = []string{}
)

type (
	// ReactionTypeSlice is an alias for a slice of pointers to ReactionType.
	// This should almost always be used instead of []ReactionType.
	ReactionTypeSlice []*ReactionType

	reactionTypeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	reactionTypeType                 = reflect.TypeOf(&ReactionType{})
	reactionTypeMapping              = queries.MakeStructMapping(reactionTypeType)
	reactionTypePrimaryKeyMapping, _ = queries.BindMapping(reactionTypeType, reactionTypeMapping, reactionTypePrimaryKeyColumns)
	reactionTypeInsertCacheMut       sync.RWMutex
	reactionTypeInsertCache          = make(map[string]insertCache)
	reactionTypeUpdateCacheMut       sync.RWMutex
	reactionTypeUpdateCache          = make(map[string]updateCache)
	reactionTypeUpsertCacheMut       sync.RWMutex
	reactionTypeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single reactionType record from the query.
func (q reactionTypeQuery) One(exec boil.Executor) (*ReactionType, error) {
	o := &ReactionType{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for reaction_type")
	}

	return o, nil
}

// All returns all ReactionType records from the query.
func (q reactionTypeQuery) All(exec boil.Executor) (ReactionTypeSlice, error) {
	var o []*ReactionType

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to ReactionType slice")
	}

	return o, nil
}

// Count returns the count of all ReactionType records in the query.
func (q reactionTypeQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count reaction_type rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q reactionTypeQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if reaction_type exists")
	}

	return count > 0, nil
}

// Reactions retrieves all the reaction's Reactions with an executor.
func (o *ReactionType) Reactions(mods ...qm.QueryMod) reactionQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("`reaction`.`reaction_type_id`=?", o.ID),
	)

	return Reactions(queryMods...)
}

// LoadReactions allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (reactionTypeL) LoadReactions(e boil.Executor, singular bool, maybeReactionType interface{}, mods queries.Applicator) error {
	var slice []*ReactionType
	var object *ReactionType

	if singular {
		var ok bool
		object, ok = maybeReactionType.(*ReactionType)
		if !ok {
			object = new(ReactionType)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeReactionType)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeReactionType))
			}
		}
	} else {
		s, ok := maybeReactionType.(*[]*ReactionType)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeReactionType)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeReactionType))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &reactionTypeR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &reactionTypeR{}
			}
			args[obj.ID] = struct{}{}
		}
	}

	if len(args) == 0 {
		return nil
	}

	argsSlice := make([]interface{}, len(args))
	i := 0
	for arg := range args {
		argsSlice[i] = arg
		i++
	}

	query := NewQuery(
		qm.From(`reaction`),
		qm.WhereIn(`reaction.reaction_type_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load reaction")
	}

	var resultSlice []*Reaction
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice reaction")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on reaction")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for reaction")
	}

	if singular {
		object.R.Reactions = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &reactionR{}
			}
			foreign.R.ReactionType = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.ReactionTypeID {
				local.R.Reactions = append(local.R.Reactions, foreign)
				if foreign.R == nil {
					foreign.R = &reactionR{}
				}
				foreign.R.ReactionType = local
				break
			}
		}
	}

	return nil
}

// AddReactions adds the given related objects to the existing relationships
// of the reaction_type, optionally inserting them as new records.
// Appends related to o.R.Reactions.
// Sets related.R.ReactionType appropriately.
func (o *ReactionType) AddReactions(exec boil.Executor, insert bool, related ...*Reaction) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.ReactionTypeID = o.ID
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE `reaction` SET %s WHERE %s",
				strmangle.SetParamNames("`", "`", 0, []string{"reaction_type_id"}),
				strmangle.WhereClause("`", "`", 0, reactionPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.ReactionTypeID = o.ID
		}
	}

	if o.R == nil {
		o.R = &reactionTypeR{
			Reactions: related,
		}
	} else {
		o.R.Reactions = append(o.R.Reactions, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &reactionR{
				ReactionType: o,
			}
		} else {
			rel.R.ReactionType = o
		}
	}
	return nil
}

// ReactionTypes retrieves all the records using an executor.
func ReactionTypes(mods ...qm.QueryMod) reactionTypeQuery {
	mods = append(mods, qm.From("`reaction_type`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`reaction_type`.*"})
	}

	return reactionTypeQuery{q}
}

// FindReactionType retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindReactionType(exec boil.Executor, iD uint64, selectCols ...string) (*ReactionType, error) {
	reactionTypeObj := &ReactionType{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `reaction_type` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, reactionTypeObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from reaction_type")
	}

	return reactionTypeObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ReactionType) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no reaction_type provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(reactionTypeColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	reactionTypeInsertCacheMut.RLock()
	cache, cached := reactionTypeInsertCache[key]
	reactionTypeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			reactionTypeAllColumns,
			reactionTypeColumnsWithDefault,
			reactionTypeColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(reactionTypeType, reactionTypeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(reactionTypeType, reactionTypeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `reaction_type` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `reaction_type` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `reaction_type` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, reactionTypePrimaryKeyColumns))
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
		return errors.Wrap(err, "model: unable to insert into reaction_type")
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
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == reactionTypeMapping["id"] {
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
		return errors.Wrap(err, "model: unable to populate default values for reaction_type")
	}

CacheNoHooks:
	if !cached {
		reactionTypeInsertCacheMut.Lock()
		reactionTypeInsertCache[key] = cache
		reactionTypeInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the ReactionType.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ReactionType) Update(exec boil.Executor, columns boil.Columns) error {
	var err error
	key := makeCacheKey(columns, nil)
	reactionTypeUpdateCacheMut.RLock()
	cache, cached := reactionTypeUpdateCache[key]
	reactionTypeUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			reactionTypeAllColumns,
			reactionTypePrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return errors.New("model: unable to update reaction_type, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `reaction_type` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, reactionTypePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(reactionTypeType, reactionTypeMapping, append(wl, reactionTypePrimaryKeyColumns...))
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
		return errors.Wrap(err, "model: unable to update reaction_type row")
	}

	if !cached {
		reactionTypeUpdateCacheMut.Lock()
		reactionTypeUpdateCache[key] = cache
		reactionTypeUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAll updates all rows with the specified column values.
func (q reactionTypeQuery) UpdateAll(exec boil.Executor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "model: unable to update all for reaction_type")
	}

	return nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ReactionTypeSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reactionTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `reaction_type` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, reactionTypePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to update all in reactionType slice")
	}

	return nil
}

var mySQLReactionTypeUniqueColumns = []string{
	"id",
	"name",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ReactionType) Upsert(exec boil.Executor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model: no reaction_type provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(reactionTypeColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLReactionTypeUniqueColumns, o)

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

	reactionTypeUpsertCacheMut.RLock()
	cache, cached := reactionTypeUpsertCache[key]
	reactionTypeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			reactionTypeAllColumns,
			reactionTypeColumnsWithDefault,
			reactionTypeColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			reactionTypeAllColumns,
			reactionTypePrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("model: unable to upsert reaction_type, could not build update column list")
		}

		ret := strmangle.SetComplement(reactionTypeAllColumns, strmangle.SetIntersect(insert, update))

		cache.query = buildUpsertQueryMySQL(dialect, "`reaction_type`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `reaction_type` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(reactionTypeType, reactionTypeMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(reactionTypeType, reactionTypeMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert for reaction_type")
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
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == reactionTypeMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(reactionTypeType, reactionTypeMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model: unable to retrieve unique values for reaction_type")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, nzUniqueCols...)
	}
	err = exec.QueryRow(cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model: unable to populate default values for reaction_type")
	}

CacheNoHooks:
	if !cached {
		reactionTypeUpsertCacheMut.Lock()
		reactionTypeUpsertCache[key] = cache
		reactionTypeUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single ReactionType record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ReactionType) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("model: no ReactionType provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), reactionTypePrimaryKeyMapping)
	sql := "DELETE FROM `reaction_type` WHERE `id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete from reaction_type")
	}

	return nil
}

// DeleteAll deletes all matching rows.
func (q reactionTypeQuery) DeleteAll(exec boil.Executor) error {
	if q.Query == nil {
		return errors.New("model: no reactionTypeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete all from reaction_type")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ReactionTypeSlice) DeleteAll(exec boil.Executor) error {
	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reactionTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `reaction_type` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, reactionTypePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete all from reactionType slice")
	}

	return nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ReactionType) Reload(exec boil.Executor) error {
	ret, err := FindReactionType(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ReactionTypeSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ReactionTypeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reactionTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `reaction_type`.* FROM `reaction_type` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, reactionTypePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in ReactionTypeSlice")
	}

	*o = slice

	return nil
}

// ReactionTypeExists checks if the ReactionType row exists.
func ReactionTypeExists(exec boil.Executor, iD uint64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `reaction_type` where `id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if reaction_type exists")
	}

	return exists, nil
}

// Exists checks if the ReactionType row exists.
func (o *ReactionType) Exists(exec boil.Executor) (bool, error) {
	return ReactionTypeExists(exec, o.ID)
}
