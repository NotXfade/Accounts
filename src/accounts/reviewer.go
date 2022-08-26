package accounts

import (
	"errors"
	"log"
	"strconv"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
)

//AssignReviewer is used to assign reviewer to users
func AssignReviewer(email string, reviewer string) error {
	defer util.Panic()
	db := config.DB
	acc, err := GetAccountForEmail(email)
	if err != nil {
		return errors.New("User does not exist")
	}
	reviewerDetails, err := GetAccountForEmail(reviewer)
	if err != nil {
		log.Println("error")
		return errors.New("Reviewer does not exist")
	}
	reviewers := []database.Reviewer{}
	db.Where("reviewerid=? AND userid=?", reviewerDetails.ID, acc.ID).Find(&reviewers)
	if len(reviewers) != 0 {
		log.Println("error")
		return errors.New("Already assigned")
	}
	assignReviewer := database.Reviewer{
		Userid:     acc.ID,
		Reviewerid: reviewerDetails.ID,
		CreateAt:   time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = db.Create(&assignReviewer).Error
	if err != nil {
		log.Println(err)
	}
	return nil
}

type reviewerResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

//GetListOfAssignedReviewers is used to list reviewers assigned to a user
func GetListOfAssignedReviewers(email string) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	acc, err := GetAccountForEmail(email)
	if err != nil {
		mapd["error"] = true
		mapd["message"] = "User account not found"
		return mapd, 404
	}
	reviewer, _ := GetReviewer(acc.ID)
	if len(reviewer) == 0 {
		mapd["error"] = true
		mapd["message"] = "No reviewers assigned yet"
		return mapd, 404
	}
	response := []reviewerResponse{}
	for i := 0; i < len(reviewer); i++ {
		account, err := GetAccountForID(reviewer[i].Reviewerid)
		if err != nil {
			log.Println(err)
		}
		response = append(response, reviewerResponse{
			Name:  account.Name,
			Email: account.Email,
		})
	}
	mapd["error"] = false
	mapd["message"] = "Operation Successful"
	mapd["list"] = response
	return mapd, 200
}

//GetListOfAssignedInterns is used to get interns list
func GetListOfAssignedInterns(email, level, limit, page string) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	//log.Println("reached here")
	list := []InternsList{}
	accounts := []database.Accounts{}
	limits, _ := strconv.Atoi(limit)
	pageno, _ := strconv.Atoi(page)
	offset := (limits * (pageno - 1))
	db := config.DB
	count := Count{}

	acc, err := GetAccountForEmail(email)
	if err != nil {
		mapd["error"] = true
		mapd["message"] = err.Error()
		return mapd, 404
	}
	queryCount := "SELECT COUNT(email) FROM reviewerslist WHERE reviewerid =" + strconv.Itoa(acc.ID) + " AND account_status='active' AND level='" + level + "'"
	err = db.Raw(queryCount).Scan(&count).Error
	if err != nil {
		log.Println(err)
	}

	query := "select * from reviewerslist where reviewerid =" + strconv.Itoa(acc.ID) + " AND account_status='active' AND level='" + level + "'"
	query = query + " ORDER BY lower(name) " + "LIMIT " + limit + " OFFSET " + strconv.Itoa(offset) + ";"
	err = db.Raw(query).Scan(&accounts).Error
	if err != nil {
		log.Println(err)
	}
	//log.Println(accounts)
	for i := 0; i < len(accounts); i++ {
		var interndata InternsList
		interndata.Email = accounts[i].Email
		interndata.Name = accounts[i].Name
		internsList := SendRequestForStatus(accounts[i].Email)
		interndata.Task = internsList.Task
		interndata.Status = internsList.Status
		interndata.Progress = internsList.Progress
		pandc := database.PandCData{}
		db.Where("userid=?", accounts[i].ID).Find(&pandc)
		interndata.Department = pandc.Team
		interndata.Scoring = internsList.Scoring
		personal := []database.PersonalInfo{}
		//log.Println(i)
		db.Where("userid=?", accounts[i].ID).Find(&personal)
		//log.Println(personal)
		if len(personal) == 0 {
			list = append(list, interndata)
			//log.Println(i)
			continue
		}
		if personal[0].ProfileImgLink != "" {
			link, err := SendRequestForLink(personal[0].ProfileImgLink)
			if err != nil {
				log.Println(err)
			}
			interndata.Link = link
			list = append(list, interndata)
			//log.Println(i)
		} else {
			list = append(list, interndata)
		}
	}
	mapd["error"] = false
	mapd["message"] = "Operation successfull"
	mapd["list"] = list
	mapd["count"] = count.Count
	return mapd, 200
}

//DeleteUserAssignedReviewer is used to delete reviewer from users table
func DeleteUserAssignedReviewer(id int) error {
	defer util.Panic()
	reviewer := database.Reviewer{}
	db := config.DB
	rows := db.Where("reviewerid=?", id).Delete(&reviewer).RowsAffected
	if rows == 0 {
		return errors.New("Operation unsuccessfull, could not find user")
	}
	return nil
}

//DeleteAssignedReviewer is used to delete reviewer from users table
func DeleteAssignedReviewer(uid, rid int) error {
	defer util.Panic()
	reviewer := database.Reviewer{}
	db := config.DB
	rows := db.Where("reviewerid=? and userid=?", rid, uid).Delete(&reviewer).RowsAffected
	if rows == 0 {
		return errors.New("Operation unsuccessfull, could not find user")
	}
	return nil
}
