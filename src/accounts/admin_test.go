package accounts

import (
	"log"
	"testing"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/token"
)

func TestInviteAdmin(t *testing.T) {
	inviteData := database.Inviteadmin{
		//case 1 : enter wrong role
		Invite: []database.AdminInvite{
			database.AdminInvite{Email: "admintest@example.com",
				Role: "L1",
			},
			//case 2 : enter wrong email
			database.AdminInvite{
				Email: "someone",
				Role:  "P&C",
			},
			//case 3 : when email already has a account
			database.AdminInvite{
				Email: "someone@",
				Role:  "P&C",
			},
			//case 4 : email with capital casing
			database.AdminInvite{
				Email: "Admintest@example.com",
				Role:  "P&C",
			},
		},
	}
	_, code := InviteAdmin(inviteData)
	if code != 200 {
		t.Error("failed on invite admin")
	}
}

func TestRegisterAdmin(t *testing.T) {
	db := config.DB
	account := database.Accounts{
		Email:          "adminregister@example.com",
		Role:           "p&c",
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
	tok := token.GenerateToken("adminregister@example.com", "adminInvite", "p&c")
	registerAdmin := database.Registeradmin{
		Name:     "xyz",
		Contact:  "9876543210",
		Password: "Test@1234",
		Token:    tok,
	}
	_, code := RegisterAdmin(registerAdmin)
	if code != 200 {
		t.Error("expected", 200, " found ", code)
	}
	//reuse token again
	_, code = RegisterAdmin(registerAdmin)
	if code != 400 {
		t.Error("expected", 400, " found ", code)
	}
}

func TestListAdmin(t *testing.T) {
	//case 1 : when role is empty
	_, code := ListAdmin("", "10", "1")
	if code != 200 {
		t.Error("expected", 200, " found ", code)
	}

	//case 2: when role is non empty
	_, code = ListAdmin("P&C", "10", "1")
	if code != 200 {
		t.Error("expected", 200, " found ", code)
	}
}

func TestDeleteAdmin(t *testing.T) {
	db := config.DB
	account := database.Accounts{
		Email:          "admindelete@example.com",
		Role:           "p&c",
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
	tok := token.GenerateToken("admindelete@example.com", "adminInvite", "p&c")
	registerAdmin := database.Registeradmin{
		Name:     "xyz",
		Contact:  "9876543210",
		Password: "Test@1234",
		Token:    tok,
	}
	_, code := RegisterAdmin(registerAdmin)
	if code != 200 {
		t.Error("expected", 400, " found ", code)
	}
	//case 1: when account doesn't exist
	_, code = DeleteAdmin("deltest@gmail.com")
	if code != 404 {
		t.Error("expected", 404, " found ", code)
	}

	//case 2: when account exist
	_, code = DeleteAdmin("admindelete@example.com")
	if code != 200 {
		t.Error("expected", 200, " found ", code)
	}
}
