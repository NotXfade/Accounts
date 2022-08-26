package database

import (
	"log"
	"os"
	"testing"

	"git.xenonstack.com/xs-onboarding/accounts/config"
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

}
func TestDatabase(t *testing.T) {

	CreateDatabase()

	CreateDatabaseTables()

}
