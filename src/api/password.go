//+build !test

package api

import (
	"log"
	"strings"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/src/token"

	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/accounts"
	"git.xenonstack.com/xs-onboarding/accounts/src/activity"

	"git.xenonstack.com/xs-onboarding/accounts/src/methods"

	"git.xenonstack.com/xs-onboarding/accounts/src/util"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

//==================== RESET PASSWORD TO SEND RESET LINK IN EMAIL FOR FORGOT/RESET PASSWORD ======================

// ResetPasswordData is a  structure for binding data in body during forget or reset password request
type ResetPasswordData struct {
	// email of user
	Email string `json:"email"`
}

// ResetPassword api is used to send email with reset password link
func ResetPassword(c *gin.Context) {
	//Handle panics and alerts
	defer util.Panic()
	var resetData ResetPasswordData
	if err := c.BindJSON(&resetData); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "Email is required field."})
		return
	}
	//Check whether entered email is valid
	isValid := methods.ValidateEmail(resetData.Email)
	if !isValid {
		log.Println("email not valid")
		c.JSON(400, gin.H{"error": true, "message": "Please enter a valid email"})
		return
	}
	//Send Email for Reset of Password
	mapd, code := accounts.ForgotPassword(strings.ToLower(resetData.Email))
	if code == 200 {
		// recording user activity of reseting password
		activity.RecordActivity(database.Activities{Email: resetData.Email,
			ActivityName: "reset-password",
			ClientIP:     c.ClientIP(),
			ClientAgent:  c.Request.Header.Get("User-Agent"),
			Timestamp:    time.Now().Unix()})
	}
	c.JSON(code, mapd)
}

//========================== SET PASSWORD FROM EMAIL TOKEN =====================================

//NewPassword struct is used to bind data when set password api is hit
type NewPassword struct {
	Password string `json:"password" binding:"required"`
	Token    string `json:"token"`
}

//SetPassword is used to set password from token that was sent in mail for forget password
func SetPassword(c *gin.Context) {
	//Handle panics and alerts
	defer util.Panic()
	var newPassData NewPassword
	if err := c.BindJSON(&newPassData); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "New password and Reset Token are required fields"})
		return
	}
	tok, err := token.VerifyToken(newPassData.Token, "reset-password")
	if err != nil {
		c.JSON(400, gin.H{"error": true, "message": "Invalid Token or Your Token has expired"})
		return
	}
	if !methods.CheckPassword(newPassData.Password) {
		c.JSON(400, gin.H{"error": true, "message": "Minimum eight characters, at least one uppercase letter, at least one lowercase letter, at least one number and at least one special character."})
		return
	}
	mapd, code := accounts.SetNewPassword(tok.Email, newPassData.Password)
	if code == 200 {
		token.DeleteToken(tok.Token)
		// recording user activity of reseting password
		activity.RecordActivity(database.Activities{Email: tok.Email,
			ActivityName: "set-password",
			ClientIP:     c.ClientIP(),
			ClientAgent:  c.Request.Header.Get("User-Agent"),
			Timestamp:    time.Now().Unix()})
	}
	c.JSON(code, mapd)
}

type Changepassword struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

//ChangePassword is used to change password of user
func ChangePassword(c *gin.Context) {
	defer util.Panic()
	passwordData := Changepassword{}
	claims := jwt.ExtractClaims(c)
	email := claims["email"].(string)
	if err := c.BindJSON(&passwordData); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "New password and Old password are required fields"})
		return
	}
	if !methods.CheckPassword(passwordData.NewPassword) {
		c.JSON(400, gin.H{"error": true, "message": "Minimum eight characters, at least one uppercase letter, at least one lowercase letter, at least one number and at least one special character."})
		return
	}
	if !methods.CheckPassword(passwordData.OldPassword) {
		c.JSON(400, gin.H{"error": true, "message": "Please enter correct current password"})
		return
	}
	mapd, code := accounts.ChangePassword(passwordData.OldPassword, passwordData.NewPassword, email)
	c.JSON(code, mapd)
}
