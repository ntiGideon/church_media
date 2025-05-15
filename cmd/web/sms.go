package main

import (
	"fmt"
	"github.com/AfricasTalkingLtd/africastalking-go/sms"
	"os"
)

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
