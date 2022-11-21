// Code generated by ent, DO NOT EDIT.

package pagesource

const (
	// Label holds the string label denoting the pagesource type in the database.
	Label = "page_source"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldSource holds the string denoting the source field in the database.
	FieldSource = "source"
	// Table holds the table name of the pagesource in the database.
	Table = "page_source"
)

// Columns holds all SQL columns for pagesource fields.
var Columns = []string{
	FieldID,
	FieldSource,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "page_source"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"page_page_source",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultSource holds the default value on creation for the "source" field.
	DefaultSource string
)
