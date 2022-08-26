package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/src/methods"
)

//CreateDatabaseTables funtion for create Table
func CreateDatabaseTables() {
	// connecting db using connection string
	db := config.DB

	// creating all tables one by one but firstly checking whether table exists or not
	if !(db.HasTable(Accounts{})) {
		db.CreateTable(Accounts{})
		//creating admin account
		adminAcc := InitAdminAccount()
		db.Create(&adminAcc)
	}
	if !(db.HasTable(Activities{})) {
		db.CreateTable(Activities{})
	}
	if !(db.HasTable(Departments{})) {
		db.CreateTable(Departments{})
	}
	if !(db.HasTable(Tokens{})) {
		db.CreateTable(Tokens{})
	}
	if !(db.HasTable(Education{})) {
		db.CreateTable(Education{})
	}
	if !(db.HasTable(Others{})) {
		db.CreateTable(Others{})
	}
	if !(db.HasTable(PersonalInfo{})) {
		db.CreateTable(PersonalInfo{})
	}
	if !(db.HasTable(EmergencyContact{})) {
		db.CreateTable(EmergencyContact{})
	}
	if !(db.HasTable(PandCData{})) {
		db.CreateTable(PandCData{})
	}
	if !(db.HasTable(ActiveSessions{})) {
		db.CreateTable(ActiveSessions{})
	}
	if !(db.HasTable(Address{})) {
		db.CreateTable(Address{})
	}
	if !(db.HasTable(Reviewer{})) {
		db.CreateTable(Reviewer{})
	}
	if !(db.HasTable(SocialToken{})) {
		db.CreateTable(SocialToken{})
	}
	// Database migration
	db.AutoMigrate(&Accounts{},
		&Activities{},
		&Tokens{},
		&Education{},
		&Others{},
		&PersonalInfo{},
		&EmergencyContact{},
		&PandCData{},
		&ActiveSessions{},
		&Address{},
		&Reviewer{},
		&SocialToken{},
		&Departments{},
	)
	//db=db.Debug()
	db.Exec("ALTER TABLE educations DROP CONSTRAINT educations_userid_key;")
	db.Exec("ALTER TABLE addresses DROP CONSTRAINT addresses_userid_key;")
	// Add foreignKeys
	db.Model(&Education{}).AddForeignKey("userid", "accounts(id)", "CASCADE", "CASCADE")
	db.Model(&Others{}).AddForeignKey("userid", "accounts(id)", "CASCADE", "CASCADE")
	db.Model(&PersonalInfo{}).AddForeignKey("userid", "accounts(id)", "CASCADE", "CASCADE")
	db.Model(&EmergencyContact{}).AddForeignKey("userid", "accounts(id)", "CASCADE", "CASCADE")
	db.Model(&PandCData{}).AddForeignKey("userid", "accounts(id)", "CASCADE", "CASCADE")
	db.Model(&ActiveSessions{}).AddForeignKey("userid", "accounts(id)", "CASCADE", "CASCADE")
	db.Model(&Address{}).AddForeignKey("userid", "accounts(id)", "CASCADE", "CASCADE")
	db.Model(&Reviewer{}).AddForeignKey("userid", "accounts(id)", "CASCADE", "CASCADE")

	// //Reviewers view
	// query := `create extension dblink;`
	// err := db.Exec(query).Error
	// log.Println("err", err)
	// query = `create view internslist as SELECT name,email,level,team,department,account_status,role,reviewerid,slug,scores, total_scores,created_at,profile_img_link,form_filled_at
	// FROM   dblink('dbname=xsonboarding_module_tracker user=xsonboarding password=X$onbo@rding','select m.userid,m.slug,f.scores,f.total_scores,f.created_at from interns_modules_registries as m left join interns_feedback_forms as f on f.module_id = m.id')
	// AS     tb2(userid text, slug text, scores numeric,total_scores numeric,form_filled_at timestamptz)
	// right join(
	// 	select a2.id,a2.name,a2.email,a2.level,a2.account_status , a2.role,a2.created_at, r2.reviewerid, pcd.team,pcd.department,pi2.profile_img_link from accounts as a2
	// 	left join reviewers as r2 on  a2.id = r2.userid
	// 	left join personal_infos  as pi2 on a2.id = pi2.userid
	// 	left join pand_c_data as pcd on a2.id = pcd.userid
	// )
	// as acc on acc.email = tb2.userid`
	// err = db.Exec(query).Error
	// log.Println(err)

	//Reviewers view
	query := "create view reviewerslist as select a2.id ,a2.name,a2.email,a2.level,a2.account_status , a2.role, r2.reviewerid from accounts as a2 inner join reviewers as r2 on r2.userid = a2.id "
	db.Exec(query)
}

// CreateDatabase Initializing Database
func CreateDatabase() error {
	// connecting with postgres database root db
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Conf.Database.Host,
		config.Conf.Database.Port,
		config.Conf.Database.User,
		config.Conf.Database.Pass,
		"postgres", config.Conf.Database.Ssl))

	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()
	db = db.Debug()
	log.Println(config.Conf.Database.Name)
	// executing create database query.
	db.Exec(fmt.Sprintf("create database %s;", config.Conf.Database.Name))
	return nil
}

// InitAdminAccount is a function used to create admin account
func InitAdminAccount() Accounts {

	// fetching info from env variables
	adminEmail := config.Conf.Admin.Email
	if adminEmail == "" {
		adminEmail = "admin@xenonstack.com"
	}
	adminPass := config.Conf.Admin.Pass
	if adminPass == "" {
		adminPass = "admin"
	}
	// return struct with details of admin
	return Accounts{
		Name:          "Admin",
		Password:      methods.HashForNewPassword(adminPass),
		Email:         adminEmail,
		ContactNo:     "",
		Role:          "admin",
		AccountStatus: "active",
		Timestamp:     time.Now().Unix(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
