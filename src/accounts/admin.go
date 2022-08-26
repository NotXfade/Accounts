package accounts

import (
	"encoding/json"
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

//InviteResponse is used for response
type InviteResponse struct {
	Email   string `json:"email"`
	Role    string `json:"role"`
	Message string `json:"message"`
}

//InviteAdmin is used to invite admin
func InviteAdmin(data database.Inviteadmin) (map[string]interface{}, int) {
	//handler panic and Alerts
	defer util.Panic()
	mapd := make(map[string]interface{})
	nc := config.NC
	listnotsent := []InviteResponse{}
	listsent := []InviteResponse{}
	for i := 0; i < len(data.Invite); i++ {
		notsent := InviteResponse{}
		role := data.Invite[i].Role
		email := data.Invite[i].Email
		if role != "P&C" && role != "Business Analyst" && role != "Reviewer" && role != "Admin" {
			notsent.Email = email
			notsent.Role = role
			notsent.Message = "Invalid Role"
			listnotsent = append(listnotsent, notsent)
			continue
		}
		email = strings.ToLower(email)
		//check if email is valid
		isvalid := methods.ValidateEmail(email)
		if isvalid == false {
			log.Println(data.Invite[i])
			notsent.Email = data.Invite[i].Email
			notsent.Role = role
			notsent.Message = "Invalid Email"
			listnotsent = append(listnotsent, notsent)
			continue
		}
		//check if account already exists
		acc, _ := GetAccountForEmail(email)
		if acc.AccountStatus == "active" || acc.AccountStatus == "blocked" || acc.Role == "user" {
			log.Println(data.Invite[i])
			notsent.Email = data.Invite[i].Email
			notsent.Role = data.Invite[i].Role
			notsent.Message = "Account already exists"
			listnotsent = append(listnotsent, notsent)
			continue
		}
		if acc.Role != data.Invite[i].Role && acc.Role != "" {
			notsent.Email = data.Invite[i].Email
			notsent.Role = data.Invite[i].Role
			notsent.Message = "You had already invited user with role " + acc.Role
			listnotsent = append(listnotsent, notsent)
			continue
		}
		type AdminData struct {
			Email string `json:"email"`
			Link  string `json:"link"`
			Task  string `json:"task"`
			Role  string `json:"role"`
		}
		if isvalid == true {
			if acc.AccountStatus != "invited" {
				account := database.Accounts{
					Email:          email,
					Role:           strings.ToLower(role),
					Level:          "",
					AccountStatus:  "invited",
					AcceptedPolicy: false,
					Timestamp:      time.Now().Unix(),
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}
				db := config.DB
				db.Create(&account)
			}
			//generate and store token in database
			tok := token.GenerateToken(email, "adminInvite", role)
			payload := AdminData{
				Email: data.Invite[i].Email,
				Link:  "https://" + config.Conf.Service.HostAddr + "/admin-registration?token=" + tok,
				Task:  "inviteAdmin",
				Role:  role,
			}
			sendData, err := json.Marshal(payload)
			if err != nil {
				notsent.Email = data.Invite[i].Email
				notsent.Role = role
				notsent.Message = "Error while sending notifications"
				log.Print("error while marshalling data")
				go token.DeleteToken(tok)
				listnotsent = append(listnotsent, notsent)
				continue
			}
			//send request to notification service for sending invite email
			err = nc.Publish(config.Conf.NatsServer.Subject+".notifications.inviteadmin", sendData)
			if err != nil {
				//if there is error sending invite return the error and delete token as token
				//is for one time use only
				notsent.Email = data.Invite[i].Email
				notsent.Role = data.Invite[i].Role
				notsent.Message = "Error while sending notifications"
				go token.DeleteToken(tok)
				log.Print(err)
				listnotsent = append(listnotsent, notsent)
				continue
			}
			listsent = append(listsent, InviteResponse{
				Email:   data.Invite[i].Email,
				Role:    data.Invite[i].Role,
				Message: "Sent successfully",
			})
		}
	}
	mapd["error"] = false
	mapd["message"] = "Invite link's sent"
	mapd["unsuccessful_sents"] = listnotsent
	mapd["successful_sents"] = listsent
	return mapd, 200
}

//RegisterAdmin is used to register account of admin
func RegisterAdmin(data database.Registeradmin) (map[string]interface{}, int) {
	defer util.Panic()
	db := config.DB
	mapd := make(map[string]interface{})
	tokens := database.Tokens{}
	var err error
	tokens, err = token.VerifyToken(data.Token, "adminInvite")
	if err != nil {
		mapd["error"] = true
		mapd["message"] = "Your invite is expired"
		go token.DeleteToken(tokens.Token)
		return mapd, 400
	}
	account, err := GetAccountForEmail(tokens.Email)
	if err != nil {
		mapd["error"] = true
		mapd["message"] = "Invalid request"
		go token.DeleteToken(tokens.Token)
		return mapd, 400
	}
	if account.AccountStatus == "active" || account.AccountStatus == "blocked" {
		mapd["error"] = true
		mapd["message"] = "Account already exists"
		go token.DeleteToken(tokens.Token)
		return mapd, 400
	}
	account.Name = data.Name
	account.ContactNo = data.Contact
	account.Password = methods.HashForNewPassword(data.Password)
	account.AccountStatus = "active"
	account.UpdatedAt = time.Now()
	err = db.Save(account).Error
	if err != nil {
		log.Println(err)
		mapd["error"] = "true"
		mapd["message"] = "Internal server error"
		return mapd, 500
	}
	go token.DeleteToken(tokens.Token)
	mapd = token.GenerateJwtToken(account)
	mapd["email"] = account.Email
	mapd["role_id"] = account.Role
	if account.Level != "" {
		mapd["Level"] = account.Level
	}
	mapd["name"] = account.Name
	mapd["error"] = "false"
	mapd["message"] = "Registration Successful"
	//mapd["link"] = link
	return mapd, 200
}

//ListAdmin is used to list admins
func ListAdmin(Role, limit, page string) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	db := config.DB
	accounts := []database.Accounts{}
	limits, _ := strconv.Atoi(limit)
	pageno, _ := strconv.Atoi(page)
	offset := (limits * (pageno - 1))
	var count int
	if Role != "" {
		err := db.Where("role=?", Role).Offset(offset).Limit(limits).Order("lower(name)").Find(&accounts).Count(&count).Error
		if err != nil {
			log.Println(err)
		}
	} else {
		err := db.Not("role", "user").Offset(offset).Limit(limits).Order("lower(role)").Find(&accounts).Count(&count).Error
		if err != nil {
			log.Println(err)
		}
	}
	if count == 0 {
		mapd["error"] = true
		mapd["message"] = "Users not found"
		return mapd, 404
	}
	mapd["error"] = false
	mapd["list"] = accounts
	mapd["count"] = count
	mapd["message"] = "Operation Successful"
	return mapd, 200
}

//DeleteAdmin is used to delete admin account
func DeleteAdmin(email string) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	db := config.DB
	accounts := []database.Accounts{}
	db.Where("email=?", email).Find(&accounts)
	if len(accounts) == 0 {
		mapd["error"] = true
		mapd["message"] = "Account does not exists"
		return mapd, 404
	}
	if accounts[0].ID == 1 {
		mapd["error"] = true
		mapd["message"] = "You are not authorised to delete this account"
		return mapd, 404
	}

	personalinfo := []database.PersonalInfo{}
	//========== get profile image link===========
	db.Where("userid=?", accounts[0].ID).Find(&personalinfo)

	//========= delete user assigned reviewers ==========
	if accounts[0].Role == "reviewer" {
		DeleteUserAssignedReviewer(accounts[0].ID)
	}
	//=========== delete account ==============
	account := database.Accounts{}
	err := db.Where("id=?", accounts[0].ID).Delete(&account).Error
	if err != nil {
		log.Println(err)
		mapd["error"] = true
		mapd["message"] = "could not delete" + err.Error()
		return mapd, 400
	}

	nc := config.NC
	if len(personalinfo) != 0 {
		//========== delete profile image ==============
		data := []byte(personalinfo[0].ProfileImgLink)
		_, err = nc.Request(config.Conf.NatsServer.Subject+".documents.delete", data, 2*time.Minute)
		if err != nil {
			log.Println(err)
		}
		personal := database.PersonalInfo{}
		db.Where("id=?", personalinfo[0].ID).Delete(&personal)
	}
	mapd["error"] = false
	mapd["message"] = "Account deleted successfully"
	return mapd, 200
}
