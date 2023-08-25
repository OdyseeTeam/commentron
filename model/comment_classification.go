// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// CommentClassification is an object representing the database table.
type CommentClassification struct {
	CommentID        string      `boil:"comment_id" json:"comment_id" toml:"comment_id" yaml:"comment_id"`
	Toxicity         float32     `boil:"toxicity" json:"toxicity" toml:"toxicity" yaml:"toxicity"`
	SevereToxicity   float32     `boil:"severe_toxicity" json:"severe_toxicity" toml:"severe_toxicity" yaml:"severe_toxicity"`
	Obscene          float32     `boil:"obscene" json:"obscene" toml:"obscene" yaml:"obscene"`
	IdentityAttack   float32     `boil:"identity_attack" json:"identity_attack" toml:"identity_attack" yaml:"identity_attack"`
	Insult           float32     `boil:"insult" json:"insult" toml:"insult" yaml:"insult"`
	Threat           float32     `boil:"threat" json:"threat" toml:"threat" yaml:"threat"`
	SexualExplicit   float32     `boil:"sexual_explicit" json:"sexual_explicit" toml:"sexual_explicit" yaml:"sexual_explicit"`
	Nazi             float32     `boil:"nazi" json:"nazi" toml:"nazi" yaml:"nazi"`
	Doxx             float32     `boil:"doxx" json:"doxx" toml:"doxx" yaml:"doxx"`
	IsReviewed       null.Bool   `boil:"is_reviewed" json:"is_reviewed,omitempty" toml:"is_reviewed" yaml:"is_reviewed,omitempty"`
	ReviewerApproved null.Bool   `boil:"reviewer_approved" json:"reviewer_approved,omitempty" toml:"reviewer_approved" yaml:"reviewer_approved,omitempty"`
	Timestamp        int         `boil:"timestamp" json:"timestamp" toml:"timestamp" yaml:"timestamp"`
	CreatedAt        time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt        time.Time   `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	ModelIdent       null.String `boil:"model_ident" json:"model_ident,omitempty" toml:"model_ident" yaml:"model_ident,omitempty"`

	R *commentClassificationR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L commentClassificationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CommentClassificationColumns = struct {
	CommentID        string
	Toxicity         string
	SevereToxicity   string
	Obscene          string
	IdentityAttack   string
	Insult           string
	Threat           string
	SexualExplicit   string
	Nazi             string
	Doxx             string
	IsReviewed       string
	ReviewerApproved string
	Timestamp        string
	CreatedAt        string
	UpdatedAt        string
	ModelIdent       string
}{
	CommentID:        "comment_id",
	Toxicity:         "toxicity",
	SevereToxicity:   "severe_toxicity",
	Obscene:          "obscene",
	IdentityAttack:   "identity_attack",
	Insult:           "insult",
	Threat:           "threat",
	SexualExplicit:   "sexual_explicit",
	Nazi:             "nazi",
	Doxx:             "doxx",
	IsReviewed:       "is_reviewed",
	ReviewerApproved: "reviewer_approved",
	Timestamp:        "timestamp",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	ModelIdent:       "model_ident",
}

var CommentClassificationTableColumns = struct {
	CommentID        string
	Toxicity         string
	SevereToxicity   string
	Obscene          string
	IdentityAttack   string
	Insult           string
	Threat           string
	SexualExplicit   string
	Nazi             string
	Doxx             string
	IsReviewed       string
	ReviewerApproved string
	Timestamp        string
	CreatedAt        string
	UpdatedAt        string
	ModelIdent       string
}{
	CommentID:        "comment_classification.comment_id",
	Toxicity:         "comment_classification.toxicity",
	SevereToxicity:   "comment_classification.severe_toxicity",
	Obscene:          "comment_classification.obscene",
	IdentityAttack:   "comment_classification.identity_attack",
	Insult:           "comment_classification.insult",
	Threat:           "comment_classification.threat",
	SexualExplicit:   "comment_classification.sexual_explicit",
	Nazi:             "comment_classification.nazi",
	Doxx:             "comment_classification.doxx",
	IsReviewed:       "comment_classification.is_reviewed",
	ReviewerApproved: "comment_classification.reviewer_approved",
	Timestamp:        "comment_classification.timestamp",
	CreatedAt:        "comment_classification.created_at",
	UpdatedAt:        "comment_classification.updated_at",
	ModelIdent:       "comment_classification.model_ident",
}

// Generated where

type whereHelperfloat32 struct{ field string }

func (w whereHelperfloat32) EQ(x float32) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperfloat32) NEQ(x float32) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelperfloat32) LT(x float32) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperfloat32) LTE(x float32) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelperfloat32) GT(x float32) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperfloat32) GTE(x float32) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelperfloat32) IN(slice []float32) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperfloat32) NIN(slice []float32) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var CommentClassificationWhere = struct {
	CommentID        whereHelperstring
	Toxicity         whereHelperfloat32
	SevereToxicity   whereHelperfloat32
	Obscene          whereHelperfloat32
	IdentityAttack   whereHelperfloat32
	Insult           whereHelperfloat32
	Threat           whereHelperfloat32
	SexualExplicit   whereHelperfloat32
	Nazi             whereHelperfloat32
	Doxx             whereHelperfloat32
	IsReviewed       whereHelpernull_Bool
	ReviewerApproved whereHelpernull_Bool
	Timestamp        whereHelperint
	CreatedAt        whereHelpertime_Time
	UpdatedAt        whereHelpertime_Time
	ModelIdent       whereHelpernull_String
}{
	CommentID:        whereHelperstring{field: "`comment_classification`.`comment_id`"},
	Toxicity:         whereHelperfloat32{field: "`comment_classification`.`toxicity`"},
	SevereToxicity:   whereHelperfloat32{field: "`comment_classification`.`severe_toxicity`"},
	Obscene:          whereHelperfloat32{field: "`comment_classification`.`obscene`"},
	IdentityAttack:   whereHelperfloat32{field: "`comment_classification`.`identity_attack`"},
	Insult:           whereHelperfloat32{field: "`comment_classification`.`insult`"},
	Threat:           whereHelperfloat32{field: "`comment_classification`.`threat`"},
	SexualExplicit:   whereHelperfloat32{field: "`comment_classification`.`sexual_explicit`"},
	Nazi:             whereHelperfloat32{field: "`comment_classification`.`nazi`"},
	Doxx:             whereHelperfloat32{field: "`comment_classification`.`doxx`"},
	IsReviewed:       whereHelpernull_Bool{field: "`comment_classification`.`is_reviewed`"},
	ReviewerApproved: whereHelpernull_Bool{field: "`comment_classification`.`reviewer_approved`"},
	Timestamp:        whereHelperint{field: "`comment_classification`.`timestamp`"},
	CreatedAt:        whereHelpertime_Time{field: "`comment_classification`.`created_at`"},
	UpdatedAt:        whereHelpertime_Time{field: "`comment_classification`.`updated_at`"},
	ModelIdent:       whereHelpernull_String{field: "`comment_classification`.`model_ident`"},
}

// CommentClassificationRels is where relationship names are stored.
var CommentClassificationRels = struct {
	Comment string
}{
	Comment: "Comment",
}

// commentClassificationR is where relationships are stored.
type commentClassificationR struct {
	Comment *Comment `boil:"Comment" json:"Comment" toml:"Comment" yaml:"Comment"`
}

// NewStruct creates a new relationship struct
func (*commentClassificationR) NewStruct() *commentClassificationR {
	return &commentClassificationR{}
}

func (r *commentClassificationR) GetComment() *Comment {
	if r == nil {
		return nil
	}
	return r.Comment
}

// commentClassificationL is where Load methods for each relationship are stored.
type commentClassificationL struct{}

var (
	commentClassificationAllColumns            = []string{"comment_id", "toxicity", "severe_toxicity", "obscene", "identity_attack", "insult", "threat", "sexual_explicit", "nazi", "doxx", "is_reviewed", "reviewer_approved", "timestamp", "created_at", "updated_at", "model_ident"}
	commentClassificationColumnsWithoutDefault = []string{"comment_id", "toxicity", "severe_toxicity", "obscene", "identity_attack", "insult", "threat", "sexual_explicit", "nazi", "doxx", "reviewer_approved", "timestamp", "model_ident"}
	commentClassificationColumnsWithDefault    = []string{"is_reviewed", "created_at", "updated_at"}
	commentClassificationPrimaryKeyColumns     = []string{"comment_id"}
	commentClassificationGeneratedColumns      = []string{}
)

type (
	// CommentClassificationSlice is an alias for a slice of pointers to CommentClassification.
	// This should almost always be used instead of []CommentClassification.
	CommentClassificationSlice []*CommentClassification

	commentClassificationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	commentClassificationType                 = reflect.TypeOf(&CommentClassification{})
	commentClassificationMapping              = queries.MakeStructMapping(commentClassificationType)
	commentClassificationPrimaryKeyMapping, _ = queries.BindMapping(commentClassificationType, commentClassificationMapping, commentClassificationPrimaryKeyColumns)
	commentClassificationInsertCacheMut       sync.RWMutex
	commentClassificationInsertCache          = make(map[string]insertCache)
	commentClassificationUpdateCacheMut       sync.RWMutex
	commentClassificationUpdateCache          = make(map[string]updateCache)
	commentClassificationUpsertCacheMut       sync.RWMutex
	commentClassificationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single commentClassification record from the query.
func (q commentClassificationQuery) One(exec boil.Executor) (*CommentClassification, error) {
	o := &CommentClassification{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for comment_classification")
	}

	return o, nil
}

// All returns all CommentClassification records from the query.
func (q commentClassificationQuery) All(exec boil.Executor) (CommentClassificationSlice, error) {
	var o []*CommentClassification

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to CommentClassification slice")
	}

	return o, nil
}

// Count returns the count of all CommentClassification records in the query.
func (q commentClassificationQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count comment_classification rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q commentClassificationQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if comment_classification exists")
	}

	return count > 0, nil
}

// Comment pointed to by the foreign key.
func (o *CommentClassification) Comment(mods ...qm.QueryMod) commentQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`comment_id` = ?", o.CommentID),
	}

	queryMods = append(queryMods, mods...)

	return Comments(queryMods...)
}

// LoadComment allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (commentClassificationL) LoadComment(e boil.Executor, singular bool, maybeCommentClassification interface{}, mods queries.Applicator) error {
	var slice []*CommentClassification
	var object *CommentClassification

	if singular {
		var ok bool
		object, ok = maybeCommentClassification.(*CommentClassification)
		if !ok {
			object = new(CommentClassification)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeCommentClassification)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeCommentClassification))
			}
		}
	} else {
		s, ok := maybeCommentClassification.(*[]*CommentClassification)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeCommentClassification)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeCommentClassification))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &commentClassificationR{}
		}
		args = append(args, object.CommentID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &commentClassificationR{}
			}

			for _, a := range args {
				if a == obj.CommentID {
					continue Outer
				}
			}

			args = append(args, obj.CommentID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`comment`),
		qm.WhereIn(`comment.comment_id in ?`, args...),
		qmhelper.WhereIsNull(`comment.deleted_at`),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Comment")
	}

	var resultSlice []*Comment
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Comment")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for comment")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for comment")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Comment = foreign
		if foreign.R == nil {
			foreign.R = &commentR{}
		}
		foreign.R.CommentClassification = object
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.CommentID == foreign.CommentID {
				local.R.Comment = foreign
				if foreign.R == nil {
					foreign.R = &commentR{}
				}
				foreign.R.CommentClassification = local
				break
			}
		}
	}

	return nil
}

// SetComment of the commentClassification to the related item.
// Sets o.R.Comment to related.
// Adds o to related.R.CommentClassification.
func (o *CommentClassification) SetComment(exec boil.Executor, insert bool, related *Comment) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `comment_classification` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"comment_id"}),
		strmangle.WhereClause("`", "`", 0, commentClassificationPrimaryKeyColumns),
	)
	values := []interface{}{related.CommentID, o.CommentID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CommentID = related.CommentID
	if o.R == nil {
		o.R = &commentClassificationR{
			Comment: related,
		}
	} else {
		o.R.Comment = related
	}

	if related.R == nil {
		related.R = &commentR{
			CommentClassification: o,
		}
	} else {
		related.R.CommentClassification = o
	}

	return nil
}

// CommentClassifications retrieves all the records using an executor.
func CommentClassifications(mods ...qm.QueryMod) commentClassificationQuery {
	mods = append(mods, qm.From("`comment_classification`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`comment_classification`.*"})
	}

	return commentClassificationQuery{q}
}

// FindCommentClassification retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCommentClassification(exec boil.Executor, commentID string, selectCols ...string) (*CommentClassification, error) {
	commentClassificationObj := &CommentClassification{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `comment_classification` where `comment_id`=?", sel,
	)

	q := queries.Raw(query, commentID)

	err := q.Bind(nil, exec, commentClassificationObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from comment_classification")
	}

	return commentClassificationObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *CommentClassification) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no comment_classification provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(commentClassificationColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	commentClassificationInsertCacheMut.RLock()
	cache, cached := commentClassificationInsertCache[key]
	commentClassificationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			commentClassificationAllColumns,
			commentClassificationColumnsWithDefault,
			commentClassificationColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(commentClassificationType, commentClassificationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(commentClassificationType, commentClassificationMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `comment_classification` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `comment_classification` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `comment_classification` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, commentClassificationPrimaryKeyColumns))
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
		return errors.Wrap(err, "model: unable to insert into comment_classification")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.CommentID,
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, identifierCols...)
	}
	err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "model: unable to populate default values for comment_classification")
	}

CacheNoHooks:
	if !cached {
		commentClassificationInsertCacheMut.Lock()
		commentClassificationInsertCache[key] = cache
		commentClassificationInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the CommentClassification.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *CommentClassification) Update(exec boil.Executor, columns boil.Columns) error {
	var err error
	key := makeCacheKey(columns, nil)
	commentClassificationUpdateCacheMut.RLock()
	cache, cached := commentClassificationUpdateCache[key]
	commentClassificationUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			commentClassificationAllColumns,
			commentClassificationPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return errors.New("model: unable to update comment_classification, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `comment_classification` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, commentClassificationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(commentClassificationType, commentClassificationMapping, append(wl, commentClassificationPrimaryKeyColumns...))
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
		return errors.Wrap(err, "model: unable to update comment_classification row")
	}

	if !cached {
		commentClassificationUpdateCacheMut.Lock()
		commentClassificationUpdateCache[key] = cache
		commentClassificationUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAll updates all rows with the specified column values.
func (q commentClassificationQuery) UpdateAll(exec boil.Executor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "model: unable to update all for comment_classification")
	}

	return nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CommentClassificationSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), commentClassificationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `comment_classification` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, commentClassificationPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to update all in commentClassification slice")
	}

	return nil
}

var mySQLCommentClassificationUniqueColumns = []string{
	"comment_id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *CommentClassification) Upsert(exec boil.Executor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model: no comment_classification provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(commentClassificationColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLCommentClassificationUniqueColumns, o)

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

	commentClassificationUpsertCacheMut.RLock()
	cache, cached := commentClassificationUpsertCache[key]
	commentClassificationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			commentClassificationAllColumns,
			commentClassificationColumnsWithDefault,
			commentClassificationColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			commentClassificationAllColumns,
			commentClassificationPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("model: unable to upsert comment_classification, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`comment_classification`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `comment_classification` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(commentClassificationType, commentClassificationMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(commentClassificationType, commentClassificationMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert for comment_classification")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(commentClassificationType, commentClassificationMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model: unable to retrieve unique values for comment_classification")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, nzUniqueCols...)
	}
	err = exec.QueryRow(cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model: unable to populate default values for comment_classification")
	}

CacheNoHooks:
	if !cached {
		commentClassificationUpsertCacheMut.Lock()
		commentClassificationUpsertCache[key] = cache
		commentClassificationUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single CommentClassification record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *CommentClassification) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("model: no CommentClassification provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), commentClassificationPrimaryKeyMapping)
	sql := "DELETE FROM `comment_classification` WHERE `comment_id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete from comment_classification")
	}

	return nil
}

// DeleteAll deletes all matching rows.
func (q commentClassificationQuery) DeleteAll(exec boil.Executor) error {
	if q.Query == nil {
		return errors.New("model: no commentClassificationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete all from comment_classification")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CommentClassificationSlice) DeleteAll(exec boil.Executor) error {
	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), commentClassificationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `comment_classification` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, commentClassificationPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete all from commentClassification slice")
	}

	return nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *CommentClassification) Reload(exec boil.Executor) error {
	ret, err := FindCommentClassification(exec, o.CommentID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CommentClassificationSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CommentClassificationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), commentClassificationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `comment_classification`.* FROM `comment_classification` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, commentClassificationPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in CommentClassificationSlice")
	}

	*o = slice

	return nil
}

// CommentClassificationExists checks if the CommentClassification row exists.
func CommentClassificationExists(exec boil.Executor, commentID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `comment_classification` where `comment_id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, commentID)
	}
	row := exec.QueryRow(sql, commentID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if comment_classification exists")
	}

	return exists, nil
}

// Exists checks if the CommentClassification row exists.
func (o *CommentClassification) Exists(exec boil.Executor) (bool, error) {
	return CommentClassificationExists(exec, o.CommentID)
}
