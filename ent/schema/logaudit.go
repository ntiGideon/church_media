package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// LogAudit holds the schema definition for the LogAudit entity.
type LogAudit struct {
	ent.Schema
}

// Fields of the LogAudit.
func (LogAudit) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("action").NotEmpty(),
		field.String("entity_type").NotEmpty(),
		field.Int("entity_id"),
		field.JSON("entity_data", map[string]interface{}{}).Optional(),
		field.Int("created_by").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.String("ip_address").MaxLen(45),
		field.String("user_agent").Optional(),
		field.String("request_id").Optional(),
		field.JSON("metadata", map[string]interface{}{}).Optional(),
	}
}

// Edges of the LogAudit.
func (LogAudit) Edges() []ent.Edge {
	return nil
}

// Indexes of the AuditLog.
func (LogAudit) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("action"),
		index.Fields("entity_type"),
		index.Fields("entity_id"),
		index.Fields("created_at"),
		index.Fields("created_by"),
	}
}
