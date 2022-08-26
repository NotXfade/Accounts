package accounts

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
)

//CreateDepartment is used for creating department
func CreateDepartment(name, shortname, email string) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	db := config.DB
	dept := database.Departments{
		Name:      name,
		ShortName: shortname,
		AddedBy:   email,
		Timestamp: time.Now().Unix(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(&dept).Error
	if err != nil {
		mapd["error"] = true
		mapd["message"] = "Department with name " + name + " already exists"
		return mapd, 400
	}
	mapd["error"] = false
	mapd["message"] = "Department created successfuly"
	return mapd, 200
}

//ListDepartment is used to list departments
func ListDepartment(limit, page string) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	db := config.DB
	dept := []database.Departments{}
	limits, _ := strconv.Atoi(limit)
	pageno, _ := strconv.Atoi(page)
	offset := (limits * (pageno - 1))
	var count int
	err := db.Find(&dept).Count(&count).Error
	if err != nil {
		log.Println(err)
	}
	err = db.Offset(offset).Limit(limits).Order("lower(name)").Find(&dept).Error
	if err != nil {
		log.Println(err)
	}
	mapd["error"] = false
	mapd["message"] = "Operation Successfull"
	mapd["list"] = dept
	mapd["count"] = count
	return mapd, 200
}

//UpdateDepartment is used for updating a department
func UpdateDepartment(name, shortname string, id int) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	db := config.DB
	dept := database.Departments{}
	db.Where("id=?", id).Find(&dept)
	if id == 0 {
		mapd["error"] = true
		mapd["message"] = "Department does not exists"
		return mapd, 400
	}
	dept.ID = id
	dept.Name = name
	dept.ShortName = shortname
	dept.UpdatedAt = time.Now()
	err := db.Save(&dept).Error
	if err != nil {
		mapd["error"] = true
		mapd["message"] = "Department with name " + name + " already exists"
		return mapd, 400
	}
	db.Model(&database.PandCData{}).Where("department=?", id).Updates(map[string]interface{}{"team": name})
	updateModule(id, name)
	mapd["error"] = false
	mapd["message"] = "Department updated successfuly"
	return mapd, 200
}

//DeleteDepartment is used for deleting department
func DeleteDepartment(id int) (map[string]interface{}, int) {
	defer util.Panic()
	mapd := make(map[string]interface{})
	db := config.DB
	dept := database.Departments{}
	err := db.Where("id=?", id).Delete(&dept).Error
	if err != nil {
		log.Println(err)
		mapd["error"] = true
		mapd["message"] = "Error occured while deleting department"
		return mapd, 200
	}
	pandc := database.PandCData{}
	err = db.Where("team=?", id).Delete(&pandc).Error
	if err != nil {
		log.Println(err)
	}
	db.Model(&database.PandCData{}).Where("department=?", id).Updates(map[string]interface{}{"team": "", "department": 0})
	updateModule(id, "")
	mapd["error"] = false
	mapd["message"] = "Department deleted successfuly"
	return mapd, 200
}

//request is used for sending payload data
type request struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//updateModule is used for updating modules
func updateModule(id int, name string) {
	req := request{
		ID:   id,
		Name: name,
	}
	payload, err := json.Marshal(&req)
	nc := config.NC
	res, err := nc.Request(config.Conf.NatsServer.Subject+".module.department", payload, 1*time.Minute)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(res.Data))
}

//GetDepartment is used to get department info
func GetDepartment(id int) database.Departments {
	db := config.DB
	dept := database.Departments{}
	db = db.Debug()
	err := db.Where("id=?", id).Find(&dept).Error
	if err != nil {
		log.Println(err)
	}
	return dept
}
