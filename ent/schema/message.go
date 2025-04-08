package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("name"),
		field.String("email"),
		field.String("phone").Optional(),
		field.Enum("subject").Values("GENERAL_ENQUIRY", "PRAYER_REQUEST", "MINISTRY_QUESTION", "EVENT_INFORMATION", "OTHER").Default("GENERAL_ENQUIRY").SchemaType(map[string]string{
			dialect.Postgres: "varchar",
		}),
		field.Enum("state").Values("READ", "UNREAD", "RESPONDED"),
		field.String("description"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("responses", Response.Type),
	}
}
