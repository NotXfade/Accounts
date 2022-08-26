//+build !test

package nats

import (
	"encoding/json"
	"log"

	"git.xenonstack.com/xs-onboarding/accounts/src/accounts"
)

type Checkemails struct {
	Emails []string `json:"emails"`
	Level  string   `json:"level"`
}

//CheckEmails is used to check and return whether email exists or not
func CheckEmails(data []byte) []byte {
	checkEmails := Checkemails{}
	err := json.Unmarshal(data, &checkEmails)
	if err != nil {
		log.Println(err)
	}
	okmails := make([]string, 0)
	notokmails := make([]string, 0)
	for i := 0; i < len(checkEmails.Emails); i++ {
		acc, err := accounts.GetAccountForEmail(checkEmails.Emails[i])
		if err != nil {
			log.Println(err)
			notokmails = append(notokmails, checkEmails.Emails[i])
		}
		if acc.AccountStatus != "active" {
			notokmails = append(notokmails, checkEmails.Emails[i])
		}
		if acc.Level != checkEmails.Level {
			notokmails = append(notokmails, checkEmails.Emails[i])
		}
		okmails = append(okmails, checkEmails.Emails[i])
	}
	mapd := make(map[string]interface{})
	mapd["okmails"] = okmails
	mapd["notokmails"] = notokmails
	resData, err := json.Marshal(mapd)
	if err != nil {
		log.Println(err)
	}
	return resData
}

//GetDepartmentName is used to get department details
func GetDepartmentName(data []byte) []byte {
	type deptName struct {
		ID int `json:"id"`
	}
	dept := deptName{}
	err := json.Unmarshal(data, &dept)
	if err != nil {
		log.Println(err)
	}
	department := accounts.GetDepartment(dept.ID)
	payload := []byte(department.Name)
	return payload
}
