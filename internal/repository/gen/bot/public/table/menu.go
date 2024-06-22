//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Menu = newMenuTable("public", "menu", "")

type menuTable struct {
	postgres.Table

	// Columns
	Alcohol postgres.ColumnBool
	Photo   postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type MenuTable struct {
	menuTable

	EXCLUDED menuTable
}

// AS creates new MenuTable with assigned alias
func (a MenuTable) AS(alias string) *MenuTable {
	return newMenuTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new MenuTable with assigned schema name
func (a MenuTable) FromSchema(schemaName string) *MenuTable {
	return newMenuTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new MenuTable with assigned table prefix
func (a MenuTable) WithPrefix(prefix string) *MenuTable {
	return newMenuTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new MenuTable with assigned table suffix
func (a MenuTable) WithSuffix(suffix string) *MenuTable {
	return newMenuTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newMenuTable(schemaName, tableName, alias string) *MenuTable {
	return &MenuTable{
		menuTable: newMenuTableImpl(schemaName, tableName, alias),
		EXCLUDED:  newMenuTableImpl("", "excluded", ""),
	}
}

func newMenuTableImpl(schemaName, tableName, alias string) menuTable {
	var (
		AlcoholColumn  = postgres.BoolColumn("alcohol")
		PhotoColumn    = postgres.StringColumn("photo")
		allColumns     = postgres.ColumnList{AlcoholColumn, PhotoColumn}
		mutableColumns = postgres.ColumnList{AlcoholColumn, PhotoColumn}
	)

	return menuTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		Alcohol: AlcoholColumn,
		Photo:   PhotoColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
