//+build !test

package database

import (
	"time"

	"github.com/lib/pq"
)

// Accounts is a strucutre to stores user information
type Accounts struct {
	// auto generated
	ID int `json:"id" gorm:"primary_key; unique_index:idx_account"`
	//If password is empty then the user registered from career forms.
	//User can register password in future using forgot password.
	Password string `json:"-" `
	//email of user
	Email string `json:"email" gorm:"not null;unique;unique_index:idx_account;"`
	//name of user
	Name string `json:"name"`
	// contact number of user
	ContactNo string `json:"contact_no"`
	//this role id is used for auth portal management
	//value can be 'admin', 'user'
	Role string `json:"role" gorm:"unique_index:idx_account;"`
	//Level wil be the level of employees like L1 intern, L2 intern
	Level string `json:"level" gorm:"unique_index:idx_account;"`
	//account status can be active, new, deleted, blocked etc.
	AccountStatus string `json:"account_status" gorm:"unique_index:idx_account;"`
	//AcceptPolicy
	AcceptedPolicy bool `json:"accepted_policy"`
	//account creation date
	Timestamp int64     `json:"timestamp"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

// Activities is a structure to record user activties
type Activities struct {
	ID int `json:"-" gorm:"primary_key"`
	//If username entered incorrect then these activities will also be recorded.
	Email string `json:"email" gorm:"index;not null;index:idx_activity;"`
	//can be login, failedlogin, signup
	ActivityName string    `json:"activity_name" gorm:"index:idx_activity;"`
	ClientIP     string    `json:"client_ip"`
	ClientAgent  string    `json:"client_agent"`
	Timestamp    int64     `json:"timestamp" gorm:"index:idx_activity;"`
	UpdatedAt    time.Time `json:"-"`
	CreatedAt    time.Time `json:"-"`
}

// Tokens is a structure to stores token for verifcation, invite link, forgot password
type Tokens struct {
	ID        int    `json:"-" gorm:"primary_key"`
	Email     string `gorm:"index:idx_ut"`
	Token     string `gorm:"index:idx_ut"` //if you require to store token secret then append that to it with '&' symbol.
	TokenTask string
	//Level wil be the level of employees like L1 intern, L2 intern
	Level     string `json:"level"`
	Status    string
	Timestamp int64
	UpdatedAt time.Time `json:"-"`
	CreateAt  time.Time `json:"-"`
}

// ActiveSessions is a structure to stores active sessions
type ActiveSessions struct {
	ID          int `json:"-" gorm:"primary_key"`
	SessionID   string
	Userid      int `gorm:"index"`
	ClientAgent string
	Start       int64
	// if value is '0' then session is remembered.
	End       int64
	UpdatedAt time.Time `json:"-"`
	CreateAt  time.Time `json:"-"`
}

//Education is a structure to store
type Education struct {
	ID int `json:"-" gorm:"primary_key;index:edu_idx"`
	//AccountID
	Userid int `json:"-" gorm:"index:edu_idx;not null"`
	//Name of the Degree/Class
	Name string `json:"name"  binding:"required"`
	//Name of School/University
	Institution string `json:"institution" binding:"required"`
	//Location of School/Univeristy
	Location string `json:"location"  binding:"required"`
	//Session
	Year string `json:"year"  binding:"required"`
	//percentage scored in Class/Degree
	Percentage float64 `json:"percentage" binding:"required"`
	//uploaded marksheet link
	MarksheetLink pq.StringArray `json:"marksheet_link" gorm:"type:text[];" binding:"required"`
	//Uploaded Degree Link in case of University
	DegreeLink string    `json:"degree_link"`
	UpdatedAt  time.Time `json:"-"`
	CreateAt   time.Time `json:"-"`
}

//Others : it is used to save extra certificates info of employee
type Others struct {
	ID     int `json:"-" gorm:"primary_key;index:oth_idx"`
	Userid int `json:"-"  gorm:"index:oth_idx"`
	//Name of the Certificate
	Name string `json:"name" binding:"required"`
	//Uploaded Certificate Link
	CertificatesLink string    `json:"certificate_link" binding:"required"`
	UpdatedAt        time.Time `json:"-"`
	CreateAt         time.Time `json:"-"`
}

//Departments is used to store department info
type Departments struct {
	ID        int       `json:"id" gorm:"primary_key`
	Name      string    `json:"name"  gorm:"unique"`
	ShortName string    `json:"short_name"`
	AddedBy   string    `json:"added_by"`
	Timestamp int64     `json:"timestamp"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

//EmployersInfo : it will store past employer info of a user
/* type EmployersInfo struct {
	ID int `json:"-" gorm:"primary_key"`
	//Account ID
	Userid int `json:"-"  gorm:"index"`
	//Name of the Past Company
	Company string `json:"company"  binding:"required"`
	//Past Role in the Company
	Role string `json:"role" binding:"required"`
	//Link of uploaded salary slips
	SalaryLink string    `json:"salary_link" binding:"required"`
	UpdatedAt  time.Time `json:"-"`
	CreateAt   time.Time `json:"-"`
} */

//BankDetails : structure of bank account details
/* type BankDetails struct {
	ID int `json:"-" gorm:"primary_key"`
	//Account ID
	Userid int `json:"-" gorm:"index"`
	//Bank Name
	BankName string `json:"bank_name" binding:"required"`
	//Name of the Account Owner
	Name string `json:"name" binding:"required"`
	//Bank Account Number
	AccountNumber string `json:"account_number" gorm:"not null;unique;" binding:"required"`
	//IFSC Code for Account Number
	IFSC      string    `json:"ifsc" binding:"required"`
	UpdatedAt time.Time `json:"-"`
	CreateAt  time.Time `json:"-"`
} */

/* //IdentityProof : structure of id proofs of employee
type IdentityProof struct {
	ID     int `json:"-" gorm:"primary_key"`
	Userid int `json:"-" gorm:"index"`
	//Pan Card Number
	PanNumber string `json:"pan_number"`
	//Uploaded Pan Card Link
	PanLink string `json:"pan_link"`
	//Aadhar Id
	AadharNumber string `json:"aadhar_number"  gorm:"not null;unique;" binding:"required"`
	//Uploaded Aadhar Card Link
	AadharLink string `json:"aadhar_link" binding:"required"`
	//PF Number
	PFNumber  string    `json:"pf_number" `
	UpdatedAt time.Time `json:"-"`
	CreateAt  time.Time `json:"-"`
}
*/
//PersonalInfo : structure of personal information of employee
type PersonalInfo struct {
	ID             int       `json:"-" gorm:"primary_key;index:personal_idx"`
	Userid         int       `json:"-"  gorm:"not null;unique;index:personal_idx"`
	FathersName    string    `json:"fathers_name" binding:"required"`
	MothersName    string    `json:"mothers_name" binding:"required"`
	Gender         string    `json:"gender" binding:"required"`
	DateOfBirth    string    `json:"dateofbirth" binding:"required"`
	BloodGroup     string    `json:"blood_group" binding:"required"`
	MaritalStatus  string    `json:"marital_Status" binding:"required"`
	SpouseName     string    `json:"spouse_name"`
	SpouseContact  string    `json:"spouse_contact"`
	ProfileImgLink string    `json:"profile_img_link" binding:"required"`
	UpdatedAt      time.Time `json:"-"`
	CreateAt       time.Time `json:"-"`
}

//Address is used to add address of user
type Address struct {
	ID          int       `json:"-" gorm:"primary_key;index:address_idx"`
	Userid      int       `json:"-"  gorm:"not null;index:address_idx"`
	AddressType string    `json:"address_type" binding:"required"` //value can be current Address or permanent address
	Address     string    `json:"address" binding:"required"`
	City        string    `json:"city" binding:"required"`
	District    string    `json:"district" binding:"required"`
	State       string    `json:"state" binding:"required"`
	Pincode     int       `json:"pincode" binding:"required"`
	Country     string    `json:"country" binding:"required"`
	UpdatedAt   time.Time `json:"-"`
	CreateAt    time.Time `json:"-"`
}

//EmergencyContact : It is a structure to store emergency contact number of employees
type EmergencyContact struct {
	ID        int       `json:"-" gorm:"primary_key;index:emergency_idx"`
	Userid    int       `json:"-" gorm:"index:emergency_idx"`
	Name      string    `json:"contact_name"  binding:"required"`
	ContactNo string    `json:"contact_number" binding:"required"`
	UpdatedAt time.Time `json:"-"`
	CreateAt  time.Time `json:"-"`
}

//PandCData is for company information which can only be filled by P&C team
type PandCData struct {
	ID             int       `json:"-" gorm:"primary_key;index:pandc_idx"`
	Userid         int       `json:"-" binding:"required" gorm:"index:pandc_idx"`
	Buddy          string    `json:"buddy" `
	Team           string    `json:"team" `
	Department     int       `json:"department"`
	Coach          string    `json:"coach" `
	HeadCoach      string    `json:"headcoach" `
	ProjectManager string    `json:"project_manager"`
	NexaOps        string    `json:"nexaops"`
	PeopleManager  string    `json:"people_manager"`
	DateOfJoining  time.Time `json:"date_of_joining" binding:"required"`
	UpdatedAt      time.Time `json:"-"`
	CreateAt       time.Time `json:"-"`
}

//Reviewer is used to assign reviewer details
type Reviewer struct {
	ID         int       `json:"-" gorm:"primary_key"`
	Userid     int       `json:"-" binding:"required"`
	Reviewerid int       `json:"reviewerid" `
	UpdatedAt  time.Time `json:"-"`
	CreateAt   time.Time `json:"-"`
}

//InviteData defining structure for binding data
type InviteData struct {
	Email string `json:"email" `
	Level string `json:"level"`
}

//MailData is used to send as payload for notification service
type MailData struct {
	Email string `json:"email"`
	Link  string `json:"link"`
	Task  string `json:"task"`
}

//EmployeeInfo is used to bind info for signup
type EmployeeInfo struct {
	PersonalInfo     PersonalInfo       `json:"personal_info"`
	UserAddress      []Address          `json:"user_address"`
	EmergencyContact []EmergencyContact `json:"emergency_contact"`
	Education        []Education        `json:"education"`
	Others           []Others           `json:"others"`
}

//DataInvite used for getting invite data
type DataInvite struct {
	Invite []InviteData `json:"invite"`
}

//AdminInvite is used to bind data
type AdminInvite struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

//Inviteadmin is used to bind data
type Inviteadmin struct {
	Invite []AdminInvite `json:"invite"`
}

//Registeradmin is used for binding data
type Registeradmin struct {
	Name     string `json:"name" binding:"required"`
	Contact  string `json:"contact" binding:"required"`
	Password string `json:"password" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

//=========================== for update ======================================

//UpdateEmployeeInfo is used to bind info for signup
type UpdateEmployeeInfo struct {
	Account          Accounts                 `json:"account_info"`
	PersonalInfo     UpdatePersonalInfo       `json:"personal_info"`
	UserAddress      []UpdateAddress          `json:"user_address"`
	EmergencyContact []UpdateEmergencyContact `json:"emergency_contact"`
	Education        []UpdateEducation        `json:"education"`
	Others           []UpdateOthers           `json:"others"`
}

//UpdatePersonalInfo : structure of personal information of employee
type UpdatePersonalInfo struct {
	ID             int       `json:"-"`
	Userid         int       `json:"-"`
	FathersName    string    `json:"fathers_name"`
	MothersName    string    `json:"mothers_name"`
	Gender         string    `json:"gender"`
	DateOfBirth    string    `json:"dateofbirth"`
	BloodGroup     string    `json:"blood_group"`
	MaritalStatus  string    `json:"marital_Status"`
	SpouseName     string    `json:"spouse_name"`
	SpouseContact  string    `json:"spouse_contact"`
	ProfileImgLink string    `json:"profile_img_link"`
	UpdatedAt      time.Time `json:"-"`
	CreateAt       time.Time `json:"-"`
}

//UpdateAddress is used to add address of user
type UpdateAddress struct {
	ID          int       `json:"-"`
	Userid      int       `json:"-"`
	AddressType string    `json:"address_type" ` //value can be current Address or permanent address
	Address     string    `json:"address"`
	City        string    `json:"city"`
	District    string    `json:"district"`
	State       string    `json:"state"`
	Pincode     int       `json:"pincode"`
	Country     string    `json:"country"`
	UpdatedAt   time.Time `json:"-"`
	CreateAt    time.Time `json:"-"`
}

//UpdateEmergencyContact : It is a structure to store emergency contact number of employees
type UpdateEmergencyContact struct {
	ID        int       `json:"-"`
	Userid    int       `json:"-"`
	Name      string    `json:"contact_name" `
	ContactNo string    `json:"contact_number"`
	UpdatedAt time.Time `json:"-"`
	CreateAt  time.Time `json:"-"`
}

//UpdateEducation is a structure to store
type UpdateEducation struct {
	ID int `json:"-" `
	//AccountID
	Userid int `json:"-"`
	//Name of the Degree/Class
	Name string `json:"name"`
	//Name of School/University
	Institution string `json:"institution"`
	//Location of School/Univeristy
	Location string `json:"location"`
	//Session
	Year string `json:"year"`
	//percentage scored in Class/Degree
	Percentage float64 `json:"percentage"`
	//uploaded marksheet link
	MarksheetLink pq.StringArray `json:"marksheet_link"  `
	//Uploaded Degree Link in case of University
	DegreeLink string    `json:"degree_link"`
	UpdatedAt  time.Time `json:"-"`
	CreateAt   time.Time `json:"-"`
}

//UpdateOthers : it is used to save extra certificates info of employee
type UpdateOthers struct {
	ID     int `json:"-" `
	Userid int `json:"-"`
	//Name of the Certificate
	Name string `json:"name"`
	//Uploaded Certificate Link
	CertificatesLink string    `json:"certificate_link" `
	UpdatedAt        time.Time `json:"-"`
	CreateAt         time.Time `json:"-"`
}

//Token is used to store token info
type SocialToken struct {
	ID int `json:"-" `

	Email string `json:"email"`
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string `json:"access_token"`

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string `json:"token_type"`

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string `json:"refresh_token"`

	// Expiry is the optional expiration time of the access token.
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry    time.Time `json:"expiry,omitempty"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

//InviteEmail is used to invite user
type InviteEmail struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	StartTime   int64    `json:"start" binding:"required"`
	EndTime     int64    `json:"end" binding:"required"`
	Email       []string `json:"email" binding:"required"`
}
