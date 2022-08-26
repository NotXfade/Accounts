package accounts

import (
	"log"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
)

//====================================== DELETE API ===========================================

//Delete is used to delete all data related to an email
func Delete(email string) (map[string]interface{}, int) {
	defer util.Panic()
	//log.Println("reached here")
	mapd := make(map[string]interface{})
	db := config.DB
	nc := config.NC
	//db = db.Debug()
	accounts := []database.Accounts{}
	db.Where("email=?", email).Find(&accounts)
	if len(accounts) == 0 {
		mapd["error"] = true
		mapd["message"] = "account does not exists"
		return mapd, 404
	}
	personalinfo := []database.PersonalInfo{}
	db.Where("userid=?", accounts[0].ID).Find(&personalinfo)
	account := database.Accounts{}
	err := db.Where("id=?", accounts[0].ID).Delete(&account).Error
	if err != nil {
		log.Println(err)
		mapd["error"] = true
		mapd["message"] = "could not delete" + err.Error()
		return mapd, 400
	}
	if len(personalinfo) != 0 {
		//========== delete profile image ==============
		data := []byte(personalinfo[0].ProfileImgLink)
		_, err = nc.Request(config.Conf.NatsServer.Subject+".documents.delete", data, 2*time.Minute)
		if err != nil {
			log.Println(err)
		}
		personal := database.PersonalInfo{}
		db.Where("id=?", personalinfo[0].ID).Delete(&personal)
	}
	userAddress := database.Address{}
	db.Where("userid=?", accounts[0].ID).Delete(&userAddress)
	emergencyContact := database.EmergencyContact{}
	db.Where("userid=?", accounts[0].ID).Delete(&emergencyContact)
	others := []database.Others{}
	db.Where("userid=?", accounts[0].ID).Find(&others)
	for i := 0; i < len(others); i++ {
		data := []byte(others[i].CertificatesLink)
		_, err = nc.Request(config.Conf.NatsServer.Subject+".documents.delete", data, 2*time.Minute)
		if err != nil {
			log.Println(err)
		}
		other := database.Others{}
		db.Where("id=?", others[0].ID).Delete(&other)
	}
	education := []database.Education{}
	db.Where("userid=?", accounts[0].ID).Find(&education)
	for i := 0; i < len(education); i++ {
		data := []byte(education[i].DegreeLink)
		_, err = nc.Request(config.Conf.NatsServer.Subject+".documents.delete", data, 2*time.Minute)
		if err != nil {
			log.Println(err)
		}
		links := education[i].MarksheetLink
		for i := 0; i < len(links); i++ {
			data := []byte(links[i])
			_, err = nc.Request(config.Conf.NatsServer.Subject+".documents.delete", data, 2*time.Minute)
			if err != nil {
				log.Println(err)
			}
		}
		eduDetail := database.Education{}
		db.Where("id=?", others[0].ID).Delete(&eduDetail)
	}
	//=============delete modules =================
	//log.Println("reached here")
	data := []byte(email)
	_, err = nc.Request(config.Conf.NatsServer.Subject+".module.delete", data, 2*time.Minute)
	if err != nil {
		log.Println(err)
	}
	//log.Println("reached here")
	mapd["error"] = false
	mapd["message"] = "account deleted successfully"
	return mapd, 200
}
