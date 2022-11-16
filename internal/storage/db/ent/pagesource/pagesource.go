// Code generated by ent, DO NOT EDIT.

package pagesource

import (
	"entgo.io/ent"
)

const (
	// Label holds the string label denoting the pagesource type in the database.
	Label = "page_source"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldPageID holds the string denoting the page_id field in the database.
	FieldPageID = "page_id"
	// FieldReferredPageID holds the string denoting the referred_page_id field in the database.
	FieldReferredPageID = "referred_page_id"
	// FieldURL holds the string denoting the url field in the database.
	FieldURL = "url"
	// FieldReferredURL holds the string denoting the referred_url field in the database.
	FieldReferredURL = "referred_url"
	// FieldSource holds the string denoting the source field in the database.
	FieldSource = "source"
	// EdgePage holds the string denoting the page edge name in mutations.
	EdgePage = "page"
	// Table holds the table name of the pagesource in the database.
	Table = "page_source"
	// PageTable is the table that holds the page relation/edge.
	PageTable = "page_source"
	// PageInverseTable is the table name for the Page entity.
	// It exists in this package in order to avoid circular dependency with the "page" package.
	PageInverseTable = "page"
	// PageColumn is the table column denoting the page relation/edge.
	PageColumn = "page_id"
)

// Columns holds all SQL columns for pagesource fields.
var Columns = []string{
	FieldID,
	FieldPageID,
	FieldReferredPageID,
	FieldURL,
	FieldReferredURL,
	FieldSource,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/drakejin/crawler/internal/storage/db/ent/runtime"
var (
	Hooks [1]ent.Hook
	// DefaultURL holds the default value on creation for the "url" field.
	DefaultURL string
	// DefaultReferredURL holds the default value on creation for the "referred_url" field.
	DefaultReferredURL string
	// DefaultSource holds the default value on creation for the "source" field.
	DefaultSource string
)
