package models

import (
	"context"
	"github.com/ogidi/church-media/ent"
	"github.com/ogidi/church-media/ent/message"
	"time"
)

type Message struct {
	ID          int
	Name        string
	Email       string
	Phone       string
	Subject     string
	Description string
	State       message.State
	CreatedAt   time.Time
}

type MessageModelInterface interface {
	CreateMessage(ctx context.Context, dto *CreateMessageDto) error
}

type MessageModel struct {
	Db *ent.Client
}

// CreateMessage helps create a new contact message
func (m *MessageModel) CreateMessage(ctx context.Context, dto CreateMessageDto) error {
	existing, err := m.Db.Message.Query().Where(message.Phone(dto.Phone)).Exist(ctx)
	if err != nil {
		return err
	}
	if existing {
		return PhoneInUse
	}
	_, err = m.Db.Message.Create().SetName(dto.Name).SetEmail(dto.Email).SetPhone(dto.Phone).SetDescription(dto.Description).SetSubject(message.Subject(dto.Subject)).SetState(message.StateUNREAD).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *MessageModel) GetMessages(ctx context.Context, keyword string) ([]*ent.Message, error) {
	query := m.Db.Message.Query()

	if keyword == "unread" {
		query.Where(message.StateEQ(message.StateUNREAD))
	}

	if keyword == "responded" {
		query.Where(message.StateEQ(message.StateRESPONDED))
	}

	if keyword == "prayerRequest" {
		query.Where(message.SubjectEQ(message.SubjectPRAYER_REQUEST))
	}

	messages, err := query.Order(ent.Desc(message.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// GetUnreadMessagesCount counts all messages that are not unread in the system
func (m *MessageModel) GetUnreadMessagesCount(ctx context.Context) int {
	messages, _ := m.Db.Message.Query().Where(message.StateEQ(message.StateUNREAD)).Count(ctx)
	return messages
}
