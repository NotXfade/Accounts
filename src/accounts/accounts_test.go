package accounts

import (
	"log"
	"os"
	"testing"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/token"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	os.Remove(os.Getenv("HOME") + "/account-testing.db")
	db, err := gorm.Open("sqlite3", os.Getenv("HOME")+"/account-testing.db")
	if err != nil {
		log.Println(err)
		log.Println("Exit")
		os.Exit(1)
	}
	config.DB = db
	//create table
	database.CreateDatabaseTables()
}

func TestRegisterAccount(t *testing.T) {
	//case 1: when token is correct
	db := config.DB
	account := database.Accounts{
		Email:          "someoneregister@example.com",
		Role:           "user",
		Level:          "L1",
		AccountStatus:  "invited",
		AcceptedPolicy: false,
		Timestamp:      time.Now().Unix(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	//If token is true create account of user
	err := db.Create(&account).Error
	if err != nil {
		log.Println(err)
	}
	tok := token.GenerateToken("someoneregister@example.com", "onboarding_invite", "L1")
	_, code := RegisterAccount(tok, "12345789", "someone", "123456789")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	//case 2: try to use same token again
	_, code = RegisterAccount(tok, "12345789", "someone", "123456789")
	if code != 400 {
		t.Errorf("expected %d , found %d", 400, code)
	}
	//case 3: try to create account with email which already has account
	tok = token.GenerateToken("someoneregister@example.com", "onboarding_invite", "L1")
	_, code = RegisterAccount(tok, "12345789", "someone", "123456789")
	if code != 400 {
		t.Errorf("expected %d , found %d", 400, code)
	}
	//case 4: try to create account by giving empty level value
	tok = token.GenerateToken("someoneregister@example.com", "onboarding_invite", "")
	_, code = RegisterAccount(tok, "12345789", "someone", "123456789")
	if code != 400 {
		t.Errorf("expected %d , found %d", 400, code)
	}
	//case 5: with test token
	config.Conf.Service.Environment = "development"
	_, code = RegisterAccount("XSTestV1.something@example.com", "12345789", "someone", "123456789")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	//case: if email entered is invalid
	_, code = RegisterAccount("XSTestV1.something.example.com", "12345789", "someone", "123456789")
	if code == 200 {
		t.Errorf("expected %d , found %d", 400, code)
	}
}

func TestGetAccountForEmail(t *testing.T) {
	//create account for an email
	_, code := RegisterAccount("XSTestV1.testing@example.com", "12345789", "someone", "123456789")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	//case 1 : when account for an email exists
	_, err := GetAccountForEmail("testing@example.com")
	if err != nil {
		t.Error(err)
	}
	//case 2 : when account for email doesnot exists
	_, err = GetAccountForEmail("user@example.com")
	if err == nil {
		t.Error(err)
	}
}

func TestAccountStatus(t *testing.T) {
	_, code := RegisterAccount("XSTestV1.useraccount@example.com", "12345789", "someone", "123456789")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	//case 1 : when user is active block the account
	_, code = AccountStatus("useraccount@example.com", "blocked")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	//case 2: when account is already blocked
	_, code = AccountStatus("useraccount@example.com", "blocked")
	if code != 400 {
		t.Errorf("expected %d , found %d", 400, code)
	}
	//case 3: Make blocked account active
	_, code = AccountStatus("useraccount@example.com", "active")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}

	//case 4 : when user doesnot exists
	_, code = AccountStatus("someone1@example.com", "active")
	if code != 404 {
		t.Errorf("expected %d , found %d", 404, code)
	}
}

func TestChangeLevel(t *testing.T) {
	_, code := RegisterAccount("XSTestV1.user1@example.com", "12345789", "someone", "123456789")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	//case 1 : when user exists
	_, code = ChangeLevel("user1@example.com", "L2", "Enterprise Application")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	//case 2 : when user doesnot exists
	_, code = AccountStatus("user1@example.com", "blocked")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	_, code = ChangeLevel("user1@example.com", "L1", "Enterprise Application")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	//case 1 : when user exists and pandc data also present
	_, code = ChangeLevel("user1@example.com", "L2", "PlatformOps")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
}

func TestAcceptPolicy(t *testing.T) {
	_, code := RegisterAccount("XSTestV1.user2@example.com", "12345789", "someone", "123456789")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	//case 1 : when user exists
	_, code = AcceptPolicy(2)
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
}

func TestInvite(t *testing.T) {
	//case : Enter wrong email
	inviteData := database.InviteData{
		Email: "test#123",
		Level: "L1",
	}
	dataInvite := database.DataInvite{}
	dataInvite.Invite = append(dataInvite.Invite, inviteData)

	//case : Invalid Level
	inviteData.Email = "test@example.com"
	inviteData.Level = "level 1"
	dataInvite.Invite = append(dataInvite.Invite, inviteData)

	//case: Account Already exists
	inviteData.Email = "invitetest@gmail.com"
	inviteData.Level = "L1"
	dataInvite.Invite = append(dataInvite.Invite, inviteData)

	//case: Invite a valid email whose account does not exist
	inviteData.Email = "invite1test@gmail.com"
	inviteData.Level = "L1"
	dataInvite.Invite = append(dataInvite.Invite, inviteData)
	_, code := Invite(dataInvite)
	if code != 500 {
		t.Errorf("expected %d , found %d", 500, code)
	}
}

func TestCheckIfAccountExistsInCareerPortal(t *testing.T) {
	config.Conf.CareerPortal.FrontEndAddress = "http://test.example"
	isExist := CheckIfAccountExistsInCareerPortal("xyz@example.com", "L1", JwtToken{})
	if isExist != "" {
		t.Error("Did not return error")
	}
}

func TestLoginCarrerPortal(t *testing.T) {
	config.Conf.CareerPortal.FrontEndAddress = "http://test.example"
	config.Conf.CareerPortal.Username = "username"
	config.Conf.CareerPortal.Password = "password"
	_, err := LoginCarrerPortal()
	if err == "" {
		t.Error("Did not return error")
	}
}

func TestGetProfilePicture(t *testing.T) {
	_, code := ProfilePicture("profileuser@gmail.com")
	if code != 404 {
		t.Errorf("expected %d , found %d", 404, code)
	}
}
