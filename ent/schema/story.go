package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// Story holds the schema definition for the Story entity.
type Story struct {
	ent.Schema
}

// Fields of the Story.
func (Story) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("title").
			MaxLen(100).
			NotEmpty(),
		field.String("body").
			NotEmpty(),
		field.String("image").
			Optional(),
		field.String("excerpt").
			MaxLen(500).
			Optional(),
		field.Enum("status").
			Values("draft", "published", "archived").
			Default("draft"),
		field.Time("published_at").
			Optional(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Int("author_id"),
	}
}

// Edges of the Story.
func (Story) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", User.Type).
			Ref("stories").
			Field("author_id").
			Unique().
			Required(),
	}
}
