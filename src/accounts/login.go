package accounts

import (
	"log"
	"strconv"
	"strings"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/methods"
	"git.xenonstack.com/xs-onboarding/accounts/src/token"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
)

//============================================ LOGIN ========================================

//LoginEndPoint is used to login
func LoginEndPoint(email, password string) (map[string]interface{}, int) {
	defer util.Panic()
	db := config.DB
	mapd := make(map[string]interface{})
	email = strings.ToLower(email)
	account := []database.Accounts{}
	db.Where("email=?", email).Find(&account)
	if len(account) == 0 {
		mapd["error"] = true
		mapd["message"] = "Account does not exist"
		return mapd, 404
	}
	if account[0].AccountStatus == "blocked" {
		mapd["error"] = true
		mapd["message"] = "Your account has been deactivated"
		return mapd, 400
	}
	msg, isLocked, count := checkPreviousFailedLogins(account[0])
	if isLocked {
		db.Model(&database.Accounts{}).Where("email=?", email).Update("account_status", "locked")
		mapd["error"] = true
		mapd["message"] = msg
		return mapd, 400
	}
	iscorrect := methods.CheckHashForPassword(account[0].Password, password)
	if iscorrect == false {
		mapd["error"] = true
		mapd["message"] = "Invalid Credentials. You have " + strconv.Itoa(count) + " login attempts left"
		return mapd, 400
	}
	mapd = token.GenerateJwtToken(account[0])
	personal := database.PersonalInfo{}
	err := db.Where("userid=?", account[0].ID).Find(&personal).Error
	if err != nil {
		log.Println(err)
	}
	mapd["error"] = false
	mapd["message"] = "Login Successful"
	mapd["role_id"] = account[0].Role
	if account[0].Level != "" {
		mapd["Level"] = account[0].Level
	}
	mapd["name"] = account[0].Name
	mapd["email"] = account[0].Email
	mapd["link"] = personal.ProfileImgLink
	return mapd, 200
}

func checkPreviousFailedLogins(account database.Accounts) (string, bool, int) {
	defer util.Panic()
	// declaring variables
	var lockFor int64 = 3600
	var failedloginCount int
	var msg string
	var isLocked bool

	// connecting to db
	db := config.DB

	var LastFailedAttempt int64
	// extracting activities on bsis of id
	var activities []database.Activities
	db.Raw("select * from activities where email= '" + account.Email + "' order by timestamp desc limit 5;").Scan(&activities)
	for i := 0; i < len(activities); i++ {
		// if activity name is failed login and checking time interval is less then lockfor
		if activities[i].ActivityName == "failedlogin" && (time.Now().Unix()-activities[i].Timestamp) < lockFor {
			if i == 0 {
				// setting last failed attemp
				LastFailedAttempt = activities[i].Timestamp
			}
			// incrementing failed login count
			failedloginCount++
		}
	}
	log.Println(LastFailedAttempt)
	// is count is more then equal to 5
	if failedloginCount >= 5 {
		msg = "Your account has been locked due to five invalid attempts. Please reset your password by clicking forgot password "
		isLocked = true
		return msg, isLocked, 0
	}

	return "", false, 5 - failedloginCount
}
