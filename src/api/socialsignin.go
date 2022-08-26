package api

import (
	"log"
	"strings"

	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/googlecalendar"
	"github.com/gin-gonic/gin"
)

// GoogleLogin is an API handler to redirect the service to google.
func GoogleLogin(c *gin.Context) {
	// fetch and set the redriect
	redirect := c.Query("redirect")
	if redirect == "" {
		redirect = c.Request.Host
		if strings.Contains(redirect, ":") {
			redirect = "http://" + redirect
		} else {
			redirect = "https://" + redirect
		}

	}

	url := googlecalendar.GoogleLogin(googlecalendar.State{Redirect: redirect})
	c.Redirect(302, url)
}

// GoogleCallback is an API handler for call back of google signin
func GoogleCallback(c *gin.Context) {

	// log.Println(c.Query("code"), c.Query("state"))
	url := googlecalendar.GoogleCallback(c.Query("code"), c.Query("state"))
	c.Redirect(302, url)
}

//CalendarInvite is used to create calendar invite using google apis
func CalendarInvite(c *gin.Context) {
	var inviteData database.InviteEmail
	if err := c.BindJSON(&inviteData); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "Please enter required fields"})
		return
	}
	if inviteData.EndTime-inviteData.StartTime <= 0 {
		c.JSON(400, gin.H{"error": true, "message": "Please enter valid timings"})
		return
	}
	mapd, code := googlecalendar.CreateInvite(inviteData)
	c.JSON(code, mapd)
}

//CheckCredentials is used to check if google token is present in db
func CheckCredentials(c *gin.Context) {
	mapd, code := googlecalendar.CheckCredentials()
	c.JSON(code, mapd)
}
