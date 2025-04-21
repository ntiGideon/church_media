package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// Member holds the schema definition for the Member entity.
type Member struct {
	ent.Schema
}

// Fields of the Member.
func (Member) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("form_number").Unique().Optional(),
		// Identification Fields
		field.String("id_number").Unique().
			Comment("National ID number"),
		field.String("surname").NotEmpty().
			Comment("Member's surname in block letters"),
		field.String("other_names").NotEmpty().
			Comment("Other names in block letters"),
		field.Time("dob").
			Comment("Date of birth"),
		field.Enum("gender").
			Values("male", "female", "other").
			Comment("Gender identity"),

		// Contact Information
		field.String("hometown").
			Comment("Home town"),
		field.String("region").
			Comment("Geographical region"),
		field.String("residence").Optional().
			Comment("Current residence address"),
		field.String("address").
			Comment("House number/digital address"),
		field.String("mobile").
			Comment("Mobile phone number"),
		field.String("email").Optional().
			Comment("Email address"),

		// Church Information
		field.String("sunday_school_class").Optional().
			Comment("Sunday school class name"),
		field.String("occupation").NotEmpty().
			Comment("Professional occupation"),
		field.Bool("has_title_card").Default(false).
			Comment("Whether member has title card"),
		field.String("title_card_number").Optional().
			Comment("Title card number if applicable"),
		field.String("day_born").Optional().
			Comment("Day of week born (e.g., Monday)"),

		// Spouse Information
		field.Bool("has_spouse").Default(false).
			Comment("Whether member has spouse"),
		field.String("spouse_id_number").Optional().
			Comment("Spouse's ID number if member"),
		field.String("spouse_name").Optional().
			Comment("Spouse's full name"),
		field.String("spouse_occupation").Optional().
			Comment("Spouse's occupation"),
		field.String("spouse_contact").Optional().
			Comment("Spouse's contact number"),

		// Baptism Information
		field.Bool("is_baptized").Default(false).
			Comment("Whether member is baptized"),
		field.String("baptized_by").Optional().
			Comment("Who performed the baptism"),
		field.String("baptism_church").Optional().
			Comment("Church where baptized"),
		field.String("baptism_cert_number").Optional().
			Comment("Baptism certificate number"),
		field.Time("baptism_date").Optional().
			Comment("Date of baptism"),
		field.Int("membership_year").Range(1900, time.Now().Year()).
			Comment("Year joined church membership"),

		// Photo Information
		field.String("photo_url").Optional().
			Comment("URL to member's photo"),
		field.Bytes("photo_data").Optional().
			Comment("Raw photo data (alternative to URL)"),
		field.String("photo_hash").Optional().
			Comment("Hash of photo for change detection"),

		// Administrative Fields
		field.Bool("is_active").Default(true).
			Comment("Whether membership is active"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("When record was created"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("When record was last updated"),
	}
}

// Edges of the Member.
func (Member) Edges() []ent.Edge {
	return nil
}

// Indexes of the Member.
func (Member) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id_number").Unique(),
		index.Fields("surname"),
		index.Fields("mobile"),
		index.Fields("email").Unique(),
		index.Fields("membership_year"),
	}
}
