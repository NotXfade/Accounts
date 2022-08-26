// +build !test

package scheduler

import (
	"log"
	"strconv"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"github.com/robfig/cron"
)

// Start is a function to start cronjobs
func Start() {
	log.Println("running database cleaner scheduler")
	DeleteUsers()
	c := cron.New()
	c.AddFunc("0 30 * * *", DeleteUsers)
	c.Start()
}

//DeleteUsers is used to delete tokens
func DeleteUsers() {
	db := config.DB
	db.Exec("delete from accounts where account_status='invited' AND timestamp<" + strconv.FormatInt(time.Now().Unix()-604800, 10) + ";")
	db.Exec("delete from tokens where token_task='onboarding_invite' AND timestamp<" + strconv.FormatInt(time.Now().Unix()-604800, 10) + ";")
	db.Exec("delete from tokens where token_task='adminInvite' AND timestamp<" + strconv.FormatInt(time.Now().Unix()-604800, 10) + ";")
	db.Exec("delete from tokens where token_task='reset-password' AND timestamp<" + strconv.FormatInt(time.Now().Unix()-86400, 10) + ";")
}
