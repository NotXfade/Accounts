//+build !test

package api

import (
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt"

	"git.xenonstack.com/xs-onboarding/accounts/src/accounts"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
	"github.com/gin-gonic/gin"
)

//Export is used to return csv file
func Export(c *gin.Context) {
	defer util.Panic()
	level := c.Query("level")
	claims := jwt.ExtractClaims(c)
	role := claims["role"].(string)
	if role != "business analyst" && role != "p&c" && role != "admin" {
		c.Abort()
		c.JSON(401, gin.H{"message": "You are not authorised."})
		return
	}
	data, code, err := accounts.ExportData(level)
	//log.Println(string(data), code, err)
	if code == 200 {
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename="+"internsList_"+strconv.FormatInt(time.Now().Unix(), 10)+".csv")
		c.Data(code, "text/csv", data)
		return
	}
	c.JSON(code, gin.H{
		"error":   true,
		"message": err,
	})
}
