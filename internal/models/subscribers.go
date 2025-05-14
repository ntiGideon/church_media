package models

import (
	"context"
	"github.com/ogidi/church-media/ent"
	"github.com/ogidi/church-media/ent/subscribe"
)

type SubscriberModel struct {
	Db *ent.Client
}

func (m *SubscriberModel) Subscribe(ctx context.Context, email string) error {
	emailExist, err := m.Db.Subscribe.Query().Where(subscribe.EmailEQ(email)).Exist(ctx)
	if err != nil {
		return err
	}
	if emailExist {
		return EmailExistsError
	}
	_, err = m.Db.Subscribe.Create().SetEmail(email).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *SubscriberModel) GetSubscribers(ctx context.Context) ([]*ent.Subscribe, error) {
	subscribers, err := m.Db.Subscribe.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return subscribers, nil
}
