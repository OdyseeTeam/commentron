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

// DelegatedModerator is an object representing the database table.
type DelegatedModerator struct {
	ID               uint64    `boil:"id" json:"id" toml:"id" yaml:"id"`
	ModChannelID     string    `boil:"mod_channel_id" json:"mod_channel_id" toml:"mod_channel_id" yaml:"mod_channel_id"`
	CreatorChannelID string    `boil:"creator_channel_id" json:"creator_channel_id" toml:"creator_channel_id" yaml:"creator_channel_id"`
	Permissons       uint64    `boil:"permissons" json:"permissons" toml:"permissons" yaml:"permissons"`
	CreatedAt        time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt        time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *delegatedModeratorR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L delegatedModeratorL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var DelegatedModeratorColumns = struct {
	ID               string
	ModChannelID     string
	CreatorChannelID string
	Permissons       string
	CreatedAt        string
	UpdatedAt        string
}{
	ID:               "id",
	ModChannelID:     "mod_channel_id",
	CreatorChannelID: "creator_channel_id",
	Permissons:       "permissons",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
}

var DelegatedModeratorTableColumns = struct {
	ID               string
	ModChannelID     string
	CreatorChannelID string
	Permissons       string
	CreatedAt        string
	UpdatedAt        string
}{
	ID:               "delegated_moderator.id",
	ModChannelID:     "delegated_moderator.mod_channel_id",
	CreatorChannelID: "delegated_moderator.creator_channel_id",
	Permissons:       "delegated_moderator.permissons",
	CreatedAt:        "delegated_moderator.created_at",
	UpdatedAt:        "delegated_moderator.updated_at",
}

// Generated where

var DelegatedModeratorWhere = struct {
	ID               whereHelperuint64
	ModChannelID     whereHelperstring
	CreatorChannelID whereHelperstring
	Permissons       whereHelperuint64
	CreatedAt        whereHelpertime_Time
	UpdatedAt        whereHelpertime_Time
}{
	ID:               whereHelperuint64{field: "`delegated_moderator`.`id`"},
	ModChannelID:     whereHelperstring{field: "`delegated_moderator`.`mod_channel_id`"},
	CreatorChannelID: whereHelperstring{field: "`delegated_moderator`.`creator_channel_id`"},
	Permissons:       whereHelperuint64{field: "`delegated_moderator`.`permissons`"},
	CreatedAt:        whereHelpertime_Time{field: "`delegated_moderator`.`created_at`"},
	UpdatedAt:        whereHelpertime_Time{field: "`delegated_moderator`.`updated_at`"},
}

// DelegatedModeratorRels is where relationship names are stored.
var DelegatedModeratorRels = struct {
	ModChannel     string
	CreatorChannel string
}{
	ModChannel:     "ModChannel",
	CreatorChannel: "CreatorChannel",
}

// delegatedModeratorR is where relationships are stored.
type delegatedModeratorR struct {
	ModChannel     *Channel `boil:"ModChannel" json:"ModChannel" toml:"ModChannel" yaml:"ModChannel"`
	CreatorChannel *Channel `boil:"CreatorChannel" json:"CreatorChannel" toml:"CreatorChannel" yaml:"CreatorChannel"`
}

// NewStruct creates a new relationship struct
func (*delegatedModeratorR) NewStruct() *delegatedModeratorR {
	return &delegatedModeratorR{}
}

func (r *delegatedModeratorR) GetModChannel() *Channel {
	if r == nil {
		return nil
	}
	return r.ModChannel
}

func (r *delegatedModeratorR) GetCreatorChannel() *Channel {
	if r == nil {
		return nil
	}
	return r.CreatorChannel
}

// delegatedModeratorL is where Load methods for each relationship are stored.
type delegatedModeratorL struct{}

var (
	delegatedModeratorAllColumns            = []string{"id", "mod_channel_id", "creator_channel_id", "permissons", "created_at", "updated_at"}
	delegatedModeratorColumnsWithoutDefault = []string{"mod_channel_id", "creator_channel_id"}
	delegatedModeratorColumnsWithDefault    = []string{"id", "permissons", "created_at", "updated_at"}
	delegatedModeratorPrimaryKeyColumns     = []string{"id"}
	delegatedModeratorGeneratedColumns      = []string{}
)

type (
	// DelegatedModeratorSlice is an alias for a slice of pointers to DelegatedModerator.
	// This should almost always be used instead of []DelegatedModerator.
	DelegatedModeratorSlice []*DelegatedModerator

	delegatedModeratorQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	delegatedModeratorType                 = reflect.TypeOf(&DelegatedModerator{})
	delegatedModeratorMapping              = queries.MakeStructMapping(delegatedModeratorType)
	delegatedModeratorPrimaryKeyMapping, _ = queries.BindMapping(delegatedModeratorType, delegatedModeratorMapping, delegatedModeratorPrimaryKeyColumns)
	delegatedModeratorInsertCacheMut       sync.RWMutex
	delegatedModeratorInsertCache          = make(map[string]insertCache)
	delegatedModeratorUpdateCacheMut       sync.RWMutex
	delegatedModeratorUpdateCache          = make(map[string]updateCache)
	delegatedModeratorUpsertCacheMut       sync.RWMutex
	delegatedModeratorUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single delegatedModerator record from the query.
func (q delegatedModeratorQuery) One(exec boil.Executor) (*DelegatedModerator, error) {
	o := &DelegatedModerator{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for delegated_moderator")
	}

	return o, nil
}

// All returns all DelegatedModerator records from the query.
func (q delegatedModeratorQuery) All(exec boil.Executor) (DelegatedModeratorSlice, error) {
	var o []*DelegatedModerator

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to DelegatedModerator slice")
	}

	return o, nil
}

// Count returns the count of all DelegatedModerator records in the query.
func (q delegatedModeratorQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count delegated_moderator rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q delegatedModeratorQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if delegated_moderator exists")
	}

	return count > 0, nil
}

// ModChannel pointed to by the foreign key.
func (o *DelegatedModerator) ModChannel(mods ...qm.QueryMod) channelQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`claim_id` = ?", o.ModChannelID),
	}

	queryMods = append(queryMods, mods...)

	return Channels(queryMods...)
}

// CreatorChannel pointed to by the foreign key.
func (o *DelegatedModerator) CreatorChannel(mods ...qm.QueryMod) channelQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`claim_id` = ?", o.CreatorChannelID),
	}

	queryMods = append(queryMods, mods...)

	return Channels(queryMods...)
}

// LoadModChannel allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (delegatedModeratorL) LoadModChannel(e boil.Executor, singular bool, maybeDelegatedModerator interface{}, mods queries.Applicator) error {
	var slice []*DelegatedModerator
	var object *DelegatedModerator

	if singular {
		var ok bool
		object, ok = maybeDelegatedModerator.(*DelegatedModerator)
		if !ok {
			object = new(DelegatedModerator)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeDelegatedModerator)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeDelegatedModerator))
			}
		}
	} else {
		s, ok := maybeDelegatedModerator.(*[]*DelegatedModerator)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeDelegatedModerator)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeDelegatedModerator))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &delegatedModeratorR{}
		}
		args[object.ModChannelID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &delegatedModeratorR{}
			}

			args[obj.ModChannelID] = struct{}{}

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
		qm.From(`channel`),
		qm.WhereIn(`channel.claim_id in ?`, argsSlice...),
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
		foreign.R.ModChannelDelegatedModerators = append(foreign.R.ModChannelDelegatedModerators, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ModChannelID == foreign.ClaimID {
				local.R.ModChannel = foreign
				if foreign.R == nil {
					foreign.R = &channelR{}
				}
				foreign.R.ModChannelDelegatedModerators = append(foreign.R.ModChannelDelegatedModerators, local)
				break
			}
		}
	}

	return nil
}

// LoadCreatorChannel allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (delegatedModeratorL) LoadCreatorChannel(e boil.Executor, singular bool, maybeDelegatedModerator interface{}, mods queries.Applicator) error {
	var slice []*DelegatedModerator
	var object *DelegatedModerator

	if singular {
		var ok bool
		object, ok = maybeDelegatedModerator.(*DelegatedModerator)
		if !ok {
			object = new(DelegatedModerator)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeDelegatedModerator)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeDelegatedModerator))
			}
		}
	} else {
		s, ok := maybeDelegatedModerator.(*[]*DelegatedModerator)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeDelegatedModerator)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeDelegatedModerator))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &delegatedModeratorR{}
		}
		args[object.CreatorChannelID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &delegatedModeratorR{}
			}

			args[obj.CreatorChannelID] = struct{}{}

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
		qm.From(`channel`),
		qm.WhereIn(`channel.claim_id in ?`, argsSlice...),
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
		object.R.CreatorChannel = foreign
		if foreign.R == nil {
			foreign.R = &channelR{}
		}
		foreign.R.CreatorChannelDelegatedModerators = append(foreign.R.CreatorChannelDelegatedModerators, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.CreatorChannelID == foreign.ClaimID {
				local.R.CreatorChannel = foreign
				if foreign.R == nil {
					foreign.R = &channelR{}
				}
				foreign.R.CreatorChannelDelegatedModerators = append(foreign.R.CreatorChannelDelegatedModerators, local)
				break
			}
		}
	}

	return nil
}

// SetModChannel of the delegatedModerator to the related item.
// Sets o.R.ModChannel to related.
// Adds o to related.R.ModChannelDelegatedModerators.
func (o *DelegatedModerator) SetModChannel(exec boil.Executor, insert bool, related *Channel) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `delegated_moderator` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"mod_channel_id"}),
		strmangle.WhereClause("`", "`", 0, delegatedModeratorPrimaryKeyColumns),
	)
	values := []interface{}{related.ClaimID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ModChannelID = related.ClaimID
	if o.R == nil {
		o.R = &delegatedModeratorR{
			ModChannel: related,
		}
	} else {
		o.R.ModChannel = related
	}

	if related.R == nil {
		related.R = &channelR{
			ModChannelDelegatedModerators: DelegatedModeratorSlice{o},
		}
	} else {
		related.R.ModChannelDelegatedModerators = append(related.R.ModChannelDelegatedModerators, o)
	}

	return nil
}

// SetCreatorChannel of the delegatedModerator to the related item.
// Sets o.R.CreatorChannel to related.
// Adds o to related.R.CreatorChannelDelegatedModerators.
func (o *DelegatedModerator) SetCreatorChannel(exec boil.Executor, insert bool, related *Channel) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `delegated_moderator` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"creator_channel_id"}),
		strmangle.WhereClause("`", "`", 0, delegatedModeratorPrimaryKeyColumns),
	)
	values := []interface{}{related.ClaimID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CreatorChannelID = related.ClaimID
	if o.R == nil {
		o.R = &delegatedModeratorR{
			CreatorChannel: related,
		}
	} else {
		o.R.CreatorChannel = related
	}

	if related.R == nil {
		related.R = &channelR{
			CreatorChannelDelegatedModerators: DelegatedModeratorSlice{o},
		}
	} else {
		related.R.CreatorChannelDelegatedModerators = append(related.R.CreatorChannelDelegatedModerators, o)
	}

	return nil
}

// DelegatedModerators retrieves all the records using an executor.
func DelegatedModerators(mods ...qm.QueryMod) delegatedModeratorQuery {
	mods = append(mods, qm.From("`delegated_moderator`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`delegated_moderator`.*"})
	}

	return delegatedModeratorQuery{q}
}

// FindDelegatedModerator retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindDelegatedModerator(exec boil.Executor, iD uint64, selectCols ...string) (*DelegatedModerator, error) {
	delegatedModeratorObj := &DelegatedModerator{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `delegated_moderator` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, delegatedModeratorObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from delegated_moderator")
	}

	return delegatedModeratorObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *DelegatedModerator) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no delegated_moderator provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(delegatedModeratorColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	delegatedModeratorInsertCacheMut.RLock()
	cache, cached := delegatedModeratorInsertCache[key]
	delegatedModeratorInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			delegatedModeratorAllColumns,
			delegatedModeratorColumnsWithDefault,
			delegatedModeratorColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(delegatedModeratorType, delegatedModeratorMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(delegatedModeratorType, delegatedModeratorMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `delegated_moderator` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `delegated_moderator` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `delegated_moderator` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, delegatedModeratorPrimaryKeyColumns))
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
		return errors.Wrap(err, "model: unable to insert into delegated_moderator")
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
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == delegatedModeratorMapping["id"] {
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
		return errors.Wrap(err, "model: unable to populate default values for delegated_moderator")
	}

CacheNoHooks:
	if !cached {
		delegatedModeratorInsertCacheMut.Lock()
		delegatedModeratorInsertCache[key] = cache
		delegatedModeratorInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the DelegatedModerator.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *DelegatedModerator) Update(exec boil.Executor, columns boil.Columns) error {
	var err error
	key := makeCacheKey(columns, nil)
	delegatedModeratorUpdateCacheMut.RLock()
	cache, cached := delegatedModeratorUpdateCache[key]
	delegatedModeratorUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			delegatedModeratorAllColumns,
			delegatedModeratorPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return errors.New("model: unable to update delegated_moderator, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `delegated_moderator` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, delegatedModeratorPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(delegatedModeratorType, delegatedModeratorMapping, append(wl, delegatedModeratorPrimaryKeyColumns...))
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
		return errors.Wrap(err, "model: unable to update delegated_moderator row")
	}

	if !cached {
		delegatedModeratorUpdateCacheMut.Lock()
		delegatedModeratorUpdateCache[key] = cache
		delegatedModeratorUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAll updates all rows with the specified column values.
func (q delegatedModeratorQuery) UpdateAll(exec boil.Executor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "model: unable to update all for delegated_moderator")
	}

	return nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o DelegatedModeratorSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), delegatedModeratorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `delegated_moderator` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, delegatedModeratorPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to update all in delegatedModerator slice")
	}

	return nil
}

var mySQLDelegatedModeratorUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *DelegatedModerator) Upsert(exec boil.Executor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model: no delegated_moderator provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(delegatedModeratorColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLDelegatedModeratorUniqueColumns, o)

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

	delegatedModeratorUpsertCacheMut.RLock()
	cache, cached := delegatedModeratorUpsertCache[key]
	delegatedModeratorUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			delegatedModeratorAllColumns,
			delegatedModeratorColumnsWithDefault,
			delegatedModeratorColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			delegatedModeratorAllColumns,
			delegatedModeratorPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("model: unable to upsert delegated_moderator, could not build update column list")
		}

		ret := strmangle.SetComplement(delegatedModeratorAllColumns, strmangle.SetIntersect(insert, update))

		cache.query = buildUpsertQueryMySQL(dialect, "`delegated_moderator`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `delegated_moderator` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(delegatedModeratorType, delegatedModeratorMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(delegatedModeratorType, delegatedModeratorMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert for delegated_moderator")
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
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == delegatedModeratorMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(delegatedModeratorType, delegatedModeratorMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model: unable to retrieve unique values for delegated_moderator")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, nzUniqueCols...)
	}
	err = exec.QueryRow(cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model: unable to populate default values for delegated_moderator")
	}

CacheNoHooks:
	if !cached {
		delegatedModeratorUpsertCacheMut.Lock()
		delegatedModeratorUpsertCache[key] = cache
		delegatedModeratorUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single DelegatedModerator record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *DelegatedModerator) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("model: no DelegatedModerator provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), delegatedModeratorPrimaryKeyMapping)
	sql := "DELETE FROM `delegated_moderator` WHERE `id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete from delegated_moderator")
	}

	return nil
}

// DeleteAll deletes all matching rows.
func (q delegatedModeratorQuery) DeleteAll(exec boil.Executor) error {
	if q.Query == nil {
		return errors.New("model: no delegatedModeratorQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete all from delegated_moderator")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o DelegatedModeratorSlice) DeleteAll(exec boil.Executor) error {
	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), delegatedModeratorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `delegated_moderator` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, delegatedModeratorPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete all from delegatedModerator slice")
	}

	return nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *DelegatedModerator) Reload(exec boil.Executor) error {
	ret, err := FindDelegatedModerator(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *DelegatedModeratorSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := DelegatedModeratorSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), delegatedModeratorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `delegated_moderator`.* FROM `delegated_moderator` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, delegatedModeratorPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in DelegatedModeratorSlice")
	}

	*o = slice

	return nil
}

// DelegatedModeratorExists checks if the DelegatedModerator row exists.
func DelegatedModeratorExists(exec boil.Executor, iD uint64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `delegated_moderator` where `id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if delegated_moderator exists")
	}

	return exists, nil
}

// Exists checks if the DelegatedModerator row exists.
func (o *DelegatedModerator) Exists(exec boil.Executor) (bool, error) {
	return DelegatedModeratorExists(exec, o.ID)
}
