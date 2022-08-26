package api

import (
	"log"

	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/accounts"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

//InviteAdmin is used to invite admin
func InviteAdmin(c *gin.Context) {
	//handler panic and Alerts
	defer util.Panic()
	var data database.Inviteadmin
	if err := c.BindJSON(&data); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "Email and Role are required fields"})
		return
	}
	mapd, code := accounts.InviteAdmin(data)
	c.JSON(code, mapd)
}

//RegisterAdmin is used to register other roles than user
func RegisterAdmin(c *gin.Context) {
	//handler panic and Alerts
	defer util.Panic()
	var data database.Registeradmin
	if err := c.BindJSON(&data); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "Please enter all the required fields"})
		return
	}
	mapd, code := accounts.RegisterAdmin(data)
	c.JSON(code, mapd)
}

//ListAdmin is used to list admin
func ListAdmin(c *gin.Context) {
	//handler panic and Alerts
	defer util.Panic()
	//limit is used to add limit for pagination
	limit := c.Query("limit")
	if limit == "" {
		limit = "20"
	}

	//page is pagenumber that is offser
	page := c.Query("page")
	if page == "" {
		page = "1"
	}

	role := c.Query("role")
	mapd, code := accounts.ListAdmin(role, limit, page)
	c.JSON(code, mapd)
}

//DeleteAdmin account is used for deleting admin account
func DeleteAdmin(c *gin.Context) {
	//handler panic and Alerts
	defer util.Panic()
	email := c.Param("email")
	claims := jwt.ExtractClaims(c)
	if claims["email"] == email {
		c.JSON(200, gin.H{
			"error":   true,
			"message": "You cannot delete your own account",
		})
	}
	mapd, code := accounts.DeleteAdmin(email)
	c.JSON(code, mapd)
}
