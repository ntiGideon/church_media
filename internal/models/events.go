package models

import (
	"context"
	"fmt"
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

func (m *EventModel) EditEvent(ctx context.Context, dto *CreateEventDto, id int) error {
	startTime, _ := time.Parse("2006-01-02T15:04", dto.StartTime)
	endTime, _ := time.Parse("2006-01-02T15:04", dto.EndTime)
	if dto.ImageUrl != "" {
		_, err := m.Db.Event.UpdateOneID(id).
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
	} else {
		_, err := m.Db.Event.UpdateOneID(id).
			SetTitle(dto.Title).
			SetDescription(dto.Description).
			SetStartTime(startTime).
			SetEndTime(endTime).
			SetLocation(dto.Location).
			SetFeatured(dto.Featured).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	return nil
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

func (m *EventModel) GetEvents(ctx context.Context, filter EventFilter) ([]*ent.Event, int, error) {
	query := m.Db.Event.Query()

	// Apply filters
	if filter.Search != "" {
		query = query.Where(
			event.Or(
				event.TitleContains(filter.Search),
				event.DescriptionContains(filter.Search),
				event.LocationContains(filter.Search),
			),
		)
	}

	if filter.Featured != nil {
		query = query.Where(event.FeaturedEQ(*filter.Featured))
	}

	if !filter.FromDate.IsZero() {
		query = query.Where(event.StartTimeGTE(filter.FromDate))
	}

	if !filter.ToDate.IsZero() {
		query = query.Where(event.StartTimeLTE(filter.ToDate))
	}

	// Apply sorting
	switch filter.SortBy {
	case "title":
		if filter.Order == "asc" {
			query = query.Order(ent.Asc(event.FieldTitle))
		} else {
			query = query.Order(ent.Desc(event.FieldTitle))
		}
	case "date":
		if filter.Order == "asc" {
			query = query.Order(ent.Asc(event.FieldStartTime))
		} else {
			query = query.Order(ent.Desc(event.FieldStartTime))
		}
	case "location":
		if filter.Order == "asc" {
			query = query.Order(ent.Asc(event.FieldLocation))
		} else {
			query = query.Order(ent.Desc(event.FieldLocation))
		}
	default:
		query = query.Order(ent.Desc(event.FieldStartTime))
	}

	// Count total before pagination
	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}

	// Apply pagination
	if filter.PageSize > 0 {
		query = query.Offset((filter.Page - 1) * filter.PageSize).Limit(filter.PageSize)
	}

	events, err := query.All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query events: %w", err)
	}

	return events, total, nil
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

func (m *EventModel) GetEventByID(ctx context.Context, id int) (*ent.Event, error) {
	eventData, err := m.Db.Event.Query().Where(event.ID(id)).First(ctx)
	if err != nil {
		return nil, err
	}
	return eventData, nil
}
