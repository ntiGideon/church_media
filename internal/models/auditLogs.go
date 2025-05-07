package models

import (
	"context"
	"github.com/ogidi/church-media/ent"
)

type AuditLogsModel struct {
	Db *ent.Client
}

func (m *AuditLogsModel) CreateLogs(ctx context.Context, entry AuditLogEntry) (*ent.LogAudit, error) {
	return m.Db.LogAudit.Create().
		SetAction(entry.Action).
		SetEntityType(entry.EntityType).
		SetEntityID(entry.EntityID).
		SetNillableCreatedBy(entry.CreatedBy).
		SetIPAddress(entry.IPAddress).
		SetUserAgent(entry.UserAgent).
		SetRequestID(entry.RequestID).
		SetEntityData(entry.EntityData).
		SetMetadata(entry.Metadata).
		Save(ctx)
}
