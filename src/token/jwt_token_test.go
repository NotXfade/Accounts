package token

import (
	"log"
	"os"
	"testing"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

const (
	email string = "test@test.com"
	name  string = "test"
)

var (
	id int
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
	acc := database.Accounts{}
	//save data in account table
	acc.Name = name
	acc.Email = email
	db.Create(&acc)
	id = acc.ID
}

func TestGinJwtToken(t *testing.T) {

	//create table
	database.CreateDatabaseTables()
	acc := database.Accounts{}
	mapd := make(map[string]interface{})
	//save data in account table
	mapd["name"] = name
	mapd["role"] = "user"
	mapd["email"] = email
	mapd["id"] = acc.ID
	acc.ID = id
	acc.Email = email
	acc.Name = name
	acc.Role = "user"
	acc.Level = "L1"
	_ = GenerateJwtToken(acc)
	//test csae 1 -> generate token with account only
	mapd, _ = GinJwtToken(mapd)
	if mapd["token"].(string) == "" {
		t.Error("test case fail")
	}

	//test case 2 -> Refresh JWT token
	mapd = JwtRefreshToken(mapd)
	if mapd["token"].(string) == "" {
		t.Error("test case fail")
	}

	//test case 4 -> Delete the token
	DeleteTokenFromDb(mapd["token"].(string))
}
