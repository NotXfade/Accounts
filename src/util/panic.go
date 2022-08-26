//+build !test

package util

import (
	"encoding/json"
	"log"
	"runtime/debug"

	"git.xenonstack.com/xs-onboarding/accounts/config"
)

//Panic function used when the panic is iccur then it send the alert
func Panic() {
	if recover() != nil {
		panic := string(debug.Stack())
		log.Println(panic)
		mails := []string{"neha.sharma@xenonstack.com"}

		for i := 0; i < len(mails); i++ {
			go SendMail(mails[i], panic)
		}

	}
}

//MailData is used to send
type MailData struct {
	Email string `json:"email"`
	Body  string `json:"body"`
}

//SendMail function is used to send normal mail to user
func SendMail(email, body string) {
	nc := config.NC
	sendData := MailData{
		Email: email,
		Body:  body,
	}
	data, _ := json.Marshal(sendData)
	nc.Publish(config.Conf.NatsServer.Subject+".notifications.panics", data)
}
