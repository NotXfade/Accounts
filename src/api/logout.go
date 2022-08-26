//+build !test

package api

import (
	"strings"

	"git.xenonstack.com/xs-onboarding/accounts/src/token"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
	"github.com/gin-gonic/gin"
)

//===================================== LOGOUT =========================================

//Logout : It is used to Logout of portal and delete session
func Logout(c *gin.Context) {
	//handler panic and Alerts
	defer util.Panic()

	tok := c.Request.Header.Get("Authorization")
	// trim bearer from token
	tok = strings.TrimPrefix(tok, "Bearer ")
	token.DeleteTokenFromDb(tok)
	c.JSON(200, gin.H{
		"error":   false,
		"message": "Logged Out Successfully",
	})
}
