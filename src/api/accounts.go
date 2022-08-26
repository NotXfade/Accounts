//+build !test

package api

import (
	"log"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/accounts"
	"git.xenonstack.com/xs-onboarding/accounts/src/activity"
	"git.xenonstack.com/xs-onboarding/accounts/src/methods"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
)

//===================================== CREATE ACCOUNT ============================================

// TokenPassword is a structure for binding data from body during set new password request
type TokenPassword struct {
	Name     string `json:"name" binding:"required"`
	Contact  string `json:"contact" binding:"required"`
	Password string `json:"password" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

//Register : function register is used to register account of user
func Register(c *gin.Context) {
	//handler panic and Alerts
	defer util.Panic()
	var tp TokenPassword
	if err := c.BindJSON(&tp); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "Name, Contact, Password, Token are required field."})
		return
	}
	// check password is valid
	if !methods.CheckPassword(tp.Password) {
		c.JSON(400, gin.H{"error": true, "message": "Minimum eight characters, at least one uppercase letter, at least one lowercase letter, at least one number and at least one special character."})
		return
	}
	// update password in database
	mapd, code := accounts.RegisterAccount(tp.Token, methods.HashForNewPassword(tp.Password), tp.Name, tp.Contact)
	if code == 200 {
		// recording user activity of reseting password
		activity.RecordActivity(database.Activities{Email: mapd["email"].(string),
			ActivityName: "registration",
			ClientIP:     c.ClientIP(),
			ClientAgent:  c.Request.Header.Get("User-Agent"),
			Timestamp:    time.Now().Unix()})
	}
	c.JSON(code, mapd)
}

//===================================== INVITE FOR ACCOUNT CREATION =========================================

//Invite function is used to send invite only signup link to new users
func Invite(c *gin.Context) {
	//handler panic and Alerts
	defer util.Panic()
	var data database.DataInvite
	if err := c.BindJSON(&data); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "Email and Level are required fields"})
		return
	}
	mapd, code := accounts.Invite(data)
	c.JSON(code, mapd)
}

//====================================== ACTIVATE OR DEACTIVATE ACCOUNT ================================================

//AccountStatus is used to activate or deactivate account
func AccountStatus(c *gin.Context) {
	defer util.Panic()
	email := c.Param("email")
	status := c.Param("status")
	if status != "active" && status != "blocked" {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass correct status values, Values can be active or blocked",
		})
		return
	}
	mapd, code := accounts.AccountStatus(email, status)
	c.JSON(code, mapd)
}

//======================================= CHANGE LEVEL ===============================================

//LevelChange struct for binding data
type LevelChange struct {
	Email      string `json:"email" binding:"required"`
	Level      string `json:"level" binding:"required"`
	Department int    `json:"department"`
}

//ChangeLevel to call function to change level of employee
func ChangeLevel(c *gin.Context) {
	defer util.Panic()
	var data LevelChange
	if err := c.BindJSON(&data); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "level is required field"})
		return
	}
	Level := data.Level
	if Level != "L1" && Level != "L2" && Level != "L3" {
		c.JSON(400, gin.H{"error": true, "message": "Please pass correct level value"})
		return
	}
	mapd, code := accounts.ChangeLevel(data.Email, data.Level, data.Department)
	c.JSON(code, mapd)
}

//=============================== Accept Policy =======================================

//AcceptPolicy is used to update that user has accepted policy
func AcceptPolicy(c *gin.Context) {
	defer util.Panic()
	claims := jwt.ExtractClaims(c)
	id := int(claims["id"].(float64))
	mapd, code := accounts.AcceptPolicy(id)
	c.JSON(code, mapd)
}

//=============================== LIST INTERNS =======================================

//InternsList is use to list employees details based on level and limit offset
func InternsList(c *gin.Context) {
	defer util.Panic()
	level := c.Param("level")
	if level == "" {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass level of employee",
		})
	}
	//limit is used to add limit for pagination
	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}

	//page is pagenumber that is offser
	page := c.Query("page")
	if page == "" {
		page = "1"
	}

	claims := jwt.ExtractClaims(c)
	role := claims["role"].(string)
	email := claims["email"].(string)
	if role == "reviewer" {
		mapd, code := accounts.GetListOfAssignedInterns(email, level, limit, page)
		c.JSON(code, mapd)
		return
	}

	progressValue := c.Query("progress")
	var progress float64
	var err error
	if progressValue != "" {
		progress, err = strconv.ParseFloat(progressValue, 64)
		if err != nil {
			log.Println(err)
		}
	} else {
		progress = -1
	}
	//sort will be used in case of sorting
	sort := c.Query("sort")
	search := c.Query("search")
	//log.Println(sort)
	mapd, code := accounts.GetListInterns(level, limit, page, sort, search, progress)
	c.JSON(code, mapd)
}

// //InternsList is use to list employees details based on level and limit offset
// func InternsList(c *gin.Context) {
// 	defer util.Panic()
// 	level := c.Param("level")
// 	if level == "" {
// 		c.JSON(400, gin.H{
// 			"error":   true,
// 			"message": "Please pass level of employee",
// 		})
// 	}
// 	//limit is used to add limit for pagination
// 	limit := c.Query("limit")
// 	if limit == "" {
// 		limit = "10"
// 	}

// 	//page is pagenumber that is offser
// 	page := c.Query("page")
// 	if page == "" {
// 		page = "1"
// 	}

// 	claims := jwt.ExtractClaims(c)
// 	role := claims["role"].(string)
// 	rid := claims["id"]
// 	reviewerid := int(rid.(float64))
// 	if role != "reviewer" {
// 		reviewerid = 0
// 	}
// 	id := strconv.Itoa(reviewerid)
// 	//-- to be implemented
// 	//sort := c.Query("sort")
// 	search := c.Query("search")
// 	module := c.Query("module")
// 	batch := c.Query("batch")
// 	department := c.Query("department")
// 	score := c.Query("score")
// 	mapd, code := accounts.GetListInterns(level, limit, page, search, module)
// 	c.JSON(code, mapd)
// }

//=============================== VIEW USER PROFILE for admin ======================================

//ViewUserProfile is used to get details of user
func ViewUserProfile(c *gin.Context) {
	defer util.Panic()
	email := c.Param("email")
	mapd, code := accounts.GetUserInfo(email)
	c.JSON(code, mapd)
}

//============================= VIEW PROFILE ==============================================

//ViewProfile is used to viewing profile by user
func ViewProfile(c *gin.Context) {
	defer util.Panic()
	claims := jwt.ExtractClaims(c)
	email := claims["email"].(string)
	mapd, code := accounts.UserInfo(email)
	c.JSON(code, mapd)
}

//============================ DELETE ACCOUNT =============================================

//Delete Account is used for deleting account
func Delete(c *gin.Context) {
	defer util.Panic()
	email := c.Param("email")
	log.Println("reached here")
	mapd, code := accounts.Delete(email)
	c.JSON(code, mapd)
}

//=========================== GetProfilePicture ==========================================================

//GetProfilePicture is used to get profile picture
func GetProfilePicture(c *gin.Context) {
	defer util.Panic()
	claims := jwt.ExtractClaims(c)
	email := claims["email"].(string)
	mapd, code := accounts.ProfilePicture(email)
	c.JSON(code, mapd)
}

//InternList is used to list out invited users
func InternList(c *gin.Context) {

	page := c.Query("page")
	if c.Query("page") == "" {
		page = "1"
	}
	err, users, count := accounts.GetInternList(c.Query("limit"), page, c.Query("status"))
	if err != nil {
		c.JSON(400, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":       false,
		"total_users": count,
		"user_list":   users,
	})
}
