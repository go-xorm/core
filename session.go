// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"database/sql"

	"github.com/go-xorm/builder"
)

// IterFunc only use by Iterate
type IterFunc func(idx int, bean interface{}) error

// Session keep a pointer to sql.DB and provides all execution of all
// kind of database operations.
type Session interface {
	// Clone copy all the session's content and return a new session
	Clone() Session

	// Init reset the session as the init status.
	Init()

	// Close release the connection from pool
	Close()

	// Prepare set a flag to session that should be prepare statment before execute query
	Prepare() Session

	// Sql provides raw sql input parameter. When you have a complex SQL statement
	// and cannot use Where, Id, In and etc. Methods to describe, you can use SQL.
	//
	// Deprecated: use SQL instead.
	Sql(query string, args ...interface{}) Session

	// SQL provides raw sql input parameter. When you have a complex SQL statement
	// and cannot use Where, Id, In and etc. Methods to describe, you can use SQL.
	SQL(query interface{}, args ...interface{}) Session

	// Where provides custom query condition.
	Where(query interface{}, args ...interface{}) Session

	// And provides custom query condition.
	And(query interface{}, args ...interface{}) Session

	// Or provides custom query condition.
	Or(query interface{}, args ...interface{}) Session

	// Id provides converting id as a query condition
	//
	// Deprecated: use ID instead
	Id(id interface{}) Session

	// ID provides converting id as a query condition
	ID(id interface{}) Session

	// Before Apply before Processor, affected bean is passed to closure arg
	Before(closures func(interface{})) Session

	// After Apply after Processor, affected bean is passed to closure arg
	After(closures func(interface{})) Session

	// Table can input a string or pointer to struct for special a table to operate.
	Table(tableNameOrBean interface{}) Session

	// Alias set the table alias
	Alias(alias string) Session

	// In provides a query string like "id in (1, 2, 3)"
	In(column string, args ...interface{}) Session

	// NotIn provides a query string like "id in (1, 2, 3)"
	NotIn(column string, args ...interface{}) Session

	// Incr provides a query string like "count = count + 1"
	Incr(column string, arg ...interface{}) Session

	// Decr provides a query string like "count = count - 1"
	Decr(column string, arg ...interface{}) Session

	// SetExpr provides a query string like "column =expression}"
	SetExpr(column string, expression string) Session

	// Select provides some columns to special
	Select(str string) Session

	// Cols provides some columns to special
	Cols(columns ...string) Session

	// AllCols ask all columns
	AllCols() Session

	// MustCols specify some columns must use even if they are empty
	MustCols(columns ...string) Session

	// NoCascade indicate that no cascade load child object
	NoCascade() Session

	// UseBool automatically retrieve condition according struct, but
	// if struct has bool field, it will ignore them. So use UseBool
	// to tell system to do not ignore them.
	// If no paramters, it will use all the bool field of struct, or
	// it will use paramters's columns
	UseBool(columns ...string) Session

	// Distinct use for distinct columns. Caution: when you are using cache,
	// distinct will not be cached because cache system need id,
	// but distinct will not provide id
	Distinct(columns ...string) Session

	// ForUpdate Set Read/Write locking for UPDATE
	ForUpdate() Session

	// Omit Only not use the paramters as select or update columns
	Omit(columns ...string) Session

	// Nullable Set null when column is zero-value and nullable for update
	Nullable(columns ...string) Session

	// NoAutoTime means do not automatically give created field and updated field
	// the current time on the current session temporarily
	NoAutoTime() Session

	// NoAutoCondition disable generate SQL condition from beans
	NoAutoCondition(no ...bool) Session

	// Limit provide limit and offset query condition
	Limit(limit int, start ...int) Session

	// OrderBy provide order by query condition, the input parameter is the content
	// after order by on a sql statement.
	OrderBy(order string) Session

	// Desc provide desc order by query condition, the input parameters are columns.
	Desc(colNames ...string) Session

	// Asc provide asc order by query condition, the input parameters are columns.
	Asc(colNames ...string) Session

	// StoreEngine is only avialble mysql dialect currently
	StoreEngine(storeEngine string) Session

	// Charset is only avialble mysql dialect currently
	Charset(charset string) Session

	// Cascade indicates if loading sub Struct
	Cascade(trueOrFalse ...bool) Session

	// NoCache ask this session do not retrieve data from cache system and
	// get data from database directly.
	NoCache() Session

	// Join join_operator should be one of INNER, LEFT OUTER, CROSS etc - this will be prepended to JOIN
	Join(joinOperator string, tablename interface{}, condition string, args ...interface{}) Session

	// GroupBy Generate Group By statement
	GroupBy(keys string) Session

	// Having Generate Having statement
	Having(conditions string) Session

	// DB db return the wrapper of sql.DB
	DB() *DB

	// Conds returns session query conditions
	Conds() builder.Cond

	// Begin a transaction
	Begin() error

	// Rollback When using transaction, you can rollback if any error
	Rollback() error

	// Commit When using transaction, Commit will commit all operations.
	Commit() error

	// Exec raw sql
	Exec(sqlStr string, args ...interface{}) (sql.Result, error)

	// CreateTable create a table according a bean
	CreateTable(bean interface{}) error

	// CreateIndexes create indexes
	CreateIndexes(bean interface{}) error

	// CreateUniques create uniques
	CreateUniques(bean interface{}) error

	// DropIndexes drop indexes
	DropIndexes(bean interface{}) error

	// DropTable drop table will drop table if exist, if drop failed, it will return error
	DropTable(beanOrTableName interface{}) error

	// Iterate record by record handle records from table, condiBeans's non-empty fields
	// are conditions. beans could be []Struct, []*Struct, map[int64]Struct
	// map[int64]*Struct
	Iterate(bean interface{}, fun IterFunc) error

	// Get retrieve one record from database, bean's non-empty fields
	// will be as conditions
	Get(bean interface{}) (bool, error)

	// Count counts the records. bean's non-empty fields
	// are conditions.
	Count(bean interface{}) (int64, error)

	// Sum call sum some column. bean's non-empty fields are conditions.
	Sum(bean interface{}, columnName string) (float64, error)

	// Sums call sum some columns. bean's non-empty fields are conditions.
	Sums(bean interface{}, columnNames ...string) ([]float64, error)

	// SumsInt sum specify columns and return as []int64 instead of []float64
	SumsInt(bean interface{}, columnNames ...string) ([]int64, error)

	// Find retrieve records from table, condiBeans's non-empty fields
	// are conditions. beans could be []Struct, []*Struct, map[int64]Struct
	// map[int64]*Struct
	Find(rowsSlicePtr interface{}, condiBean ...interface{}) error

	// Ping test if database is ok
	Ping() error

	// IsTableExist if a table is exist
	IsTableExist(beanOrTableName interface{}) (bool, error)

	// IsTableEmpty if table have any records
	IsTableEmpty(bean interface{}) (bool, error)

	// Query a raw sql and return records as []map[string][]byte
	Query(sqlStr string, paramStr ...interface{}) (resultsSlice []map[string][]byte, err error)

	// Insert insert one or more beans
	Insert(beans ...interface{}) (int64, error)

	// InsertMulti insert multiple records
	InsertMulti(rowsSlicePtr interface{}) (int64, error)

	// InsertOne insert only one struct into database as a record.
	// The in parameter bean must a struct or a point to struct. The return
	// parameter is inserted and error
	InsertOne(bean interface{}) (int64, error)

	// Update records, bean's non-empty fields are updated contents,
	// condiBean' non-empty filds are conditions
	// CAUTION:
	//        1.bool will defaultly be updated content nor conditions
	//         You should call UseBool if you have bool to use.
	//        2.float32 & float64 may be not inexact as conditions
	Update(bean interface{}, condiBean ...interface{}) (int64, error)

	// Delete records, bean's non-empty fields are conditions
	Delete(bean interface{}) (int64, error)

	// LastSQL returns last query information
	LastSQL() (string, []interface{})

	// Unscoped always disable struct tag "deleted"
	Unscoped() Session
}
