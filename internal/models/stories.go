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

// CountPublishedStories returns the total count of published stories
func (m *StoryModel) CountPublishedStories(ctx context.Context) (int, error) {
	return m.Db.Story.Query().
		Where(story.StatusEQ("published")).
		Count(ctx)
}

// GetPaginatedStories returns stories for a specific page
func (m *StoryModel) GetPaginatedStories(ctx context.Context, page, perPage int) ([]*ent.Story, error) {
	offset := (page - 1) * perPage
	return m.Db.Story.Query().
		Where(story.StatusEQ("published")).
		Order(ent.Desc(story.FieldPublishedAt)).
		Offset(offset).
		Limit(perPage).
		WithAuthor(func(query *ent.UserQuery) {
			query.WithContactProfile()
		}).
		All(ctx)
}

func (m *StoryModel) IncrementLike(ctx context.Context, storyID int) (int, error) {
	storyLike, err := m.Db.Story.UpdateOneID(storyID).
		AddLikes(1).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return storyLike.Likes, nil
}

func (m *StoryModel) IncrementDislike(ctx context.Context, storyID int) (int, error) {
	storyDislike, err := m.Db.Story.UpdateOneID(storyID).
		AddDislikes(1).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return storyDislike.Dislikes, nil
}

func (m *StoryModel) GetStoryByID(ctx context.Context, id int) (*ent.Story, error) {
	return m.Db.Story.Query().
		Where(story.IDEQ(id)).
		WithAuthor(func(query *ent.UserQuery) {
			query.WithContactProfile()
		}).
		Only(ctx)
}

func (m *StoryModel) DeleteStory(ctx context.Context, id int) error {
	return m.Db.Story.DeleteOneID(id).Exec(ctx)
}
