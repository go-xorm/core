// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"database/sql"
	"io"
	"reflect"
)

// Engine is the major struct of xorm, it means a database manager.
// Commonly, an application only need one engine
type Engine interface {
	// ShowSQL show SQL statment or not on logger if log level is great than INFO
	ShowSQL(show ...bool)

	// ShowExecTime show SQL statment and execute time or not on logger if log level is great than INFO
	ShowExecTime(show ...bool)

	// Logger return the logger interface
	Logger() ILogger

	// SetLogger set the new logger
	SetLogger(logger ILogger)

	// SetDisableGlobalCache disable global cache or not
	SetDisableGlobalCache(disable bool)

	// DriverName return the current sql driver's name
	DriverName() string

	// DataSourceName return the current connection string
	DataSourceName() string

	// SetMapper set the name mapping rules
	SetMapper(mapper IMapper)

	// SetTableMapper set the table name mapping rule
	SetTableMapper(mapper IMapper)

	// SetColumnMapper set the column name mapping rule
	SetColumnMapper(mapper IMapper)

	// SupportInsertMany If engine's database support batch insert records like
	// "insert into user values (name, age), (name, age)".
	// When the return is ture, then engine.Insert(&users) will
	// generate batch sql and exeute.
	SupportInsertMany() bool

	// QuoteStr Engine's database use which charactor as quote.
	// mysql, sqlite use ` and postgres use "
	QuoteStr() string

	// Quote Use QuoteStr quote the string sql
	Quote(value string) string

	// QuoteTo quotes string and writes into the buffer
	QuoteTo(buf *bytes.Buffer, value string)

	// SqlType will be depracated, please use SQLType instead
	//
	// Deprecated: use SQLType instead
	SqlType(c *Column)

	// SQLType A simple wrapper to dialect's SqlType method
	SQLType(c *Column) string

	// AutoIncrStr Database's autoincrement statement
	AutoIncrStr() string

	// SetMaxOpenConns is only available for go 1.2+
	SetMaxOpenConns(conns int)

	// SetMaxIdleConns set the max idle connections on pool, default is 2
	SetMaxIdleConns(conns int)

	// SetDefaultCacher set the default cacher. Xorm's default not enable cacher.
	SetDefaultCacher(cacher Cacher)

	// NoCache If you has set default cacher, and you want temporilly stop use cache,
	// you can use NoCache()
	NoCache() Session

	// NoCascade If you do not want to auto cascade load object
	NoCascade() Session

	// MapCacher Set a table use a special cacher
	MapCacher(bean interface{}, cacher Cacher)

	// NewDB provides an interface to operate database directly
	NewDB() (*DB, error)

	// DB return the wrapper of sql.DB
	DB() *DB

	// Dialect return database dialect
	Dialect() Dialect

	// NewSession New a session
	NewSession() Session

	// Close the engine
	Close() error

	// Ping tests if database is alive
	Ping() error

	// Sql provides raw sql input parameter. When you have a complex SQL statement
	// and cannot use Where, Id, In and etc. Methods to describe, you can use SQL.
	//
	// Deprecated: use SQL instead.
	Sql(querystring string, args ...interface{}) Session

	// SQL method let's you manualy write raw SQL and operate
	// For example:
	//
	//         engine.SQL("select * from user").Find(&users)
	//
	// This    code will execute "select * from user" and set the records to users
	SQL(query interface{}, args ...interface{}) Session

	// NoAutoTime Default if your struct has "created" or "updated" filed tag, the fields
	// will automatically be filled with current time when Insert or Update
	// invoked. Call NoAutoTime if you dont' want to fill automatically.
	NoAutoTime() Session

	// NoAutoCondition disable auto generate Where condition from bean or not
	NoAutoCondition(no ...bool) Session

	// DBMetas Retrieve all tables, columns, indexes' informations from database.
	DBMetas() ([]*Table, error)

	// DumpAllToFile dump database all table structs and data to a file
	DumpAllToFile(fp string) error

	// DumpAll dump database all table structs and data to w
	DumpAll(w io.Writer) error

	// DumpTablesToFile dump specified tables to SQL file.
	DumpTablesToFile(tables []*Table, fp string, tp ...DbType) error

	// DumpTables dump specify tables to io.Writer
	DumpTables(tables []*Table, w io.Writer, tp ...DbType) error

	// Cascade use cascade or not
	Cascade(trueOrFalse ...bool) Session

	// Where method provide a condition query
	Where(query interface{}, args ...interface{}) Session

	// Id will be depracated, please use ID instead
	Id(id interface{}) Session

	// ID mehtod provoide a condition as (id) = ?
	ID(id interface{}) Session

	// Before apply before Processor, affected bean is passed to closure arg
	Before(closures func(interface{})) Session

	// After apply after insert Processor, affected bean is passed to closure arg
	After(closures func(interface{})) Session

	// Charset set charset when create table, only support mysql now
	Charset(charset string) Session

	// StoreEngine set store engine when create table, only support mysql now
	StoreEngine(storeEngine string) Session

	// Distinct use for distinct columns. Caution: when you are using cache,
	// distinct will not be cached because cache system need id,
	// but distinct will not provide id
	Distinct(columns ...string) Session

	// Select customerize your select columns or contents
	Select(str string) Session

	// Cols only use the paramters as select or update columns
	Cols(columns ...string) Session

	// AllCols indicates that all columns should be use
	AllCols() Session

	// MustCols specify some columns must use even if they are empty
	MustCols(columns ...string) Session

	// UseBool xorm automatically retrieve condition according struct, but
	// if struct has bool field, it will ignore them. So use UseBool
	// to tell system to do not ignore them.
	// If no paramters, it will use all the bool field of struct, or
	// it will use paramters's columns
	UseBool(columns ...string) Session

	// Omit only not use the paramters as select or update columns
	Omit(columns ...string) Session

	// Nullable set null when column is zero-value and nullable for update
	Nullable(columns ...string) Session

	// In will generate "column IN (?, ?)"
	In(column string, args ...interface{}) Session

	// Incr provides a update string like "column = column + ?"
	Incr(column string, arg ...interface{}) Session

	// Decr provides a update string like "column = column - ?"
	Decr(column string, arg ...interface{}) Session

	// SetExpr provides a update string like "column = {expression}"
	SetExpr(column string, expression string) Session

	// Table temporarily change the Get, Find, Update's table
	Table(tableNameOrBean interface{}) Session

	// Alias set the table alias
	Alias(alias string) Session

	// Limit will generate "LIMIT start, limit"
	Limit(limit int, start ...int) Session

	// Desc will generate "ORDER BY column1 DESC, column2 DESC"
	Desc(colNames ...string) Session

	// Asc will generate "ORDER BY column1,column2 Asc"
	// This method can chainable use.
	//
	//        engine.Desc("name").Asc("age").Find(&users)
	//        // SELECT * FROM user ORDER BY name DESC, age ASC
	//
	Asc(colNames ...string) Session

	// OrderBy will generate "ORDER BY order"
	OrderBy(order string) Session

	// Join the join_operator should be one of INNER, LEFT OUTER, CROSS etc - this will be prepended to JOIN
	Join(joinOperator string, tablename interface{}, condition string, args ...interface{}) Session

	// GroupBy generate group by statement
	GroupBy(keys string) Session

	// Having generate having statement
	Having(conditions string) Session

	// IsTableEmpty if a table has any reocrd
	IsTableEmpty(bean interface{}) (bool, error)

	// IsTableExist if a table is exist
	IsTableExist(beanOrTableName interface{}) (bool, error)

	// IDOf get id from one struct
	IDOf(bean interface{}) PK

	// IDOfV get id from one value of struct
	IDOfV(rv reflect.Value) PK

	// CreateIndexes create indexes
	CreateIndexes(bean interface{}) error

	// CreateUniques create uniques
	CreateUniques(bean interface{}) error

	// ClearCacheBean if enabled cache, clear the cache bean
	ClearCacheBean(bean interface{}, id string) error

	// ClearCache if enabled cache, clear some tables' cache
	ClearCache(beans ...interface{}) error

	// CreateTables create tabls according bean
	CreateTables(beans ...interface{}) error

	// DropTables drop specify tables
	DropTables(beans ...interface{}) error

	// Exec raw sql
	Exec(sql string, args ...interface{}) (sql.Result, error)

	// Query a raw sql and return records as []map[string][]byte
	Query(sql string, paramStr ...interface{}) (resultsSlice []map[string][]byte, err error)

	// Insert one or more records
	Insert(beans ...interface{}) (int64, error)

	// InsertOne insert only one record
	InsertOne(bean interface{}) (int64, error)

	// Update records, bean's non-empty fields are updated contents,
	// condiBean' non-empty filds are conditions
	// CAUTION:
	//        1.bool will defaultly be updated content nor conditions
	//         You should call UseBool if you have bool to use.
	//        2.float32 & float64 may be not inexact as conditions
	Update(bean interface{}, condiBeans ...interface{}) (int64, error)

	// Delete records, bean's non-empty fields are conditions
	Delete(bean interface{}) (int64, error)

	// Get retrieve one record from table, bean's non-empty fields
	// are conditions
	Get(bean interface{}) (bool, error)

	// Find retrieve records from table, condiBeans's non-empty fields
	// are conditions. beans could be []Struct, []*Struct, map[int64]Struct
	// map[int64]*Struct
	Find(beans interface{}, condiBeans ...interface{}) error

	Iterate(bean interface{}, fun IterFunc) error

	// Rows return sql.Rows compatible Rows obj, as a forward Iterator object for iterating record by record, bean's non-empty fields
	// are conditions.
	Rows(bean interface{}) (*Rows, error)

	// Count counts the records. bean's non-empty fields are conditions.
	Count(bean interface{}) (int64, error)

	// Sum sum the records by some column. bean's non-empty fields are conditions.
	Sum(bean interface{}, colName string) (float64, error)

	// Sums sum the records by some columns. bean's non-empty fields are conditions.
	Sums(bean interface{}, colNames ...string) ([]float64, error)

	// SumsInt like Sums but return slice of int64 instead of float64.
	SumsInt(bean interface{}, colNames ...string) ([]int64, error)

	// Unscoped always disable struct tag "deleted"
	Unscoped() Session
}
