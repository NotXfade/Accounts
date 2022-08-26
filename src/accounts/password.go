package accounts

import (
	"encoding/json"
	"log"

	"git.xenonstack.com/xs-onboarding/accounts/src/methods"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/token"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
)

//======================== SEND FORGET PASSWORD LINK =========================

//ForgotPassword is used to send reset password link to the mail of user
func ForgotPassword(email string) (map[string]interface{}, int) {
	//handler panic and Alerts
	defer util.Panic()
	mapd := make(map[string]interface{})
	//checking whether email exist or not
	acc, err := GetAccountForEmail(email)
	if err != nil {
		//log.Println("when no account is there")
		mapd["error"] = true
		mapd["message"] = "Email doesn't exists."
		return mapd, 400
	}
	db := config.DB
	if acc.AccountStatus == "locked" {
		db.Model(&database.Accounts{}).Where("email=?", email).Update("account_status", "active")
		acc.AccountStatus = "active"
	}
	//Check If Account is Active
	if acc.AccountStatus != "active" {
		//log.Println("Account is not active")
		mapd["error"] = true
		mapd["message"] = "Account doesn't exist"
		return mapd, 400
	}
	//generate and store token in database
	tok := token.GenerateToken(email, "reset-password", "")
	//----------------------------- request to notification service through NATS----------------
	//payload for sending invite through email
	payload := database.MailData{
		Email: email,
		Link:  "https://" + config.Conf.Service.HostAddr + "/reset-password?token=" + tok,
		Task:  "resetPassword",
	}
	//Marshal data to bytes
	sendData, err := json.Marshal(payload)
	if err != nil {
		log.Print("error while marshalling data")
	}
	nc := config.NC
	//send request to notification service for sending invite email
	err = nc.Publish(config.Conf.NatsServer.Subject+".notifications.accounts", sendData)
	if err != nil {
		//if there is error sending invite return the error and delete token as token
		//is for one time use only
		mapd["error"] = true
		mapd["message"] = "Couldn't send message"
		go token.DeleteToken(tok)
		log.Print(err)
		return mapd, 500
	}
	mapd["error"] = false
	mapd["message"] = "Email sent successfully"
	return mapd, 200
}

//=============================== SET NEW PASSWORD ==================================

//SetNewPassword is used to update old password with new password
func SetNewPassword(email, password string) (map[string]interface{}, int) {
	//handler panic and Alerts
	defer util.Panic()
	mapd := make(map[string]interface{})
	//checking whether email exist or not
	acc, err := GetAccountForEmail(email)
	if err != nil {
		//log.Println("when no account is there")
		mapd["error"] = true
		mapd["message"] = "Email doesn't exists."
		return mapd, 400
	}
	//Check If Account is Active
	if acc.AccountStatus != "active" {
		//log.Println("Account is not active")
		mapd["error"] = true
		mapd["message"] = "Account doesn't exist"
		return mapd, 400
	}
	passHash := methods.HashForNewPassword(password)
	db := config.DB
	err = db.Model(&database.Accounts{}).Where("id=?", acc.ID).Update("password", passHash).Error
	if err != nil {
		log.Println("Database Error : ", err)
		mapd["error"] = true
		mapd["message"] = "Internal Server Error : Couldn't update value"
		return mapd, 500
	}
	mapd["error"] = false
	mapd["message"] = "New password set successfully"
	return mapd, 200
}

//================================ CHANGE PASSWORD ==========================================

//ChangePassword is used to change old password
func ChangePassword(oldPassword, newPassword, email string) (map[string]interface{}, int) {
	//handler panic and Alerts
	defer util.Panic()
	mapd := make(map[string]interface{})
	account, err := GetAccountForEmail(email)
	if err != nil {
		mapd["error"] = true
		mapd["message"] = "Account doesn't exist"
		return mapd, 400
	}
	isMatch := methods.CheckHashForPassword(account.Password, oldPassword)
	if isMatch != true {
		mapd["error"] = true
		mapd["message"] = "Please enter correct current password"
		return mapd, 400
	}
	passHash := methods.HashForNewPassword(newPassword)
	db := config.DB
	err = db.Model(&database.Accounts{}).Where("id=?", account.ID).Update("password", passHash).Error
	if err != nil {
		log.Println("Database Error : ", err)
		mapd["error"] = true
		mapd["message"] = "Internal Server Error : Couldn't update value"
		return mapd, 500
	}
	mapd["error"] = false
	mapd["message"] = "New password set successfully"
	return mapd, 200
}
