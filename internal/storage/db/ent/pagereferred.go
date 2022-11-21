// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/drakejin/crawler/internal/storage/db/ent/pagereferred"
	"github.com/google/uuid"
)

// PageReferred is the model entity for the PageReferred schema.
type PageReferred struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// SourceID holds the value of the "source_id" field.
	SourceID uuid.UUID `json:"source_id,omitempty"`
	// TargetID holds the value of the "target_id" field.
	TargetID uuid.UUID `json:"target_id,omitempty"`
	// first indexed time
	CreatedAt time.Time `json:"created_at,omitempty"`
	// first indexed time by which system
	CreatedBy string `json:"created_by,omitempty"`
	// modified time
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// modified by which system
	UpdatedBy string `json:"updated_by,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PageReferred) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case pagereferred.FieldCreatedBy, pagereferred.FieldUpdatedBy:
			values[i] = new(sql.NullString)
		case pagereferred.FieldCreatedAt, pagereferred.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case pagereferred.FieldID, pagereferred.FieldSourceID, pagereferred.FieldTargetID:
			values[i] = new(uuid.UUID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type PageReferred", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PageReferred fields.
func (pr *PageReferred) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case pagereferred.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				pr.ID = *value
			}
		case pagereferred.FieldSourceID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field source_id", values[i])
			} else if value != nil {
				pr.SourceID = *value
			}
		case pagereferred.FieldTargetID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field target_id", values[i])
			} else if value != nil {
				pr.TargetID = *value
			}
		case pagereferred.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				pr.CreatedAt = value.Time
			}
		case pagereferred.FieldCreatedBy:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field created_by", values[i])
			} else if value.Valid {
				pr.CreatedBy = value.String
			}
		case pagereferred.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				pr.UpdatedAt = value.Time
			}
		case pagereferred.FieldUpdatedBy:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field updated_by", values[i])
			} else if value.Valid {
				pr.UpdatedBy = value.String
			}
		}
	}
	return nil
}

// Update returns a builder for updating this PageReferred.
// Note that you need to call PageReferred.Unwrap() before calling this method if this PageReferred
// was returned from a transaction, and the transaction was committed or rolled back.
func (pr *PageReferred) Update() *PageReferredUpdateOne {
	return (&PageReferredClient{config: pr.config}).UpdateOne(pr)
}

// Unwrap unwraps the PageReferred entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pr *PageReferred) Unwrap() *PageReferred {
	_tx, ok := pr.config.driver.(*txDriver)
	if !ok {
		panic("ent: PageReferred is not a transactional entity")
	}
	pr.config.driver = _tx.drv
	return pr
}

// String implements the fmt.Stringer.
func (pr *PageReferred) String() string {
	var builder strings.Builder
	builder.WriteString("PageReferred(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pr.ID))
	builder.WriteString("source_id=")
	builder.WriteString(fmt.Sprintf("%v", pr.SourceID))
	builder.WriteString(", ")
	builder.WriteString("target_id=")
	builder.WriteString(fmt.Sprintf("%v", pr.TargetID))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(pr.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("created_by=")
	builder.WriteString(pr.CreatedBy)
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(pr.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_by=")
	builder.WriteString(pr.UpdatedBy)
	builder.WriteByte(')')
	return builder.String()
}

// PageReferreds is a parsable slice of PageReferred.
type PageReferreds []*PageReferred

func (pr PageReferreds) config(cfg config) {
	for _i := range pr {
		pr[_i].config = cfg
	}
}
