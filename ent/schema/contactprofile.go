package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"fmt"
	"regexp"
)

// ContactProfile holds the schema definition for the ContactProfile entity.
type ContactProfile struct {
	ent.Schema
}

// Fields of the ContactProfile.
func (ContactProfile) Fields() []ent.Field {
	return []ent.Field{
		field.String("first_name").
			NotEmpty().
			MaxLen(50),

		field.String("surname").
			NotEmpty().
			MaxLen(50),

		field.String("phone_number").
			Unique().
			Optional().
			Nillable().
			Validate(phoneValidator),

		field.String("profile_picture").
			Optional().
			Nillable(),

		field.Text("address").
			Optional().
			Nillable(),

		field.Time("date_of_birth").
			Optional().
			Nillable(),

		field.Enum("gender").
			Values("male", "female").
			Optional(),

		field.String("occupation").
			Optional().
			MaxLen(100),

		field.String("marital_status").
			Optional().
			MaxLen(20),
	}
}

// Edges of the ContactProfile.
func (ContactProfile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("contact_profile").
			Unique().
			Required(),
	}
}

func phoneValidator(phone string) error {
	if phone == "" {
		return nil
	}
	matched, _ := regexp.MatchString(`^\+?[0-9]{10,15}$`, phone)
	if !matched {
		return fmt.Errorf("invalid phone number format")
	}
	return nil
}
