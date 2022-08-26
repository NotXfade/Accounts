package accounts

import (
	"log"
	"testing"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
)

func TestDelete(t *testing.T) {
	db := config.DB
	config.Conf.Service.Environment = "development"
	_, code := RegisterAccount("XSTestV1.deletesomeonetest@example.com", "12345789", "someone", "123456789")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	account := database.Accounts{}
	db.Where("email", "deletesomeonetest@example.com").Find(&account)
	personal := database.PersonalInfo{
		Userid:      account.ID,
		FathersName: "Dad",
		MothersName: "Mom",
		Gender:      "Female",
	}
	err := db.Create(&personal).Error
	if err != nil {
		log.Println(err)
	}
	userAddress := database.Address{
		Userid:      account.ID,
		AddressType: "permanent_address",
		Address:     "some house some city some country",
		City:        "CITY",
		District:    "District",
		State:       "State",
		Pincode:     819739133,
		Country:     "India",
	}
	err = db.Create(&userAddress).Error
	if err != nil {
		log.Println(err)
	}
	education := database.Education{
		Userid:      account.ID,
		Name:        "Btech",
		Institution: "ABC",
		Location:    "ABC",
		Year:        "2020",
		Percentage:  9,
	}
	err = db.Create(&education).Error
	if err != nil {
		log.Println(err)
	}
	emergencyContact := database.EmergencyContact{
		Userid:    account.ID,
		Name:      "someone",
		ContactNo: "9876543210",
	}
	err = db.Create(&emergencyContact).Error
	if err != nil {
		log.Println(err)
	}
	others := database.Others{
		Userid:           account.ID,
		Name:             "certificate1",
		CertificatesLink: "example@gmail.com/test.pdf",
	}
	err = db.Create(&others).Error
	if err != nil {
		log.Println(err)
	}
	//if account exists
	_, code = Delete("deletesomeonetest@example.com")
	if code != 200 {
		t.Error(code)
	}
	//when account does not exists
	_, code = Delete("deletesomeonetest@example.com")
	if code != 404 {
		t.Error(code)
	}
}
