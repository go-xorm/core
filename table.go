package core

import (
	"reflect"
	"strings"
)

// database table
type Table struct {
	Name          string
	Type          reflect.Type
	columnsSeq    []string
	columnsMap    map[string][]*Column
	columns       []*Column
	Indexes       map[string]*Index
	PrimaryKeys   []string
	AutoIncrement string
	Created       map[string]bool
	Updated       string
	Deleted       string
	Version       string
	Cacher        Cacher
	StoreEngine   string
	Charset       string

	Inherits []string //By Hzm
}

func (table *Table) Columns() []*Column {
	return table.columns
}

func (table *Table) ColumnsSeq() []string {
	return table.columnsSeq
}

func NewEmptyTable() *Table {
	return NewTable("", nil)
}

func NewTable(name string, t reflect.Type) *Table {
	return &Table{Name: name, Type: t,
		columnsSeq:  make([]string, 0),
		columns:     make([]*Column, 0),
		columnsMap:  make(map[string][]*Column),
		Indexes:     make(map[string]*Index),
		Created:     make(map[string]bool),
		PrimaryKeys: make([]string, 0),
		Inherits:    make([]string, 0), //by hzm
	}
}

func (table *Table) GetColumn(name string) *Column {
	if c, ok := table.columnsMap[strings.ToLower(name)]; ok {
		return c[0]
	}
	return nil
}

func (table *Table) GetColumnIdx(name string, idx int) *Column {
	if c, ok := table.columnsMap[strings.ToLower(name)]; ok {
		if idx < len(c) {
			return c[idx]
		}
	}
	return nil
}

// if has primary key, return column
func (table *Table) PKColumns() []*Column {
	columns := make([]*Column, len(table.PrimaryKeys))
	for i, name := range table.PrimaryKeys {
		columns[i] = table.GetColumn(name)
	}
	return columns
}

func (table *Table) ColumnType(name string) reflect.Type {
	t, _ := table.Type.FieldByName(name)
	return t.Type
}

func (table *Table) AutoIncrColumn() *Column {
	return table.GetColumn(table.AutoIncrement)
}

func (table *Table) VersionColumn() *Column {
	return table.GetColumn(table.Version)
}

func (table *Table) UpdatedColumn() *Column {
	return table.GetColumn(table.Updated)
}

func (table *Table) DeletedColumn() *Column {
	return table.GetColumn(table.Deleted)
}

// add a column to table
func (table *Table) AddColumn(col *Column) {
	table.columnsSeq = append(table.columnsSeq, col.Name)
	table.columns = append(table.columns, col)
	colName := strings.ToLower(col.Name)
	if c, ok := table.columnsMap[colName]; ok {
		table.columnsMap[colName] = append(c, col)
	} else {
		table.columnsMap[colName] = []*Column{col}
	}

	if col.IsPrimaryKey {
		table.PrimaryKeys = append(table.PrimaryKeys, col.Name)
	}
	if col.IsAutoIncrement {
		table.AutoIncrement = col.Name
	}
	if col.IsCreated {
		table.Created[col.Name] = true
	}
	if col.IsUpdated {
		table.Updated = col.Name
	}
	if col.IsDeleted {
		table.Deleted = col.Name
	}
	if col.IsVersion {
		table.Version = col.Name
	}
}

// add an index or an unique to table
func (table *Table) AddIndex(index *Index) {
	table.Indexes[index.Name] = index
}

// by hzm
func (table *Table) AddInherit(aNew string) {
	table.Inherits = append(table.Inherits, aNew)
}
