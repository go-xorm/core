package core

import (
	"fmt"
	"strings"
	"time"
)

type DbType string

type Uri struct {
	DbType  DbType
	Proto   string
	Host    string
	Port    string
	DbName  string
	User    string
	Passwd  string
	Charset string
	Laddr   string
	Raddr   string
	Timeout time.Duration
}

// a dialect is a driver's wrapper
type Dialect interface {
	Init(*DB, *Uri, string, string) error
	URI() *Uri
	DB() *DB
	DBType() DbType
	SqlType(*Column) string
	FormatBytes(b []byte) string

	DriverName() string
	DataSourceName() string

	QuoteStr() string
	AndStr() string
	OrStr() string
	EqStr() string
	RollBackStr() string
	AutoIncrStr() string

	SupportInsertMany() bool
	SupportEngine() bool
	SupportCharset() bool
	IndexOnTable() bool
	ShowCreateNull() bool

	IndexCheckSql(tableName, idxName string) (string, []interface{})
	TableCheckSql(tableName string) (string, []interface{})
	//ColumnCheckSql(tableName, colName string) (string, []interface{})

	//IsTableExist(tableName string) (bool, error)
	//IsIndexExist(tableName string, idx *Index) (bool, error)
	IsColumnExist(tableName string, col *Column) (bool, error)

	CreateTableSql(table *Table, tableName, storeEngine, charset string) string
	DropTableSql(tableName string) string
	CreateIndexSql(tableName string, index *Index) string

	GetColumns(tableName string) ([]string, map[string]*Column, error)
	GetTables() ([]*Table, error)
	GetIndexes(tableName string) (map[string]*Index, error)

	// Get data from db cell to a struct's field
	//GetData(col *Column, fieldValue *reflect.Value, cellData interface{}) error
	// Set field data to db
	//SetData(col *Column, fieldValue *refelct.Value) (interface{}, error)

	Filters() []Filter
}

func OpenDialect(dialect Dialect) (*DB, error) {
	return Open(dialect.DriverName(), dialect.DataSourceName())
}

type Base struct {
	db             *DB
	dialect        Dialect
	driverName     string
	dataSourceName string
	*Uri
}

func (b *Base) DB() *DB {
	return b.db
}

func (b *Base) Init(db *DB, dialect Dialect, uri *Uri, drivername, dataSourceName string) error {
	b.db, b.dialect, b.Uri = db, dialect, uri
	b.driverName, b.dataSourceName = drivername, dataSourceName
	return nil
}

func (b *Base) URI() *Uri {
	return b.Uri
}

func (b *Base) DBType() DbType {
	return b.Uri.DbType
}

func (b *Base) FormatBytes(bs []byte) string {
	return fmt.Sprintf("0x%x", bs)
}

func (b *Base) DriverName() string {
	return b.driverName
}

func (b *Base) ShowCreateNull() bool {
	return true
}

func (b *Base) DataSourceName() string {
	return b.dataSourceName
}

func (b *Base) Quote(c string) string {
	return b.dialect.QuoteStr() + c + b.dialect.QuoteStr()
}

func (b *Base) AndStr() string {
	return "AND"
}

func (b *Base) OrStr() string {
	return "OR"
}

func (b *Base) EqStr() string {
	return "="
}

func (db *Base) RollBackStr() string {
	return "ROLL BACK"
}

func (db *Base) DropTableSql(tableName string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS `%s`", tableName)
}

func (db *Base) HasRecords(query string, args ...interface{}) (bool, error) {
	rows, err := db.DB().Query(query, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (db *Base) IsColumnExist(tableName string, col *Column) (bool, error) {
	query := "SELECT `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? AND `COLUMN_NAME` = ?"
	query = strings.Replace(query, "`", db.dialect.QuoteStr(), -1)
	return db.HasRecords(query, db.DbName, tableName, col.Name)
}

func (db *Base) CreateIndexSql(tableName string, index *Index) string {
	quote := db.Quote
	var unique string
	var idxName string
	if index.Type == UniqueType {
		unique = " UNIQUE"
		idxName = fmt.Sprintf("UQE_%v_%v", tableName, index.Name)
	} else {
		idxName = fmt.Sprintf("IDX_%v_%v", tableName, index.Name)
	}
	return fmt.Sprintf("CREATE%s INDEX %v ON %v (%v);", unique,
		quote(idxName), quote(tableName),
		quote(strings.Join(index.Cols, quote(","))))
}

func (b *Base) CreateTableSql(table *Table, tableName, storeEngine, charset string) string {
	var sql string
	sql = "CREATE TABLE IF NOT EXISTS "
	if tableName == "" {
		tableName = table.Name
	}

	sql += b.Quote(tableName) + " ("

	pkList := table.PrimaryKeys

	for _, colName := range table.ColumnsSeq() {
		col := table.GetColumn(colName)
		if col.IsPrimaryKey && len(pkList) == 1 {
			sql += col.String(b.dialect)
		} else {
			sql += col.StringNoPk(b.dialect)
		}
		sql = strings.TrimSpace(sql)
		sql += ", "
	}

	if len(pkList) > 1 {
		sql += "PRIMARY KEY ( "
		sql += b.Quote(strings.Join(pkList, b.Quote(",")))
		sql += " ), "
	}

	sql = sql[:len(sql)-2] + ")"
	if b.dialect.SupportEngine() && storeEngine != "" {
		sql += " ENGINE=" + storeEngine
	}
	if b.dialect.SupportCharset() {
		if len(charset) == 0 {
			charset = b.dialect.URI().Charset
		}
		if len(charset) > 0 {
			sql += " DEFAULT CHARSET " + charset
		}
	}
	sql += ";"
	return sql
}

var (
	dialects = map[DbType]Dialect{}
)

func RegisterDialect(dbName DbType, dialect Dialect) {
	if dialect == nil {
		panic("core: Register dialect is nil")
	}
	dialects[dbName] = dialect // !nashtsai! allow override dialect
}

func QueryDialect(dbName DbType) Dialect {
	return dialects[dbName]
}
