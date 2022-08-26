package api

import (
	"log"

	"git.xenonstack.com/xs-onboarding/accounts/src/accounts"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
	"github.com/gin-gonic/gin"
)

//Assignreviewer is used for binding data
type Assignreviewer struct {
	Reviewers []string `json:"reviewers" binding:"required"`
	Interns   []string `json:"interns" binding:"required"`
}

//AssignResponse is used to return response
type AssignResponse struct {
	Email     string     `json:"email"`
	Reviewers []response `json:"status"`
}

//response is used to send response data
type response struct {
	Reviewer string `json:"reviewer"`
	Message  string `json:"message"`
}

//AssignReviewer is used to assign reviewers to interns
func AssignReviewer(c *gin.Context) {
	defer util.Panic()
	assignReviewer := Assignreviewer{}
	if err := c.BindJSON(&assignReviewer); err != nil {
		// if there is some error passing bad status code
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "Please select atleast one intern and one reviewer"})
		return
	}
	if len(assignReviewer.Reviewers) == 0 {
		c.JSON(400, gin.H{"error": true, "message": "Please select atleast one reviewer"})
		return
	}
	if len(assignReviewer.Interns) == 0 {
		c.JSON(400, gin.H{"error": true, "message": "Please select atleast one intern"})
		return
	}
	successful := []AssignResponse{}
	unsuccessful := []AssignResponse{}

	for i := 0; i < len(assignReviewer.Interns); i++ {
		unsuccess := []response{}
		success := []response{}
		for j := 0; j < len(assignReviewer.Reviewers); j++ {
			//Call function to store data in DB
			err := accounts.AssignReviewer(assignReviewer.Interns[i], assignReviewer.Reviewers[j])
			if err != nil {
				unsuccess = append(unsuccess, response{
					Reviewer: assignReviewer.Reviewers[j],
					Message:  err.Error(),
				})
				continue
			}
			success = append(success, response{
				Reviewer: assignReviewer.Reviewers[j],
				Message:  "Assigned successfully",
			})

		}
		if len(success) != 0 {
			successful = append(successful, AssignResponse{
				Email:     assignReviewer.Interns[i],
				Reviewers: success,
			})
		}
		if len(unsuccess) != 0 {
			unsuccessful = append(unsuccessful, AssignResponse{
				Email:     assignReviewer.Interns[i],
				Reviewers: unsuccess,
			})
		}

	}

	c.JSON(200, gin.H{
		"error":        "false",
		"message":      "Operation successful",
		"successful":   successful,
		"unsuccessful": unsuccessful,
	})
}

//GetListOfAssignedReviewers is used to get reviewers list that are assigned to users
func GetListOfAssignedReviewers(c *gin.Context) {
	defer util.Panic()
	email := c.Param("email")
	mapd, code := accounts.GetListOfAssignedReviewers(email)
	c.JSON(code, mapd)
}

//DeleteUserAssignedReviewer is used to delete user assigned reviewer
func DeleteUserAssignedReviewer(c *gin.Context) {
	defer util.Panic()
	email := c.Param("email")
	acc, _ := accounts.GetAccountForEmail(email)
	uemail := c.Param("uemail")
	useracc, err := accounts.GetAccountForEmail(uemail)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	err = accounts.DeleteAssignedReviewer(useracc.ID, acc.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   true,
		"message": "Reviewer deleted successfully",
	})
}
