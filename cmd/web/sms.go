package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AfricasTalkingLtd/africastalking-go/sms"
	"net/http"
	"os"
)

type SMSRequest struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Sms     string `json:"sms"`
	Type    string `json:"type"`    // "plain" or "unicode"
	Channel string `json:"channel"` // "generic" or "dnd"
	APIKey  string `json:"api_key"`
}

type SMSResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (app *application) sendSMS(phone string) error {
	username := os.Getenv("AT_USERNAME")
	apikey := os.Getenv("AT_API_KEY")
	env := "Sandbox"

	smsService := sms.NewService(apikey, username, env)

	response, err := smsService.Send("Ascension Baptist Church - Appiadu", phone, "Happy Birthday to you!")
	if err != nil {
		return err
	}
	fmt.Println(response)
	return nil
}

func (app *application) sendTermiiSMS(to string, message string) error {
	apiKey := os.Getenv("TERMII_API_KEY") // Set this in your .env or OS
	sender := "ABC-Appiadu"               // Sender ID (or "Termii" for test)

	// Build request
	payload := SMSRequest{
		To:      to,
		From:    sender,
		Sms:     message,
		Type:    "plain",
		Channel: "generic",
		APIKey:  apiKey,
	}

	body, _ := json.Marshal(payload)

	termiiUrl := fmt.Sprintf("%s/api/sms/send", os.Getenv("TERMII_URL"))

	req, err := http.NewRequest("POST", termiiUrl, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer res.Body.Close()

	smsResp := SMSResponse{
		Code:    208683,
		Message: message,
	}
	err = json.NewDecoder(res.Body).Decode(&smsResp)
	if err != nil {
		app.logger.Error("failed to parse response: %w", err.Error())
		return errors.New(err.Error())
	}

	if smsResp.Code != 200 {
		app.logger.Error("failed to send SMS: %s", smsResp.Message)
		return errors.New(smsResp.Message)
	}

	app.logger.Info("SMS sent successfully:", smsResp.Message)
	return nil
}
