package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"time"
)

type EmailDto struct {
	Firstname       string
	To              string
	Subject         string
	Token           string
	Expiration      int // in hours
	RegistrationURL string
	TemplatePath    string
}

type TemplateData struct {
	Firstname       string
	Token           string
	Expiration      int
	RegistrationURL string
	CurrentYear     int
}

var mimeHeaders = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

func EmailLogics(subject, templatePath string, emailDto *EmailDto, templateData interface{}) error {

	auth := smtp.PlainAuth("", os.Getenv("MAIL_FROM"), os.Getenv("MAIL_PASSWORD"), os.Getenv("MAIL_HOST"))

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	to := []string{emailDto.To}

	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))

	err = t.Execute(&body, templateData)
	if err != nil {
		return err
	}

	err = smtp.SendMail(os.Getenv("MAIL_HOST")+":"+os.Getenv("MAIL_PORT"), auth, os.Getenv("MAIL_FROM"), to, body.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (app *application) sendInvitationEmail(emailDto *EmailDto, subject string) error {
	data := TemplateData{
		Token:           emailDto.Token,
		Expiration:      emailDto.Expiration,
		RegistrationURL: emailDto.RegistrationURL,
		Firstname:       emailDto.Firstname,
		CurrentYear:     time.Now().Year(),
	}

	err := EmailLogics(subject, emailDto.TemplatePath, emailDto, data)
	if err != nil {
		app.logger.Error("couldn't send email to : ", emailDto.To, " an error at ::::> "+err.Error())
		return err
	}
	app.logger.Info("successfully sent invitation email to : ", "error", emailDto.To)
	return nil
}
