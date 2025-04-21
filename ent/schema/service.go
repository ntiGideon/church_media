package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Service holds the schema definition for the Service entity.
type Service struct {
	ent.Schema
}

// Fields of the Service.
func (Service) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.Time("date"),
		field.String("name").Default("Sunday Service"),
		field.Enum("type").Values("first_service", "second_service", "friday_service", "wednesday_service", "children_service").StorageKey("type"),
	}
}

// Edges of the Service.
func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("attendance_records", AttendanceRecord.Type),
	}
}
