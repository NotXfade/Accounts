package accounts

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/src/methods"

	"git.xenonstack.com/xs-onboarding/accounts/config"

	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/token"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
)

//================================== Register Account from invite Token ===============================

//RegisterAccount is used to regsiter Member's Account and Activate it
func RegisterAccount(tok, password, name, contact string) (map[string]interface{}, int) {
	//handler panic and Alerts
	defer util.Panic()
	db := config.DB
	testToken := "XSTestV1"
	toks := strings.Split(tok, ".")
	mapd := make(map[string]interface{})
	tokens := database.Tokens{}
	var err error
	var account database.Accounts
	if len(toks) == 1 {
		//Call VerifyToken function to check if token is valid and accessing right info
		tokens, err = token.VerifyToken(toks[0], "onboarding_invite")
		if err != nil {
			mapd["error"] = true
			mapd["message"] = "Your invite is expired"
			return mapd, 400
		}
		//log.Println(tokens)
		account, err = GetAccountForEmail(tokens.Email)
		if err != nil {
			mapd["error"] = true
			mapd["message"] = "Invalid Request"
			return mapd, 400
		}
		if account.AccountStatus == "active" || account.AccountStatus == "blocked" {
			mapd["error"] = true
			mapd["message"] = "Account already exists"
			return mapd, 400
		}
		account.Name = name
		account.ContactNo = contact
		account.Password = password
		account.AccountStatus = "active"
		account.UpdatedAt = time.Now()
		err = db.Save(account).Error
		if err != nil {
			mapd["error"] = true
			mapd["message"] = "Internal server error"
			return mapd, 500
		}
	} else {
		if config.Conf.Service.Environment != "production" {
			if toks[0] == testToken && len(toks) == 3 {
				email := toks[1] + "." + toks[2]
				isvalid := methods.ValidateEmail(email)
				if isvalid == false {
					mapd["error"] = true
					mapd["message"] = "Invalid Email"
					return mapd, 400
				}
				tokens.Email = email
				tokens.Level = "L1"
			} else {
				mapd["error"] = true
				mapd["message"] = "Invalid Token"
				return mapd, 400
			}
		} else {
			mapd["error"] = true
			mapd["message"] = "Invalid Token"
			return mapd, 400
		}
		if tokens.Level == "P&C" || tokens.Level == "Business Analyst" || tokens.Level == "Reviewer" {
			mapd["error"] = true
			mapd["message"] = "Admin accounts are not allowed to register"
			return mapd, 400
		}
		account = database.Accounts{
			Email:          tokens.Email,
			Password:       password,
			Name:           name,
			ContactNo:      contact,
			Role:           "user",
			Level:          "L1",
			AccountStatus:  "active",
			AcceptedPolicy: false,
			Timestamp:      time.Now().Unix(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		//If token is true create account of user
		err = db.Create(&account).Error
		if err != nil {
			//if there is error creating account return the error and delete token as token
			//is for one time use only
			if toks[0] != testToken {
				go token.DeleteToken(tok)
			}
			mapd["error"] = "true"
			mapd["message"] = "either the account on this email already exist or you entered wrong data"
			return mapd, 400
		}
		if toks[0] != testToken {
			go token.DeleteToken(tok)
		}
	}
	go token.DeleteToken(toks[0])
	mapd = token.GenerateJwtToken(account)
	mapd["email"] = account.Email
	mapd["role_id"] = account.Role
	if account.Level != "" {
		mapd["Level"] = account.Level
	}
	mapd["name"] = account.Name
	mapd["error"] = "false"
	mapd["message"] = "Registration Successfull"
	return mapd, 200
}

//=================================== INVITE FOR CREATION OF ACCOUNT =============================

//Sents is used to send back the emails where mail could not be sent and are sent
type Sents struct {
	Email   string `json:"email"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

//Invite : Invite function is used to send invite to user
func Invite(data database.DataInvite) (map[string]interface{}, int) {
	//handler panic and Alerts
	defer util.Panic()
	mapd := make(map[string]interface{})
	nc := config.NC
	listnotsent := []Sents{}
	listsent := []Sents{}
	// jwtToken, errr := LoginCarrerPortal()
	// if errr != "" {
	// 	log.Println(errr)
	// 	mapd["error"] = true
	// 	mapd["message"] = "Internal server error"
	// 	return mapd, 500
	// }
	//log.Println(len(data.Email))
	for i := 0; i < len(data.Invite); i++ {
		notsent := Sents{}
		Level := data.Invite[i].Level
		//log.Println(Level)
		if Level != "L1" && Level != "L2" && Level != "L3" {
			notsent.Email = data.Invite[i].Email
			notsent.Level = data.Invite[i].Level
			notsent.Message = "Invalid level value"
			listnotsent = append(listnotsent, notsent)
			continue
		}
		data.Invite[i].Email = strings.ToLower(data.Invite[i].Email)
		//check if email is valid
		isvalid := methods.ValidateEmail(data.Invite[i].Email)
		if isvalid == false {
			log.Println(data.Invite[i])
			notsent.Email = data.Invite[i].Email
			notsent.Level = data.Invite[i].Level
			notsent.Message = "Invalid Email"
			listnotsent = append(listnotsent, notsent)
			continue
		}
		//check if account already exists
		acc, _ := GetAccountForEmail(data.Invite[i].Email)
		if acc.AccountStatus == "active" || acc.AccountStatus == "blocked" {
			log.Println(data.Invite[i])
			notsent.Email = data.Invite[i].Email
			notsent.Level = data.Invite[i].Level
			notsent.Message = "Account already exists"
			listnotsent = append(listnotsent, notsent)
			continue
		}
		if acc.Role != "user" && acc.Role != "" {
			notsent.Email = data.Invite[i].Email
			notsent.Level = data.Invite[i].Level
			notsent.Message = "Account already exists"
			listnotsent = append(listnotsent, notsent)
			continue
		}
		if acc.Level != data.Invite[i].Level && acc.Level != "" {
			notsent.Email = data.Invite[i].Email
			notsent.Level = data.Invite[i].Level
			notsent.Message = "You had already invited user with level " + acc.Level
			listnotsent = append(listnotsent, notsent)
			continue
		}
		// isExist := CheckIfAccountExistsInCareerPortal(data.Invite[i].Email, data.Invite[i].Level, jwtToken)
		isExist := "false"
		if isExist == "" {
			log.Println(data.Invite[i])
			continue
		} else if isExist == "true" {
			listsent = append(listsent, Sents{
				Email:   data.Invite[i].Email,
				Level:   data.Invite[i].Level,
				Message: "Sent Successfully",
			})
			continue
		}
		if isExist == "false" {
			//log.Println("reached")
			if acc.AccountStatus != "invited" {
				//log.Println("reached")
				account := database.Accounts{
					Email:          data.Invite[i].Email,
					Role:           "user",
					Level:          data.Invite[i].Level,
					AccountStatus:  "invited",
					AcceptedPolicy: false,
					Timestamp:      time.Now().Unix(),
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}
				db := config.DB
				db.Create(&account)
			}
			//log.Println("reached")
			//generate and store token in database
			tok := token.GenerateToken(data.Invite[i].Email, "onboarding_invite", data.Invite[i].Level)
			payload := database.MailData{
				Email: data.Invite[i].Email,
				Link:  "https://" + config.Conf.Service.HostAddr + "/member-registration?token=" + tok,
				Task:  "inviteLink",
			}
			sendData, err := json.Marshal(payload)
			if err != nil {
				notsent.Email = data.Invite[i].Email
				notsent.Level = data.Invite[i].Level
				notsent.Message = "Error while sending notifications"
				log.Print("error while marshalling data")
				go token.DeleteToken(tok)
				listnotsent = append(listnotsent, notsent)
				continue
			}
			//send request to notification service for sending invite email
			err = nc.Publish(config.Conf.NatsServer.Subject+".notifications.accounts", sendData)
			if err != nil {
				//if there is error sending invite return the error and delete token as token
				//is for one time use only
				notsent.Email = data.Invite[i].Email
				notsent.Level = data.Invite[i].Level
				notsent.Message = "Error while sending notifications"
				go token.DeleteToken(tok)
				log.Print(err)
				listnotsent = append(listnotsent, notsent)
				continue
			}
			listsent = append(listsent, Sents{
				Email:   data.Invite[i].Email,
				Level:   data.Invite[i].Level,
				Message: "Sent Successfully",
			})
		}
	}
	mapd["error"] = false
	mapd["message"] = "Invite Link's Sent"
	mapd["unsuccessful_sents"] = listnotsent
	mapd["successful_sents"] = listsent
	return mapd, 200
}

//=============================== GET ACCOUNT DETAILS FROM EMAIL =============================================

//GetAccountForEmail is used to get account details by passing email
func GetAccountForEmail(email string) (database.Accounts, error) {
	//handler panic and Alerts
	defer util.Panic()
	// connecting to db
	db := config.DB

	// intialize variable with type accounts
	var acs []database.Accounts
	// fetching data on basis of email
	db.Where("email=?", email).Find(&acs)

	// if there is account pass the first element of array
	if len(acs) != 0 {
		return acs[0], nil
	}
	// if there is no account pass null values
	return database.Accounts{}, errors.New("No account found")
}

//============================ GetAccountForID =====================================

//GetAccountForID is used to get account details by passing email
func GetAccountForID(id int) (database.Accounts, error) {
	//handler panic and Alerts
	defer util.Panic()
	// connecting to db
	db := config.DB

	// intialize variable with type accounts
	var acs []database.Accounts
	// fetching data on basis of email
	db.Where("id=?", id).Find(&acs)

	// if there is account pass the first element of array
	if len(acs) != 0 {
		return acs[0], nil
	}
	// if there is no account pass null values
	return database.Accounts{}, errors.New("No account found")
}

//================================ Deactivate Account ===============================================================

//AccountStatus is used to change status of account based on email
func AccountStatus(email, status string) (map[string]interface{}, int) {
	//handler panic and Alerts
	defer util.Panic()
	// connecting to db
	db := config.DB
	mapd := make(map[string]interface{})
	accounts := []database.Accounts{}
	err := db.Where("email=?", email).Find(&accounts)
	if err != nil {
		log.Println(err)
	}
	if len(accounts) == 0 {
		mapd["error"] = true
		mapd["message"] = "No user found"
		return mapd, 404
	}
	if accounts[0].AccountStatus == status {
		mapd["error"] = true
		mapd["message"] = "Account of this email is already " + status
		return mapd, 400
	}
	db.Model(&database.Accounts{}).Where("email=?", email).Update("account_status", status)
	mapd["error"] = false
	mapd["message"] = "Operation successful"
	return mapd, 200
}

//================================= Change Level of Intern =========================================================

//SendChangeLevelEmail is used to send change level emails
type SendChangeLevelEmail struct {
	Email      string `json:"email"`
	Level      string `json:"level"`
	Department string `json:"department"`
}

//ChangeLevel is used to change level and department of interns and notify
//them about that
func ChangeLevel(email, level string, department int) (map[string]interface{}, int) {
	//handler panic and Alerts
	defer util.Panic()
	// connecting to db
	db := config.DB
	mapd := make(map[string]interface{})
	rows := db.Model(&database.Accounts{}).Where("email=?", email).Update("level", level).RowsAffected
	if rows == 0 {
		mapd["error"] = false
		mapd["message"] = "No user found"
		return mapd, 404
	}
	account := database.Accounts{}
	db.Where("email=?", email).Find(&account)
	deptdata := []database.PandCData{}
	err := db.Where("userid=?", account.ID).Find(&deptdata).Error
	if err != nil {
		log.Print(err)
	}
	dept := database.Departments{}
	db.Where("id=?", department).Find(&dept)
	//if details not exist
	if len(deptdata) == 0 {
		var data database.PandCData
		data.Userid = account.ID
		data.Team = dept.Name
		data.Department = dept.ID
		data.UpdatedAt = time.Now()
		data.CreateAt = time.Now()
		log.Println(data)
		db.Create(&data)
	} else {
		deptdata[0].Team = dept.Name
		deptdata[0].Department = dept.ID
		deptdata[0].UpdatedAt = time.Now()
		//if details exist update the database
		err = db.Save(&deptdata[0]).Error
		if err != nil {
			log.Print(err)
		}
	}
	nc := config.NC
	payload := SendChangeLevelEmail{
		Email:      account.Email,
		Level:      level,
		Department: dept.Name,
	}
	sendData, err := json.Marshal(payload)
	if err != nil {
		log.Println(nil)
	}
	err = nc.Publish(config.Conf.NatsServer.Subject+".notifications.changelevel", sendData)
	if err != nil {
		log.Println(nil)
	}
	mapd["error"] = false
	mapd["message"] = "Level updated successfully"
	return mapd, 200
}

//====================================== ACCEPT POLICY =====================================================

//AcceptPolicy is used to store if user has accepted policy
func AcceptPolicy(uid int) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	db := config.DB
	err := db.Model(&database.Accounts{}).Where("id=?", uid).Update("AcceptedPolicy", true).Error
	if err != nil {
		log.Println(err)
		mapd["error"] = true
		mapd["message"] = "Internal server error"
		return mapd, 500
	}
	mapd["error"] = false
	mapd["message"] = "Operation Successfull"
	return mapd, 200
}

//==================================== Check If Account Exists ==========================================

type JwtToken struct {
	Start  string
	End    string
	Expire string
	Error  bool
	Token  string
}

type ResponseData struct {
	Personal Personal `json:"profile"`
}
type Personal struct {
	Fname     string      `json:"fname"`
	Mname     string      `json:"mname"`
	Lname     string      `json:"lname"`
	Email     string      `json:"email"`
	Contact   string      `json:"contact"`
	Password  string      `json:"password"`
	Country   string      `json:"country"`
	State     string      `json:"state"`
	City      string      `json:"city"`
	Postal    string      `json:"postal"`
	Education []Education `json:"education"`
}
type Education struct {
	School string `json:"school"`
	Degree string `json:"degree"`
	End    string `json:"end"`
}

func LoginCarrerPortal() (JwtToken, string) {
	defer util.Panic()
	//log.Println("reached here")
	jwtToken := JwtToken{}
	reqBody, err := json.Marshal(map[string]string{
		"username": config.Conf.CareerPortal.Username,
		"password": config.Conf.CareerPortal.Password,
	})
	if err != nil {
		log.Println(err)
		return jwtToken, err.Error()
	}
	resp, err := http.Post(config.Conf.CareerPortal.FrontEndAddress+"/api/auth/v1/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println(err)
		return jwtToken, err.Error()
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return jwtToken, "Invalid credentials"
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return jwtToken, err.Error()
	}
	//log.Println(string(body))
	err = json.Unmarshal(body, &jwtToken)
	if err != nil {
		log.Println(err)
		return jwtToken, err.Error()
	}
	return jwtToken, ""
}

//CheckIfAccountExistsInCareerPortal is used to check if in career portal account exists if yes then transfer
//already present data here
func CheckIfAccountExistsInCareerPortal(email, level string, jwtToken JwtToken) string {
	defer util.Panic()
	req, err := http.NewRequest("GET", config.Conf.CareerPortal.FrontEndAddress+"/api/auth/v1/users/"+email+"?onboard=true", nil)
	if err != nil {
		log.Println(err)
		return ""
	}
	req.Header.Add("Authorization", "Bearer "+jwtToken.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "false"
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	//log.Println(string(body))
	resData := ResponseData{}
	err = json.Unmarshal(body, &resData)
	if err != nil {
		log.Println(err)
		return ""
	}
	db := config.DB
	account := database.Accounts{
		Email:          email,
		Password:       resData.Personal.Password,
		Name:           resData.Personal.Fname + " " + resData.Personal.Mname + " " + resData.Personal.Lname,
		ContactNo:      resData.Personal.Contact,
		Role:           "user",
		Level:          level,
		AccountStatus:  "active",
		AcceptedPolicy: false,
		Timestamp:      time.Now().Unix(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	//If token is true create account of user
	err = db.Create(&account).Error
	if err != nil {
		log.Println(err)
	}
	postal, err := strconv.Atoi(resData.Personal.Postal)
	if err != nil {
		log.Println(err)
	}
	address := database.Address{
		Userid:      account.ID,
		AddressType: "Permanent address",
		City:        resData.Personal.City,
		State:       resData.Personal.State,
		Pincode:     postal,
		Country:     resData.Personal.Country,
		CreateAt:    time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = db.Create(&address).Error
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < len(resData.Personal.Education); i++ {
		school := strings.Split(resData.Personal.Education[i].School, ",")
		year := strings.Split(resData.Personal.Education[i].End, "-")
		education := database.Education{
			Userid:        account.ID,
			Name:          resData.Personal.Education[i].Degree,
			Institution:   school[0],
			Location:      school[1],
			Year:          year[1],
			MarksheetLink: []string{},
			DegreeLink:    "",
			CreateAt:      time.Now(),
			UpdatedAt:     time.Now(),
		}
		err = db.Create(&education).Error
		if err != nil {
			log.Println(err)
		}
	}
	payload := database.MailData{
		Email: email,
		Link:  "https://" + config.Conf.Service.HostAddr,
		Task:  "inviteCareers",
	}
	sendData, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}
	nc := config.NC
	//send request to notification service for sending invite email
	err = nc.Publish(config.Conf.NatsServer.Subject+".notifications.accounts", sendData)
	if err != nil {
		//if there is error sending invite return the error and delete token as token
		//is for one time use only
		log.Println(err)
	}
	return "true"
}

//=================================== Profile Picture ===================================

//Profile picture function is used to return name as well as img link
func ProfilePicture(email string) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	var link string
	acc, err := GetAccountForEmail(email)
	if err != nil {
		mapd["error"] = true
		mapd["message"] = "Account not found"
		return mapd, 404
	}
	personalInfo, err := GetPersonalInfo(acc.ID)
	if err != nil {
		log.Println(err)
	}
	if personalInfo.ProfileImgLink != "" {
		link, err = SendRequestForLink(personalInfo.ProfileImgLink)
		log.Println(err)
	}
	mapd["error"] = false
	mapd["message"] = "Operation successful"
	mapd["name"] = acc.Name
	mapd["link"] = link
	return mapd, 200
}

//GetInternList is a function to get the list of intern on the basis of status
func GetInternList(limit, page, status string) (error, []database.Accounts, int) {
	pageno, _ := strconv.Atoi(page)
	limits, _ := strconv.Atoi(limit)
	count := 0
	var users []database.Accounts
	err := config.DB.Model(&database.Accounts{}).Where("account_status=?", status).Count(&count).Error
	if err != nil {
		log.Println(err)
		return err, users, count
	}
	err = config.DB.Raw("SELECT * FROM public.accounts where account_status = '" + status + "' ORDER BY created_at desc limit " + limit + " OFFSET " + fmt.Sprint((limits * (pageno - 1)))).Find(&users).Error
	if err != nil {
		log.Println(err)
		return err, users, count
	}
	return nil, users, count
}
