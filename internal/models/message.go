package models

import (
	"context"
	"github.com/ogidi/church-media/ent"
	"github.com/ogidi/church-media/ent/message"
	"github.com/ogidi/church-media/ent/response"
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
	_, err := m.Db.Message.Create().SetName(dto.Name).SetEmail(dto.Email).SetPhone(dto.Phone).SetDescription(dto.Description).SetSubject(message.Subject(dto.Subject)).SetState(message.StateUNREAD).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *MessageModel) GetMessages(ctx context.Context, filter string) ([]*ent.Message, error) {
	query := m.Db.Message.Query().WithResponses().Order(ent.Desc(message.FieldCreatedAt))

	switch filter {
	case "unread":
		query.Where(message.StateEQ(message.StateUNREAD))
	case "responded":
		query.Where(message.StateEQ(message.StateRESPONDED))
	case "prayer":
		query.Where(message.SubjectEQ(message.SubjectPRAYER_REQUEST))
	}

	return query.All(ctx)
}

func (m *MessageModel) GetMessageWithResponses(ctx context.Context, id int) (*ent.Message, error) {
	return m.Db.Message.Query().
		Where(message.IDEQ(id)).
		WithResponses(func(q *ent.ResponseQuery) {
			q.Order(ent.Desc(response.FieldCreatedAt))
		}).
		Only(ctx)
}

func (m *MessageModel) GetUnreadCount(ctx context.Context) (int, error) {
	return m.Db.Message.Query().
		Where(message.StateEQ(message.StateUNREAD)).
		Count(ctx)
}

func (m *MessageModel) MarkAsRead(ctx context.Context, id int) error {
	_, err := m.Db.Message.UpdateOneID(id).
		SetState(message.StateREAD).
		Save(ctx)
	return err
}

func (m *MessageModel) MarkAsResponded(ctx context.Context, id int) error {
	_, err := m.Db.Message.UpdateOneID(id).
		SetState(message.StateRESPONDED).
		Save(ctx)
	return err
}

func (m *MessageModel) SearchMessages(ctx context.Context, filter, query string) ([]*ent.Message, error) {
	q := m.Db.Message.Query().
		Where(
			message.Or(
				message.NameContains(query),
				message.EmailContains(query),
				message.DescriptionContains(query),
				message.SubjectEQ(message.Subject(query)),
			),
		)

	switch filter {
	case "unread":
		q.Where(message.StateEQ(message.StateUNREAD))
	case "responded":
		q.Where(message.StateEQ(message.StateRESPONDED))
	case "prayer":
		q.Where(message.SubjectEQ(message.SubjectPRAYER_REQUEST))
	}

	return q.Order(ent.Desc(message.FieldCreatedAt)).All(ctx)
}

func (m *MessageModel) DeleteMessage(ctx context.Context, id int) error {
	return m.Db.Message.DeleteOneID(id).Exec(ctx)
}

func (m *MessageModel) UpdateState(ctx context.Context, id int, state string) error {
	_, err := m.Db.Message.UpdateOneID(id).
		SetState(message.State(state)).
		Save(ctx)
	return err
}

func (m *MessageModel) GetResponses(ctx context.Context, messageID int) ([]*ent.Response, error) {
	return m.Db.Response.Query().
		Where(response.HasMessageWith(message.IDEQ(messageID))).
		Order(ent.Desc(response.FieldCreatedAt)).
		All(ctx)
}

func (m *MessageModel) CreateResponse(ctx context.Context, messageID, userID int, email, subject, content string) (*ent.Response, error) {
	return m.Db.Response.Create().
		SetMessageID(messageID).
		SetUserID(userID).
		SetEmail(email).
		SetSubject(response.Subject(subject)).
		SetDescription(content).
		Save(ctx)
}

// GetUnreadMessagesCount counts all messages that are not unread in the system
func (m *MessageModel) GetUnreadMessagesCount(ctx context.Context) int {
	messages, _ := m.Db.Message.Query().Where(message.StateEQ(message.StateUNREAD)).Count(ctx)
	return messages
}
