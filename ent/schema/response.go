package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// Response holds the schema definition for the Response entity.
type Response struct {
	ent.Schema
}

// Fields of the Response.
func (Response) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("email"),
		field.Enum("subject").Values("GENERAL_ENQUIRY", "PRAYER_REQUEST", "MINISTRY_QUESTION", "EVENT_INFORMATION", "OTHER").Default("GENERAL_ENQUIRY").SchemaType(map[string]string{
			dialect.Postgres: "varchar",
		}),
		field.Int("user_id"),
		field.String("description"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Response.
func (Response) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("message", Message.Type).Ref("responses").Unique().Required().Field("user_id"),
		edge.From("user", User.Type).Ref("responses").Unique().Required().Field("user_id"),
	}
}
