package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"fmt"
	"net/mail"
	"regexp"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			Unique().
			Optional().
			NotEmpty().
			MinLen(3).
			MaxLen(20).
			Match(usernameRegex()),

		field.String("verify_token").
			Unique().
			Optional(),

		field.String("email").
			Unique().
			NotEmpty().
			Validate(emailValidator),

		field.String("password").
			Optional().
			Sensitive().
			NotEmpty().
			MinLen(8),

		field.String("registration_token").
			Unique().
			Optional().
			Nillable(),

		field.Enum("role").
			Values("member", "deacon", "pastor", "admin", "content_admin", "secretary", "superadmin").
			Default("member"),

		field.Time("token_expires_at").
			Optional().
			Nillable(),
		field.Enum("state").Values("FRESH", "DELETED", "VERIFIED", "ACCEPTED", "PENDING", "EXPIRED").Default("FRESH").SchemaType(map[string]string{
			dialect.Postgres: "varchar",
		}),
		field.Time("created_at").
			Default(time.Now),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("responses", Response.Type),
		edge.To("contact_profile", ContactProfile.Type),
	}
}

// Indexes set index on email field to improve speed of data retrieval
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username"),
		index.Fields("email"),
		index.Fields("registration_token"),
		index.Fields("role"),
	}
}

// Mixin for the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func usernameRegex() *regexp.Regexp {
	return regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
}

func emailValidator(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email format")
	}
	return nil
}
