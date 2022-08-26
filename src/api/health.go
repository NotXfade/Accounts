//+build !test

package api

import (
	"log"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/src/health"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
	"github.com/gin-gonic/gin"
)

// Healthz is an api handler to check health of service
func Healthz(c *gin.Context) {

	//handler panic and Alerts
	defer util.Panic()

	// call health service check function
	err := health.ServiceHealth()
	if err != nil {
		// if any error is there
		log.Println(err)
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// if no error is there
	c.JSON(200, gin.H{
		"error":       false,
		"message":     "All is okay",
		"build":       config.Conf.Service.Build,
		"environment": config.Conf.Service.Environment,
	})
}
