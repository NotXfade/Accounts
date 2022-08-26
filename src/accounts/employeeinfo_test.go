package accounts

import (
	"testing"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
)

func TestStoreEmployeeInfo(t *testing.T) {
	empdata := database.EmployeeInfo{}
	//case : if uid is 0
	_, code := StoreEmployeeInfo(empdata, 0)
	if code != 400 {
		t.Errorf("expected %d , found %d", 400, code)
	}
	//store info with info
	empdata.PersonalInfo = database.PersonalInfo{
		FathersName: "Dad",
		MothersName: "Mom",
		Gender:      "Female",
	}
	userAddress := database.Address{
		AddressType: "permanent_address",
		Address:     "some house some city some country",
		City:        "CITY",
		District:    "District",
		State:       "State",
		Pincode:     819739133,
		Country:     "India",
	}
	empdata.UserAddress = append(empdata.UserAddress, userAddress)
	education := database.Education{
		Name:        "Btech",
		Institution: "ABC",
		Location:    "ABC",
		Year:        "2020",
		Percentage:  9,
	}
	empdata.Education = append(empdata.Education, education)
	emergencyContact := database.EmergencyContact{
		Name:      "someone",
		ContactNo: "9876543210",
	}
	empdata.EmergencyContact = append(empdata.EmergencyContact, emergencyContact)
	others := database.Others{
		Name:             "certificate1",
		CertificatesLink: "example@gmail.com/test.pdf",
	}
	empdata.Others = append(empdata.Others, others)
	_, code = StoreEmployeeInfo(empdata, 2)
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
}

func TestGetListInterns(t *testing.T) {
	_, code := GetListInterns("L1", "10", "1", "", "", 0)
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	_, code = GetListInterns("L1", "10", "1", "dateofjoining", "", 0)
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	_, code = GetListInterns("L1", "10", "1", "scoring", "", 0)
	if code != 404 {
		t.Errorf("expected %d , found %d", 400, code)
	}
}

func TestExportData(t *testing.T) {
	_, _, err := ExportData("")
	if err != nil {
		t.Error(err)
	}
	_, _, err = ExportData("Level1")
	if err != nil {
		t.Error(err)
	}
	_, _, err = ExportData("L1")
	if err != nil {
		t.Error(err)
	}
}

func TestGetEducationDetails(t *testing.T) {
	_, err := GetEducationDetails(2)
	if err != nil {
		t.Error(err)
	}
}

func TestGetOthersDetails(t *testing.T) {
	_, err := GetOthersDetails(2)
	if err != nil {
		t.Error(err)
	}
}

func TestGetPersonalInfo(t *testing.T) {
	_, err := GetPersonalInfo(2)
	if err != nil {
		t.Error(err)
	}
}

func TestGetEmergencyContact(t *testing.T) {
	_, err := GetEmergencyContact(2)
	if err != nil {
		t.Error(err)
	}
}

func TestGetAddressDetails(t *testing.T) {
	_, err := GetAddressDetails(2)
	if err != nil {
		t.Error(err)
	}
}

func TestGetPandCInfo(t *testing.T) {
	_, err := GetPandCInfo(2)
	if err == nil {
		t.Error(err)
	}
}

func TestGetUserInfo(t *testing.T) {
	_, code := GetUserInfo("someone@example.com")
	if code != 200 {
		t.Error(code)
	}
}

func TestSendRequestForLink(t *testing.T) {
	_, err := SendRequestForLink("file.txt")
	if err == nil {
		t.Error(err)
	}
}

func TestSendRequestForStatus(t *testing.T) {
	_ = SendRequestForStatus("someone@example.com")
}

func TestUpdateEmployeeInfo(t *testing.T) {
	empData := database.UpdateEmployeeInfo{
		PersonalInfo: database.UpdatePersonalInfo{
			FathersName: "Dad",
			MothersName: "Mom",
			Gender:      "Female",
		},
		UserAddress: []database.UpdateAddress{
			database.UpdateAddress{
				AddressType: "permanent_address",
				Address:     "some house some city some country",
				City:        "CITY",
				District:    "District",
				State:       "State",
				Pincode:     819739133,
				Country:     "India",
			},
		},
		EmergencyContact: []database.UpdateEmergencyContact{
			database.UpdateEmergencyContact{
				Name:      "someone",
				ContactNo: "9876543210",
			},
		},
		Education: []database.UpdateEducation{
			database.UpdateEducation{
				Name:        "Btech",
				Institution: "ABC",
				Location:    "ABC",
				Year:        "2020",
				Percentage:  9,
			},
		},
		Others: []database.UpdateOthers{
			database.UpdateOthers{
				Name:             "certificate1",
				CertificatesLink: "example@gmail.com/test.pdf",
			},
		},
	}
	_, code := UpdateEmployeeInfo(empData, 2)
	if code != 200 {
		t.Error(code)
	}
}

func TestSavePersonalInfo(t *testing.T) {
	personal := database.UpdatePersonalInfo{
		FathersName: "Dad",
		MothersName: "Mom",
		Gender:      "Female",
	}
	err := SavePersonalInfo(personal, 2)
	if err != nil {
		t.Error(err)
	}
}

func TestSaveUserAddress(t *testing.T) {
	address := database.UpdateAddress{
		AddressType: "permanent_address",
		Address:     "some house some city some country",
		City:        "CITY",
		District:    "District",
		State:       "State",
		Pincode:     819739133,
		Country:     "India",
	}
	userAddress := []database.UpdateAddress{}
	userAddress = append(userAddress, address)
	err := SaveUserAddress(userAddress, 2)
	if err != nil {
		t.Error(err)
	}
}

func TestSaveEmergencyContacts(t *testing.T) {
	emergencyContact := database.UpdateEmergencyContact{
		Name:      "someone",
		ContactNo: "9876543210",
	}
	contacts := []database.UpdateEmergencyContact{}
	contacts = append(contacts, emergencyContact)
	err := SaveEmergencyContacts(contacts, 2)
	if err != nil {
		t.Error(err)
	}
}

func TestSaveEducation(t *testing.T) {
	education := database.UpdateEducation{
		Name:        "Btech",
		Institution: "ABC",
		Location:    "ABC",
		Year:        "2020",
		Percentage:  9,
	}
	educations := []database.UpdateEducation{}
	educations = append(educations, education)
	err := SaveEducation(educations, 2)
	if err != nil {
		t.Error(err)
	}
}

func TestSaveOtherInformation(t *testing.T) {
	other := database.UpdateOthers{
		Name:             "certificate1",
		CertificatesLink: "example@gmail.com/test.pdf",
	}
	others := []database.UpdateOthers{}
	others = append(others, other)
	err := SaveOtherInformation(others, 2)
	if err != nil {
		t.Error(err)
	}
}

func TestProfileStatus(t *testing.T) {
	db := config.DB
	accounts := database.Accounts{}
	db.Where("id=?", 2).Find(&accounts)
	ProfileStatus(accounts.Email)
}

func TestUserInfo(t *testing.T) {
	_, code := RegisterAccount("XSTestV1.userinfo@example.com", "12345789", "someone", "123456789")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	_, code = UserInfo("userinfo@example.com")
	if code != 200 {
		t.Error("expected ", 200, " Found ", code)
	}
}

func BenchmarkUserInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UserInfo("userinfo@example.com")
	}
}
