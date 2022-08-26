package api

import (
	"log"
	"strconv"

	"git.xenonstack.com/xs-onboarding/accounts/src/accounts"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

//===================================== Create Department ============================================

// createDepartment is a structure for binding data from body during creating new department
type createDepartment struct {
	Name      string `json:"name" binding:"required"`
	ShortName string `json:"short_name" `
}

//CreateDepartment is used for creating new department
func CreateDepartment(c *gin.Context) {
	//handler panic and Alerts
	defer util.Panic()
	var dept createDepartment
	if err := c.BindJSON(&dept); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "Please enter name of department"})
		return
	}
	claims := jwt.ExtractClaims(c)
	email := claims["email"].(string)
	mapd, code := accounts.CreateDepartment(dept.Name, dept.ShortName, email)
	c.JSON(code, mapd)
}

//UpdateDepartment is using for updating departments
func UpdateDepartment(c *gin.Context) {
	//handler panic and Alerts
	defer util.Panic()
	var dept createDepartment
	if err := c.BindJSON(&dept); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "Please enter name of department"})
		return
	}
	deptid := c.Param("id")
	id, _ := strconv.Atoi(deptid)
	mapd, code := accounts.UpdateDepartment(dept.Name, dept.ShortName, id)
	c.JSON(code, mapd)
}

//ListDepartment is used for listing departments
func ListDepartment(c *gin.Context) {
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
	mapd, code := accounts.ListDepartment(limit, page)
	c.JSON(code, mapd)
}

//DeleteDepartment is used for deleting departments
func DeleteDepartment(c *gin.Context) {
	deptid := c.Param("id")
	id, _ := strconv.Atoi(deptid)
	mapd, code := accounts.DeleteDepartment(id)
	c.JSON(code, mapd)
}
