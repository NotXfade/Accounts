package activity

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
	log.Println(os.Getenv("HOME"))
	db, err := gorm.Open("sqlite3", os.Getenv("HOME")+"/account-testing.db")
	if err != nil {

		log.Println(err)
		log.Println("Exit")
		os.Exit(1)
	}
	config.DB = db

	database.CreateDatabaseTables()
	acc := database.Accounts{}

	acc.Name = name
	acc.Email = email
	db.Create(&acc)
	id = acc.ID

}

func TestRecordActivity(t *testing.T) {
	data := database.Activities{}
	data.Email = email
	data.ActivityName = name
	RecordActivity(data)
}

func TestGetLoginActivities(t *testing.T) {

	data, err := GetLoginActivities(email)
	if err != nil {
		t.Error("test case fail", data)
	}
}
