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
	startTime, _ := time.Parse("2006-01-02T15:04", dto.StartTime)
	endTime, _ := time.Parse("2006-01-02T15:04", dto.EndTime)
	_, err := m.Db.Event.Create().
		SetTitle(dto.Title).
		SetDescription(dto.Description).
		SetStartTime(startTime).
		SetEndTime(endTime).
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

func (m *EventModel) FeaturedEvents(ctx context.Context) ([]*ent.Event, error) {
	eventsData, err := m.Db.Event.Query().Where(event.Featured(true)).All(ctx)
	if err != nil {
		return nil, err
	}
	return eventsData, nil
}

func (m *EventModel) AllEvents(ctx context.Context) ([]*ent.Event, error) {
	eventsData, err := m.Db.Event.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return eventsData, nil
}

func (m *EventModel) EventExist(ctx context.Context, id int) (*ent.Event, error) {
	existingEvent, err := m.Db.Event.Query().Where(event.IDEQ(id)).First(ctx)
	if err != nil {
		return nil, err
	}

	return existingEvent, nil
}

func (m *EventModel) DeleteEvent(ctx context.Context, id int) error {
	err := m.Db.Event.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *EventModel) ViewEvent(ctx context.Context, id int) (*ent.Event, error) {
	eventData, err := m.Db.Event.Query().Where(event.ID(id)).First(ctx)
	if err != nil {
		return nil, err
	}
	return eventData, nil
}

func (m *EventModel) ListEvents(ctx context.Context, page int, pageSize int) ([]*ent.Event, int, error) {
	eventList, err := m.Db.Event.Query().Order(ent.Asc(event.FieldStartTime)).Offset((page - 1) * pageSize).Limit(pageSize).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	total, err := m.Db.Event.Query().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return eventList, total, nil
}
