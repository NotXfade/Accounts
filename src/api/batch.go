package api

import (
	"git.xenonstack.com/xs-onboarding/accounts/src/accounts"
	"github.com/gin-gonic/gin"
)

//GetBatch is used to get unique list of batch
func GetBatch(c *gin.Context) {
	mapd, code := accounts.BatchList()
	c.JSON(code, mapd)
}

//GetReports is used to get reports list of interns
func GetReports(c *gin.Context) {
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

	sortby := c.Query("sortby")
	batch := c.Query("batch")
	module := c.Query("module")
	mapd, code := accounts.GetReports(level, batch, module, sortby, limit, page)
	c.JSON(code, mapd)
}
