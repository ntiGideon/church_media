package models

import (
	"context"
	"github.com/ogidi/church-media/ent"
	"github.com/ogidi/church-media/ent/event"
	"time"
)

type EventModelInterface interface {
	CreateEvent(ctx context.Context, dto *CreateEventDto) error
}

type EventModel struct {
	Db *ent.Client
}

func (m *EventModel) CreateEvent(ctx context.Context, dto *CreateEventDto) error {
	_, err := m.Db.Event.Create().
		SetTitle(dto.Title).
		SetDescription(dto.Description).
		SetStartTime(dto.StartTime).
		SetEndTime(dto.EndTime).
		SetLocation(dto.Location).
		SetImageURL(dto.ImageUrl).
		SetFeatured(dto.Featured).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *EventModel) UpcomingEvents(ctx context.Context, limit int) ([]*ent.Event, error) {
	upcomingEvents, err := m.Db.Event.Query().
		Where(event.StartTimeGTE(time.Now())).
		Where(event.StartTimeLTE(time.Now().AddDate(0, 0, 30))).
		Order(ent.Asc(event.FieldStartTime)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return upcomingEvents, nil
}
