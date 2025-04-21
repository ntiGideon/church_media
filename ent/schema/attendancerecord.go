package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// AttendanceRecord holds the schema definition for the AttendanceRecord entity.
type AttendanceRecord struct {
	ent.Schema
}

// Fields of the AttendanceRecord.
func (AttendanceRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.Int("males"),
		field.Int("females"),
		field.Float("offering"),
		field.Float("tithe"),
	}
}

// Edges of the AttendanceRecord.
func (AttendanceRecord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("service", Service.Type).Ref("attendance_records").Unique().Required(),
	}
}
