package models

import (
	"context"
	"github.com/ogidi/church-media/ent"
	"github.com/ogidi/church-media/ent/story"
	"time"
)

type StoryModelInterface interface {
}

type StoryModel struct {
	Db *ent.Client
}

func (m *StoryModel) CreateStory(ctx context.Context, storiesDto StoriesDto) error {
	_, err := m.Db.Story.Create().
		SetAuthorID(storiesDto.AuthorId).
		SetBody(storiesDto.Body).
		SetExcerpt(storiesDto.Excerpt).
		SetImage(storiesDto.Image).
		SetStatus(storiesDto.Status).
		SetTitle(storiesDto.Title).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *StoryModel) UpdateStory(ctx context.Context, id int, storiesDto StoriesDto) error {
	_, err := m.Db.Story.UpdateOneID(id).
		SetAuthorID(storiesDto.AuthorId).
		SetBody(storiesDto.Body).
		SetExcerpt(storiesDto.Excerpt).
		SetImage(storiesDto.Image).
		SetTitle(storiesDto.Title).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *StoryModel) PublishStory(ctx context.Context, id int) (*ent.Story, error) {
	return m.Db.Story.UpdateOneID(id).
		SetStatus("published").
		SetPublishedAt(time.Now()).
		Save(ctx)
}

func (m *StoryModel) GetPublishedStories(ctx context.Context, limit int) ([]*ent.Story, error) {
	return m.Db.Story.Query().
		Where(story.StatusEQ("published")).
		Order(ent.Desc(story.FieldPublishedAt)).
		Limit(limit).
		WithAuthor(func(q *ent.UserQuery) {
			q.WithContactProfile()
		}).
		All(ctx)
}

func (m *StoryModel) GetStoryByID(ctx context.Context, id int) (*ent.Story, error) {
	return m.Db.Story.Query().
		Where(story.IDEQ(id)).
		WithAuthor().
		Only(ctx)
}

func (m *StoryModel) DeleteStory(ctx context.Context, id int) error {
	return m.Db.Story.DeleteOneID(id).Exec(ctx)
}
