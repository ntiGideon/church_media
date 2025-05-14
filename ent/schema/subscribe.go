package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// Subscribe holds the schema definition for the Subscribe entity.
type Subscribe struct {
	ent.Schema
}

// Fields of the Subscribe.
func (Subscribe) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("email").NotEmpty().Unique(),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Subscribe.
func (Subscribe) Edges() []ent.Edge {
	return nil
}
