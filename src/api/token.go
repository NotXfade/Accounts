//+build !test

package api

import (
	"git.xenonstack.com/xs-onboarding/accounts/src/token"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

//RefreshToken is used to refresh JWT Token
func RefreshToken(c *gin.Context) {
	defer util.Panic()
	claims := jwt.ExtractClaims(c)
	// generating new token
	mapd := token.JwtRefreshToken(claims)
	c.JSON(200, mapd)
}
