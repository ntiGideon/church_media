package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// Event holds the schema definition for the Event entity.
type Event struct {
	ent.Schema
}

// Fields of the Event.
func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.Text("description"),
		field.Time("start_time"),
		field.Time("end_time"),
		field.String("location"),
		field.String("image_url").Optional(),
		field.Bool("featured").Default(false),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the Event.
func (Event) Edges() []ent.Edge {
	return nil
}
