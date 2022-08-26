package accounts

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
)

// ================================ STORE EMPLOYEE INFO ====================================

//StoreEmployeeInfo is used to store employee data required
func StoreEmployeeInfo(empdata database.EmployeeInfo, id int) (map[string]interface{}, int) {
	defer util.Panic()
	db := config.DB
	mapd := make(map[string]interface{})
	//Add Personal Info in DB
	if id <= 0 {
		mapd["error"] = true
		mapd["message"] = "Invalid user"
		return mapd, 400
	}
	empdata.PersonalInfo.Userid = id
	err := db.Create(&empdata.PersonalInfo).Error
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < len(empdata.Education); i++ {
		//Add education details
		empdata.Education[i].Userid = id
		err = db.Create(&empdata.Education[i]).Error
		if err != nil {
			log.Println(err)
		}
	}
	for i := 0; i < len(empdata.Others); i++ {
		//store other certificates data
		empdata.Others[i].Userid = id
		err = db.Create(&empdata.Others[i]).Error
		if err != nil {
			log.Println(err)
		}
	}
	for i := 0; i < len(empdata.EmergencyContact); i++ {
		//store emergency contacts
		empdata.EmergencyContact[i].Userid = id
		err = db.Create(&empdata.EmergencyContact[i]).Error
		if err != nil {
			log.Println(err)
		}
	}
	for i := 0; i < len(empdata.UserAddress); i++ {
		//store other certificates data
		empdata.UserAddress[i].Userid = id
		if empdata.UserAddress[i].AddressType != "permanent_address" && empdata.UserAddress[i].AddressType != "current_address" {
			err = db.Create(&empdata.UserAddress[i]).Error
			if err != nil {
				log.Println(err)
			}
		}
	}
	mapd["error"] = "false"
	mapd["message"] = "Information saved successfully"
	mapd["link"] = empdata.PersonalInfo.ProfileImgLink
	return mapd, 200
}

//============================ SEND REQUEST TO VERIFY LINKS ===========================

//Checklink is used to check whether link is present in db or not
/* func Checklink(link string) bool {
	nc := config.NC
	res, err := nc.Request("xsonboarding.documents.checklink", []byte(link), 1*time.Minute)
	if err != nil {
		log.Println(err)
	}
	val := string(res.Data)
	if val == "true" {
		return true
	} else if val == "false" {
		return false
	}
	log.Print(val)
	return false
} */

//========================================== Get List of Interns ===================================

//InternsList is a type for sending interns information
type InternsList struct {
	Email      string  `json:"email"`
	Name       string  `json:"name"`
	Department string  `json:"department"`
	Task       string  `json:"task"`
	Status     string  `json:"status"`
	Progress   float64 `json:"progress"`
	Scoring    float64 `json:"scores"`
	Link       string  `json:"link"`
}

//listInterns is used to get details of List Interns
type listInterns struct {
	Email          string  `json:"email"`
	Name           string  `json:"name"`
	Level          string  `json:"level"`
	Team           string  `json:"team"`
	Slug           string  `json:"slug"`
	Score          float64 `json:"scores"`
	CreatedAt      string  `json:"created_at"`
	ProfileImgLink string  `json:"profile_img_link"`
}

//Count is used to get count
type Count struct {
	Count int
}

//GetListInterns is used to fetch list of employees according to their level and pagination implemented
func GetListInterns(level, limit, page, sorting, search string, progress float64) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	list := []InternsList{}
	accounts := []database.Accounts{}
	limits, _ := strconv.Atoi(limit)
	pageno, _ := strconv.Atoi(page)
	offset := (limits * (pageno - 1))
	db := config.DB

	query := "select * from accounts where level='" + level + "' AND account_status='active' "
	if sorting == "dateofjoining" { //sort according to date of joining
		query = query + "order by created_at DESC "
	}
	if search != "" {
		query = query + "AND (name ::text LIKE '%" + search + "%' or UPPER(name) ::text LIKE '%" + strings.ToUpper(search) + "%' or email ::text LIKE '%" + search + "%') "
	}
	query = query + "ORDER BY lower(name) " + "LIMIT " + limit + " OFFSET " + strconv.Itoa(offset) + ";"
	//db = db.Debug()
	db.Raw(query).Scan(&accounts)

	if len(accounts) == 0 {
		mapd["error"] = false
		mapd["message"] = "Operation Successful"
		mapd["list"] = list
		mapd["count"] = len(accounts)
		return mapd, 200
	}
	if progress >= 0 {
		for i := 0; i < len(accounts); i++ {
			var interndata InternsList
			internsList := SendRequestForStatus(accounts[i].Email)
			if internsList.Scoring == progress {
				interndata.Task = internsList.Task
				interndata.Status = internsList.Status
				interndata.Progress = internsList.Progress
				interndata.Scoring = internsList.Scoring
				interndata.Email = accounts[i].Email
				interndata.Name = accounts[i].Name
				pandc := database.PandCData{}
				db.Where("userid=?", accounts[i].ID).Find(&pandc)
				interndata.Department = pandc.Team
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
				//log.Println(i)
			}
		}
	} else {
		for i := 0; i < len(accounts); i++ {
			var interndata InternsList
			interndata.Email = accounts[i].Email
			interndata.Name = accounts[i].Name
			internsList := SendRequestForStatus(accounts[i].Email)
			//log.Println(internsList)
			interndata.Task = internsList.Task
			interndata.Status = internsList.Status
			interndata.Progress = internsList.Progress
			interndata.Scoring = internsList.Scoring
			//log.Println(interndata)
			pandc := database.PandCData{}
			db.Where("userid=?", accounts[i].ID).Find(&pandc)
			interndata.Department = pandc.Team
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
			//log.Println(i)
		}
	}
	//if sort by progress
	if sorting == "progress" {
		sort.Slice(list, func(i, j int) bool {
			return list[i].Scoring > list[j].Scoring
		})
	}
	mapd["error"] = false
	mapd["message"] = "Operation Successful"
	mapd["list"] = list
	mapd["count"] = len(accounts)
	return mapd, 200
}

// //GetListInterns is used to get list of Interns
// func GetListInterns(level, limit, page, search, module, batch, department, score, reviewer string) (map[string]interface{}, int) {
// 	defer util.Panic()
// 	mapd := make(map[string]interface{})
// 	list := []listInterns{}
// 	limits, _ := strconv.Atoi(limit)
// 	pageno, _ := strconv.Atoi(page)
// 	offset := (limits * (pageno - 1))
// 	db := config.DB
// 	db = db.Debug()
// 	count := Count{}
// 	query := `level='` + level + `' AND account_status='active' `
// 	scoreQuery := ""
// 	log.Println(reviewer)
// 	if reviewer != "0" {
// 		query = query + `AND reviewerid=` + reviewer
// 	}
// 	if module != "" {
// 		query = query + ` AND slug='` + module + `' `
// 	}
// 	if department != "" {
// 		query = query + `AND departmentid=` + department + `' `
// 	}
// 	if batch != "" {
// 		query = query + " AND  created_at ::date ='" + batch + "' "
// 	}
// 	if search != "" {
// 		query = query + " AND ( name ::text LIKE '%" + search + "%' or email ::text LIKE '%" + search + "%') "
// 	}
// 	if score != "" {
// 		scoreQuery = "where score >= 75"
// 	}
// 	endQuery := " group by email,name,level,team,profile_img_link ,created_at, form_filled_at ,scores,total_scores"
// 	finalQuery := "select distinct email,name,level,team,scores,total_scores,profile_img_link,created_at,form_filled_at from internslist where " + query + endQuery
// 	listQuery := "from (select email,name,level,team,COALESCE(round(sum(scores) / NULLIF(sum(total_scores),0)*100,2), 0) as score,profile_img_link,created_at from ( " + finalQuery + " ) as list "
// 	listQuery = listQuery + " " + " group by email,name,level,team,profile_img_link,created_at ) as tb " + scoreQuery

// 	err := db.Raw("select count(email) " + listQuery).Scan(&count).Error
// 	if err != nil {
// 		mapd["error"] = true
// 		mapd["message"] = "Internal Server Error"
// 		return mapd, 500
// 	}
// 	listQuery = listQuery + " ORDER BY lower(name) " + "LIMIT " + limit + " OFFSET " + strconv.Itoa(offset) + ";"
// 	err = db.Raw("select * " + listQuery).Scan(&list).Error
// 	if err != nil {
// 		mapd["error"] = true
// 		mapd["message"] = "Internal Server Error"
// 		return mapd, 500
// 	}
// 	internsList := []InternsList{}
// 	for i := 0; i < len(list); i++ {
// 		var interndata InternsList
// 		interndata.Email = list[i].Email
// 		interndata.Name = list[i].Name
// 		internsStatus := SendRequestForStatus(list[i].Email)
// 		interndata.Task = internsStatus.Task
// 		interndata.Status = internsStatus.Status
// 		interndata.Progress = internsStatus.Progress
// 		interndata.Scoring = internsStatus.Scoring
// 		interndata.Department = list[i].Team
// 		if list[i].ProfileImgLink != "" {
// 			link, err := SendRequestForLink(list[i].ProfileImgLink)
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			interndata.Link = link
// 		}
// 		internsList = append(internsList, interndata)
// 	}
// 	mapd["error"] = false
// 	mapd["message"] = "Operation successful"
// 	mapd["list"] = internsList
// 	mapd["count"] = count.Count
// 	return mapd, 200
// }

//======================================== EXPORT DATA TO CSV ===================================================

//ExportData is used to return data in csv format
func ExportData(level string) ([]byte, int, error) {
	defer util.Panic()
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	db := config.DB
	accounts := []database.Accounts{}
	//if level is empty fetch all the details of interns
	if level == "" {
		err := db.Where("role=? AND account_status=?", "user", "active").Find(&accounts).Error
		if err != nil {
			log.Println(err)
			return b.Bytes(), 400, err
		}
	} else { //fetch details according to level
		err := db.Where("level=? AND account_status=?", level, "active").Find(&accounts).Error
		if err != nil {
			log.Println(err)
			return b.Bytes(), 400, err
		}
	}
	record := []string{"Name", "Email", "ContactNo", "Level", "Father's Name", "Mother's Name", "Date Of Birth", "Gender", "Blood Group", "Marital Status", "Spouse Name", "Spouse Contact"}
	//write record in buffer
	if err := w.Write(record); err != nil {
		log.Println("error writing record to csv:", err)
	}
	//log.Println(accounts)
	for i := 0; i < len(accounts); i++ {
		record := make([]string, 0)
		//----------- added account information -------------------
		record = append(record, accounts[i].Name, accounts[i].Email, accounts[i].ContactNo, accounts[i].Level)
		//--------- Get Personal Information --------------
		personal, err := GetPersonalInfo(accounts[i].ID)
		if err != nil {
			log.Println(err)
		}
		record = append(record, personal.FathersName, personal.MothersName, personal.DateOfBirth, personal.Gender, personal.BloodGroup, personal.MaritalStatus, personal.SpouseName, personal.SpouseContact)
		//write record in buffer
		if err := w.Write(record); err != nil {
			log.Println("error writing record to csv:", err)
		}
	}
	//call flush to
	w.Flush()
	//log.Println(string(b.Bytes()))
	return b.Bytes(), 200, nil
}

//=================================== FETCH EDUCATION DETAILS FROM DATABASE ==========================================

//GetEducationDetails is used to return details based on uid
func GetEducationDetails(uid int) ([]database.Education, error) {
	defer util.Panic()
	db := config.DB
	education := []database.Education{}
	err := db.Where("userid=?", uid).Find(&education).Error
	if err != nil {
		return education, err
	}
	//log.Println(education)
	return education, nil
}

//==================================  FETCH OTHER CERTIICATES DETAILS FROM DATABASE ====================================

//GetOthersDetails is used to return details based on uid
func GetOthersDetails(uid int) ([]database.Others, error) {
	defer util.Panic()
	db := config.DB
	others := []database.Others{}
	err := db.Where("userid=?", uid).Find(&others).Error
	if err != nil {
		return others, err
	}
	//log.Println(uid, others)
	return others, nil
}

//================================== FETCH PERSONAL DETAILS FROM DATABASE ===================================

//GetPersonalInfo is used to return details based on uid
func GetPersonalInfo(uid int) (database.PersonalInfo, error) {
	defer util.Panic()
	db := config.DB
	personal := database.PersonalInfo{}
	err := db.Where("userid=?", uid).Find(&personal).Error
	if err != nil {
		return personal, err
	}
	t := strings.Split(personal.DateOfBirth, "T")
	personal.DateOfBirth = t[0]
	//log.Println(personal)
	return personal, nil
}

//================================ GET EMERGENCY CONTACT DETAILS ===========================================

//GetEmergencyContact is used to return details based on uid
func GetEmergencyContact(uid int) ([]database.EmergencyContact, error) {
	defer util.Panic()
	db := config.DB
	emergency := []database.EmergencyContact{}
	err := db.Where("userid=?", uid).Find(&emergency).Error
	if err != nil {
		return emergency, err
	}
	//	log.Println(emergency)
	return emergency, nil
}

//================================ GET ADDRESS DETAILS ===========================================

//GetAddressDetails is used to return details based on uid
func GetAddressDetails(uid int) ([]database.Address, error) {
	defer util.Panic()
	db := config.DB
	address := []database.Address{}
	err := db.Where("userid=?", uid).Find(&address).Error
	if err != nil {
		return address, err
	}
	//	log.Println(address)
	return address, nil
}

//================================ GET P&C DETAILS ===========================================

//GetPandCInfo is used to return details based on uid
func GetPandCInfo(uid int) (database.PandCData, error) {
	defer util.Panic()
	db := config.DB
	pandc := database.PandCData{}
	err := db.Where("userid=?", uid).Find(&pandc).Error
	if err != nil {
		return pandc, err
	}
	return pandc, nil
}

//GetReviewer is used to get details of user assigned reviewer
func GetReviewer(uid int) ([]database.Reviewer, error) {
	defer util.Panic()
	db := config.DB
	reviewer := []database.Reviewer{}
	err := db.Where("userid=?", uid).Find(&reviewer).Error
	if err != nil {
		return reviewer, err
	}
	return reviewer, nil
}

//======================================= GETUSERINFO =============================================

//GetUserInfo is used to fetch user details
func GetUserInfo(email string) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	account, err := GetAccountForEmail(email)
	if err != nil {
		mapd["error"] = true
		mapd["message"] = "User Not Found"
		return mapd, 200
	}
	account.Role = strings.Title(strings.ToLower(account.Role))
	personalInfo, err := GetPersonalInfo(account.ID)
	link, err := SendRequestForLink(personalInfo.ProfileImgLink)
	if err != nil {
		log.Println(err)
	}
	//log.Println(link)
	personalInfo.ProfileImgLink = link
	educationDetail, err := GetEducationDetails(account.ID)
	if err != nil {
		log.Println(err)
	}
	otherCertificates, err := GetOthersDetails(account.ID)
	if err != nil {
		log.Println(err)
	}
	emergencyContacts, err := GetEmergencyContact(account.ID)
	if err != nil {
		log.Println(err)
	}
	pandc, err := GetPandCInfo(account.ID)
	if err != nil {
		log.Println(err)
	}
	addressInfo, err := GetAddressDetails(account.ID)
	if err != nil {
		log.Println(err)
	}
	//pandc.Reviewer = []string{""}
	for i := 0; i < (2 - len(emergencyContacts)); i++ {
		emergencyContacts = append(emergencyContacts, database.EmergencyContact{
			Name:      "",
			ContactNo: "",
		})
	}
	if len(addressInfo) != 2 {
		//log.Println("reached here")
		for i := 0; i < 2; i++ {
			flag := ""
			for j := 0; j < len(addressInfo); j++ {
				if i == 0 {
					if addressInfo[j].AddressType == "Current address" {
						flag = "Current address"
					}
				} else if i == 1 {
					if addressInfo[j].AddressType == "Permanent address" {
						flag = "Permanent address"
					}
				}
			}
			if flag == "" {
				if i == 0 {
					flag = "Current address"
				} else if i == 1 {
					flag = "Permanent address"
				}
				addressInfo = append(addressInfo, database.Address{
					AddressType: flag,
				})
			}
		}
	}
	mapd["error"] = false
	mapd["message"] = "Operation Successful"
	mapd["accounts"] = account
	mapd["personal_info"] = personalInfo
	mapd["address_info"] = addressInfo
	mapd["education_details"] = educationDetail
	mapd["others_certificates"] = otherCertificates
	mapd["emergency_contacts"] = emergencyContacts
	mapd["pandcinfo"] = pandc
	return mapd, 200
}

//==================================== GET USER INFO =============================================

//GetUserInfo is used to fetch user details
func UserInfo(email string) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	account, err := GetAccountForEmail(email)
	if err != nil {
		mapd["error"] = true
		mapd["message"] = "User Not Found"
		return mapd, 200
	}
	account.Role = strings.Title(strings.ToLower(account.Role))
	personalInfo, err := GetPersonalInfo(account.ID)
	link, err := SendRequestForLink(personalInfo.ProfileImgLink)
	if err != nil {
		log.Println(err)
	}
	//log.Println(link)
	//personalInfo.ProfileImgLink = link
	educationDetail, err := GetEducationDetails(account.ID)
	if err != nil {
		log.Println(err)
	}
	otherCertificates, err := GetOthersDetails(account.ID)
	if err != nil {
		log.Println(err)
	}
	emergencyContacts, err := GetEmergencyContact(account.ID)
	if err != nil {
		log.Println(err)
	}
	if len(emergencyContacts) < 2 {
		for i := 0; i < (2 - len(emergencyContacts)); i++ {
			emergencyContacts = append(emergencyContacts, database.EmergencyContact{
				Name:      "",
				ContactNo: "",
			})
		}
	}
	pandc, err := GetPandCInfo(account.ID)
	if err != nil {
		log.Println(err)
	}
	addressInfo, err := GetAddressDetails(account.ID)
	if err != nil {
		log.Println(err)
	}
	if len(addressInfo) != 2 {
		//log.Println("reached here")
		for i := 0; i < 2; i++ {
			flag := ""
			for j := 0; j < len(addressInfo); j++ {
				if i == 0 {
					if addressInfo[j].AddressType == "Current address" {
						flag = "Current address"
					}
				} else if i == 1 {
					if addressInfo[j].AddressType == "Permanent address" {
						flag = "Permanent address"
					}
				}
			}
			if flag == "" {
				if i == 0 {
					flag = "Current address"
				} else if i == 1 {
					flag = "Permanent address"
				}
				addressInfo = append(addressInfo, database.Address{
					AddressType: flag,
				})
			}
		}
	}
	mapdd := make(map[string]interface{})
	mapd["error"] = false
	mapdd["pandcinfo"] = pandc
	mapdd["education"] = educationDetail
	mapdd["address"] = addressInfo
	mapdd["others"] = otherCertificates
	mapdd["emergency"] = emergencyContacts
	mapdd["personal"] = personalInfo
	mapdd["account"] = account
	mapd["percentage"] = ProfileStatus(account.Email)
	mapd["message"] = "Operation Successful"
	mapd["profile"] = mapdd
	mapd["link"] = link
	return mapd, 200
}

//SendRequestForLink is used to fetch presigned url link
func SendRequestForLink(filename string) (string, error) {
	defer util.Panic()
	payload := []byte(filename)
	nc := config.NC
	res, err := nc.Request(config.Conf.NatsServer.Subject+".documents.getlink", payload, 1*time.Minute)
	if err != nil {
		log.Println(err)
		return "", err
	}
	type Response struct {
		Error   string `json:"error"`
		Message string `json:"message"`
		Link    string `json:"link"`
	}
	response := Response{}
	json.Unmarshal(res.Data, &response)
	return response.Link, errors.New(response.Error)
}

//SendRequestForStatus is used to get progress of user
func SendRequestForStatus(email string) InternsList {
	defer util.Panic()
	payload := []byte(email)
	nc := config.NC
	internsList := InternsList{}
	res, err := nc.Request(config.Conf.NatsServer.Subject+".module.progress", payload, 1*time.Minute)
	if err != nil {
		log.Println(err)
		return internsList
	}
	//log.Println(string(res.Data))
	err = json.Unmarshal(res.Data, &internsList)
	if err != nil {
		log.Println(err)
		return internsList
	}
	return internsList
}

//=============================== Update Employee Info ===============================

//UpdateEmployeeInfo is used to update employee info
func UpdateEmployeeInfo(empData database.UpdateEmployeeInfo, id int) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	SaveAccountInfo(empData.Account, id)
	SavePersonalInfo(empData.PersonalInfo, id)
	SaveUserAddress(empData.UserAddress, id)
	SaveEmergencyContacts(empData.EmergencyContact, id)
	SaveEducation(empData.Education, id)
	SaveOtherInformation(empData.Others, id)
	mapd["error"] = false
	mapd["message"] = "Information updated successfuly"
	return mapd, 200
}

//SaveAccountInfo is used to update name
func SaveAccountInfo(accountsData database.Accounts, id int) error {
	defer util.Panic()
	db := config.DB
	//db = db.Debug()
	accounts := database.Accounts{}
	err := db.Where("id=?", id).Find(&accounts).Error
	if err != nil {
		log.Println(err)
		return err
	}
	//log.Println(accountsData.ContactNo)
	accounts.Name = accountsData.Name
	accounts.ContactNo = accountsData.ContactNo
	err = db.Save(&accounts).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//SavePersonalInfo is used to create or update personal info
func SavePersonalInfo(personalData database.UpdatePersonalInfo, id int) error {
	defer util.Panic()
	db := config.DB
	//db = db.Debug()
	personalInfo := []database.PersonalInfo{}
	personalData.Userid = id
	db.Where("userid=?", id).Find(&personalInfo)
	if len(personalInfo) == 0 {
		data := database.PersonalInfo{}
		data.Userid = id
		data.FathersName = personalData.FathersName
		data.MothersName = personalData.MothersName
		data.Gender = personalData.Gender
		data.DateOfBirth = personalData.DateOfBirth
		data.BloodGroup = personalData.BloodGroup
		data.MaritalStatus = personalData.MaritalStatus
		data.SpouseName = personalData.SpouseName
		data.SpouseContact = personalData.SpouseContact
		data.ProfileImgLink = personalData.ProfileImgLink
		data.UpdatedAt = time.Now()
		data.CreateAt = time.Now()
		err := db.Create(&data).Error
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		personalInfo[0].FathersName = personalData.FathersName
		personalInfo[0].MothersName = personalData.MothersName
		personalInfo[0].Gender = personalData.Gender
		personalInfo[0].DateOfBirth = personalData.DateOfBirth
		personalInfo[0].BloodGroup = personalData.BloodGroup
		personalInfo[0].MaritalStatus = personalData.MaritalStatus
		personalInfo[0].SpouseName = personalData.SpouseName
		personalInfo[0].SpouseContact = personalData.SpouseContact
		personalInfo[0].ProfileImgLink = personalData.ProfileImgLink
		personalInfo[0].UpdatedAt = time.Now()
		err := db.Save(&personalInfo[0]).Error
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

//SaveUserAddress is used to save user's address
func SaveUserAddress(address []database.UpdateAddress, id int) error {
	defer util.Panic()
	db := config.DB
	userAddress := database.Address{}
	db.Where("userid=?", id).Delete(&userAddress)
	for i := 0; i < len(address); i++ {
		addressUpdate := database.Address{
			Userid:      id,
			AddressType: address[i].AddressType,
			Address:     address[i].Address,
			City:        address[i].City,
			District:    address[i].District,
			State:       address[i].State,
			Pincode:     address[i].Pincode,
			Country:     address[i].Country,
			UpdatedAt:   time.Now(),
			CreateAt:    time.Now(),
		}
		err := db.Create(&addressUpdate).Error
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

//SaveEmergencyContacts is used to create or save emergency contacts
func SaveEmergencyContacts(contacts []database.UpdateEmergencyContact, id int) error {
	defer util.Panic()
	db := config.DB
	emergencyContact := database.EmergencyContact{}
	db.Where("userid=?", id).Delete(&emergencyContact)
	for i := 0; i < len(contacts); i++ {
		contact := database.EmergencyContact{
			Userid:    id,
			Name:      contacts[i].Name,
			ContactNo: contacts[i].ContactNo,
			UpdatedAt: time.Now(),
			CreateAt:  time.Now(),
		}
		err := db.Create(&contact).Error
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

//SaveEducation is used to save or create entry in database
func SaveEducation(education []database.UpdateEducation, id int) error {
	defer util.Panic()
	db := config.DB
	userEducation := database.Education{}
	db.Where("userid=?", id).Delete(&userEducation)
	for i := 0; i < len(education); i++ {
		educations := database.Education{
			Userid:        id,
			Name:          education[i].Name,
			Institution:   education[i].Institution,
			Location:      education[i].Location,
			Year:          education[i].Year,
			Percentage:    education[i].Percentage,
			MarksheetLink: education[i].MarksheetLink,
			DegreeLink:    education[i].DegreeLink,
			UpdatedAt:     time.Now(),
			CreateAt:      time.Now(),
		}
		err := db.Create(&educations).Error
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

//SaveOtherInformation is used to save info for other details
func SaveOtherInformation(others []database.UpdateOthers, id int) error {
	defer util.Panic()
	db := config.DB
	otherInfo := database.Others{}
	db.Where("userid=?", id).Delete(&otherInfo)
	for i := 0; i < len(others); i++ {
		other := database.Others{
			Userid:           id,
			Name:             others[i].Name,
			CertificatesLink: others[i].CertificatesLink,
			UpdatedAt:        time.Now(),
			CreateAt:         time.Now(),
		}
		err := db.Create(&other).Error
		if err != nil {
			log.Println(err)
			//return err
		}
	}
	return nil
}

//ProfileStatus is used to return status of profile completion
func ProfileStatus(email string) map[string]interface{} {
	defer util.Panic()
	mapd := make(map[string]interface{})
	var total, complete float64
	account, err := GetAccountForEmail(email)
	if err != nil {
		mapd["error"] = true
		mapd["message"] = "User Not Found"
		return mapd
	}
	var personalTotal, personalComplete float64
	//CheckFields for personal info
	personalInfo, err := GetPersonalInfo(account.ID)
	if err != nil {
		log.Println(err)
	}
	if personalInfo.ID == 0 {
		personalTotal = personalTotal + 7
	} else {
		if personalInfo.FathersName != "" {
			personalComplete = personalComplete + 1
		}
		if personalInfo.MothersName != "" {
			personalComplete = personalComplete + 1
		}
		if personalInfo.Gender != "" {
			personalComplete = personalComplete + 1
		}
		if personalInfo.DateOfBirth != "" {
			personalComplete = personalComplete + 1
		}
		if personalInfo.BloodGroup != "" {
			personalComplete = personalComplete + 1
		}
		if personalInfo.MaritalStatus != "" {
			personalComplete = personalComplete + 1
		}
		if personalInfo.ProfileImgLink != "" {
			personalComplete = personalComplete + 1
		}
		personalTotal = personalTotal + 7
	}

	//Get Address Details
	addressInfo, err := GetAddressDetails(account.ID)
	if err != nil {
		log.Println(err)
	}
	if len(addressInfo) == 0 {
		personalTotal = personalTotal + 14
	} else {
		for i := 0; i < len(addressInfo); i++ {
			if addressInfo[i].AddressType != "" {
				personalComplete = personalComplete + 1
			}
			if addressInfo[i].Address != "" {
				personalComplete = personalComplete + 1
			}
			if addressInfo[i].City != "" {
				personalComplete = personalComplete + 1
			}
			if addressInfo[i].District != "" {
				personalComplete = personalComplete + 1
			}
			if addressInfo[i].State != "" {
				personalComplete = personalComplete + 1
			}
			if addressInfo[i].Pincode != 0 {
				personalComplete = personalComplete + 1
			}
			if addressInfo[i].Country != "" {
				personalComplete = personalComplete + 1
			}
		}
		personalTotal = personalTotal + float64(14)
	}
	//log.Println(complete, total)
	//Get Education Details
	var educationalTotal, educationalComplete float64
	educationDetail, err := GetEducationDetails(account.ID)
	if err != nil {
		log.Println(err)
	}
	if len(educationDetail) == 0 {
		educationalTotal = educationalTotal + 6
	} else {
		for i := 0; i < len(educationDetail); i++ {
			if educationDetail[i].Name != "" {
				educationalComplete = educationalComplete + 1
			}
			if educationDetail[i].Institution != "" {
				educationalComplete = educationalComplete + 1
			}
			if educationDetail[i].Location != "" {
				educationalComplete = educationalComplete + 1
			}
			if educationDetail[i].Year != "" {
				educationalComplete = educationalComplete + 1
			}
			if educationDetail[i].Percentage != 0 {
				educationalComplete = educationalComplete + 1
			}
			if len(educationDetail[i].MarksheetLink) != 0 {
				educationalComplete = educationalComplete + 1
			}
		}
		educationalTotal = educationalTotal + float64(6*len(educationDetail))
	}
	//log.Println(complete, total)
	var otherTotal, otherComplete float64
	//Get Other Details
	otherCertificates, err := GetOthersDetails(account.ID)
	if err != nil {
		log.Println(err)
	}
	if len(otherCertificates) == 0 {
		otherTotal = otherTotal + 2
	} else {
		for i := 0; i < len(otherCertificates); i++ {
			if otherCertificates[i].Name != "" && otherCertificates[i].CertificatesLink != "" {
				otherComplete = otherComplete + 2
			}
		}
		otherTotal = otherTotal + float64(2*len(otherCertificates))
	}
	//log.Println(complete, total)
	//Get Emergency Contacts
	emergencyContacts, err := GetEmergencyContact(account.ID)
	if err != nil {
		log.Println(err)
	}
	if len(emergencyContacts) == 0 {
		otherTotal = otherTotal + 4
	} else {
		for i := 0; i < len(emergencyContacts); i++ {
			if emergencyContacts[i].Name != "" && emergencyContacts[i].ContactNo != "" {
				otherComplete = otherComplete + 2
			}
		}
		otherTotal = otherTotal + 4
	}
	complete = personalComplete + educationalComplete + otherComplete
	total = personalTotal + educationalTotal + otherTotal
	percentage := (complete / total) * 100
	mapd["total"] = math.Round(percentage*100) / 100
	mapd["personal_info"] = math.Round(((personalComplete/personalTotal)*100)*100) / 100
	mapd["others"] = math.Round(((otherComplete/otherTotal)*100)*100) / 100
	mapd["educational_info"] = math.Round(((educationalComplete/educationalTotal)*100)*100) / 100
	return mapd
}
