//+build !test

package api

import (
	"log"

	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/accounts"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

//============================ STORE EMPLOYEE INFO ====================================

//StoreEmployeeInfo is used to store information of employee
//related to education documents and other relevant information
func StoreEmployeeInfo(c *gin.Context) {
	//handler panic and Alerts
	defer util.Panic()
	claims := jwt.ExtractClaims(c)
	//log.Print(claims)
	var empData database.EmployeeInfo
	if err := c.BindJSON(&empData); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "Please enter all the required fields."})
		return
	}
	id := int(claims["id"].(float64))
	//log.Print(id)
	mapd, code := accounts.StoreEmployeeInfo(empData, id)
	c.JSON(code, mapd)
}

//============================ UPDATE USER PROFILE ===========================

//UpdateProfile api is used to update profile
func UpdateProfile(c *gin.Context) {
	//handler panic and Alerts
	defer util.Panic()
	claims := jwt.ExtractClaims(c)
	id := int(claims["id"].(float64))
	var empData database.UpdateEmployeeInfo
	if err := c.BindJSON(&empData); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "Please enter all the required fields."})
		return
	}
	mapd, code := accounts.UpdateEmployeeInfo(empData, id)
	c.JSON(code, mapd)
}
