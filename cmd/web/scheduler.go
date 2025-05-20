package main

import (
	"context"
	"github.com/ogidi/church-media/ent"
	"github.com/robfig/cron/v3"
	"time"
)

func (app *application) Start() {
	loc, err := time.LoadLocation("Africa/Accra")
	if err != nil {
		app.logger.Error("failed to load location", "error", err.Error())
		return
	}

	c := cron.New(
		cron.WithLocation(loc),
	)

	id, err := c.AddFunc("40 11 * * *", app.sendBirthdayMessages)
	if err != nil {
		app.logger.Error("failed to add birthday messages", "error", err.Error())
		return
	}

	app.logger.Info("Added birthday messages job", "entryID", id, "schedule", "0 11 * * *", "timezone", loc.String())

	go func() {
		c.Start()
		select {} // Keep the goroutine alive
	}()

	app.logger.Info("Birthday scheduler started successfully")
}

func (app *application) sendBirthdayMessages() {
	app.logger.Info("Starting birthday message processing")
	defer func() {
		app.logger.Info("Completed birthday message processing")
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	membersWithBirthday, err := app.memberClient.GetBirthdayMembers(ctx)
	if err != nil {
		app.logger.Error("failed to get birthday members", "error", err.Error())
		return
	}

	app.logger.Info("Found birthday members", "count", len(membersWithBirthday))

	for i, member := range membersWithBirthday {
		app.logger.Debug("Processing member",
			"index", i,
			"memberID", member.ID,
			"mobile", member.Mobile,
		)

		message := generateBirthdayMessage(member)
		if member.Mobile == "" {
			app.logger.Warn("Skipping member with no mobile number", "memberID", member.ID)
			continue
		}

		if err := app.sendMessage(member.Mobile, message); err != nil {
			app.logger.Error("failed to send message",
				"error", err.Error(),
				"memberID", member.ID,
				"mobile", member.Mobile,
			)
			continue
		}

		app.logger.Info("Successfully sent birthday message",
			"memberID", member.ID,
			"mobile", member.Mobile,
		)
	}
}

func generateBirthdayMessage(m *ent.Member) string {
	return "Dear " + m.OtherNames + ", Ascension Baptist Church - Appiadu wishes you a blessed birthday! May God's grace abound in your life today and always."
}
