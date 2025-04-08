package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("name"),
		field.String("email"),
		field.String("password").Sensitive(),
		field.String("phone").Unique().Optional(),
		field.Enum("state").Values("FRESH", "DELETED", "VERIFIED").Default("FRESH").SchemaType(map[string]string{
			dialect.Postgres: "varchar",
		}),
		field.Time("created_at").Default(time.Now()),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("responses", Response.Type),
	}
}

// Indexes set index on email field to improve speed of data retrieval
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").Unique(),
	}
}
