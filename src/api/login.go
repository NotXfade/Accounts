//+build !test

package api

import (
	"log"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/src/activity"

	"git.xenonstack.com/xs-onboarding/accounts/src/methods"

	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/accounts"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
	"github.com/gin-gonic/gin"
)

//=========================================== LOGIN =======================================================

//Login is used in LoginEndPoint function to bind incoming data
type Login struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

//LoginEndPoint : func Login is used to perform login operations in portal
func LoginEndPoint(c *gin.Context) {
	//handler panic and Alerts
	defer util.Panic()
	//log.Println(c)
	var loginData Login
	if err := c.BindJSON(&loginData); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "email and password are required field."})
		return
	}
	if !methods.ValidateEmail(loginData.Email) {
		c.JSON(400, gin.H{"error": true, "message": "Please enter a valid email"})
		return
	}
	//=============================================
	// recording user activity
	activityData := database.Activities{Email: loginData.Email,
		ClientIP:    c.ClientIP(),
		CreatedAt:   time.Now(),
		ClientAgent: c.Request.Header.Get("User-Agent"),
		Timestamp:   time.Now().Unix()}
	//=============================================
	mapd, code := accounts.LoginEndPoint(loginData.Email, loginData.Password)
	if code == 200 {
		// recording user activity of login
		activityData.ActivityName = "login"
		activityData.CreatedAt = time.Now()
		activity.RecordActivity(activityData)
	} else {
		// recording user activity of failed login
		activityData.ActivityName = "failedlogin"
		activityData.CreatedAt = time.Now()
		activity.RecordActivity(activityData)
	}
	c.JSON(code, mapd)
}
