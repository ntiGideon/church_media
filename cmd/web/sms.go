package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AfricasTalkingLtd/africastalking-go/sms"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
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

func (app *application) sendMessage(to string, message string) error {
	senderId := "Ascension"
	apiKey := os.Getenv("AR_API_KEY")

	baseURL := "https://sms.arkesel.com/sms/api"
	params := url.Values{}
	params.Add("action", "send-sms")
	params.Add("api_key", apiKey)
	params.Add("to", to)
	params.Add("from", senderId)
	params.Add("sms", message)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(fullURL)
	if err != nil {
		app.logger.Error("Failed to send SMS request: %v", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		app.logger.Error("Failed to read response body: %v", err)
		return err
	}

	app.logger.Info("SMS response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		app.logger.Error("Failed to send SMS, status: %d, response: %s", resp.StatusCode, string(body))
		return fmt.Errorf("SMS API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (app *application) sendTermiiSMS(to string, message string) error {
	apiKey := os.Getenv("TERMII_API_KEY")
	sender := "ABC-Appiadu"

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
		app.logger.Error("failed to parse response: %w", "err", err.Error())
		return errors.New(err.Error())
	}

	if smsResp.Code != 200 {
		app.logger.Error("failed to send SMS: %s", "err", smsResp.Message)
		return errors.New(smsResp.Message)
	}

	app.logger.Info("SMS sent successfully:", "success", smsResp.Message)
	return nil
}
