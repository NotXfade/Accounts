package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	"github.com/nats-io/nats.go"
)

// Config is a structure for configuration
type Config struct {
	Database     Database
	Service      Service
	Admin        Admin
	JWT          JWT
	NatsServer   NatsServer
	CareerPortal CareerPortal
	Google       Google
}

// Google oAUth configuration structure
type Google struct {
	ClientID  string
	ClientKey string
	Scopes    string
	Redirect  string
}

//NatsServer : for nats connection parameters
type NatsServer struct {
	URL      string
	Token    string
	Username string
	Password string
	Subject  string
	Queue    string
}

// JWT is structure for jwt token specific configuration
type JWT struct {
	PrivateKey    string
	JWTExpireTime time.Duration
}

// Database is a structure for postgres database configuration
type Database struct {
	Name  string
	Host  string
	Port  string
	User  string
	Pass  string
	Ssl   string
	Ideal string
}

// Service is a structure for service specific related configuration
type Service struct {
	Port              string
	Environment       string
	Build             string
	IsLogoutOthers    string
	VerifyLinkTimeout int64
	InviteLinkTimeout int64
	HostAddr          string
}

//CareerPortal is used for careerportal related config
type CareerPortal struct {
	FrontEndAddress string
	Username        string
	Password        string
}

// falseStr is a constant to remove duplicacy code
const falseStr string = "false"

// Admin is a structure for admin account credentials
type Admin struct {
	Email string
	Pass  string
}

var (
	// Conf is a global variable for configuration
	Conf Config
	// TomlFile is a global variable for toml file path
	TomlFile string
	// DB Database client
	DB *gorm.DB

	//NC for nats connection
	NC *nats.Conn
)

// ConfigurationWithEnv is a method to initialize configuration with environment variables
func ConfigurationWithEnv() {
	// postgres database configuration
	Conf.Database.Host = os.Getenv("AUTH_DB_HOST")
	Conf.Database.Port = os.Getenv("AUTH_DB_PORT")
	Conf.Database.User = os.Getenv("AUTH_DB_USER")
	Conf.Database.Pass = os.Getenv("AUTH_DB_PASS")
	Conf.Database.Name = os.Getenv("AUTH_DB_NAME")
	Conf.Database.Ideal = os.Getenv("AUTH_DB_IDEAL_CONNECTIONS")
	Conf.Database.Ssl = "disable"

	// admin account credentials configuration
	Conf.Admin.Email = os.Getenv("AUTH_ADMIN_EMAIL")
	Conf.Admin.Pass = os.Getenv("AUTH_ADMIN_PASS")

	Conf.Service.Environment = os.Getenv("ENVIRONMENT")
	Conf.Service.Build = os.Getenv("BUILD_IMAGE")
	Conf.Service.HostAddr = os.Getenv("HOST_ADDRESS")
	//service specific configuration
	//default value of is logout other is true
	Conf.Service.IsLogoutOthers = "true"
	//Conf.Service.VerifyLinkTimeout = os.Getenv("AUTH_SERVICE_VERIFY_LINK_TIMEOUT")
	//Conf.Service.InviteLinkTimeout = os.Getenv("AUTH_SERVICE_INVITE_LINK_TIMEOUT")
	if os.Getenv("IS_LOGOUT_OTHER") == falseStr {
		Conf.Service.IsLogoutOthers = falseStr
	}
	if os.Getenv("AUTH_SERVICE_PORT") != "" {
		Conf.Service.Port = os.Getenv("AUTH_SERVICE_PORT")
	} else {
		Conf.Service.Port = "8000"
	}

	// set constants
	//JWT Token Timeout in minutes
	Conf.JWT.JWTExpireTime = time.Minute * 30
	Conf.JWT.PrivateKey = os.Getenv("PRIVATE_KEY")
	// constants
	//Link Expiration time in seconds
	if Conf.Service.VerifyLinkTimeout == 0 {
		Conf.Service.VerifyLinkTimeout = 1800
	}
	if Conf.Service.InviteLinkTimeout == 0 {
		Conf.Service.InviteLinkTimeout = 86400
	}
	if Conf.Service.Environment == "development" {
		Conf.NatsServer.Subject = "xsonboarding.develop"
		Conf.NatsServer.Queue = "accounts-develop"
	} else if Conf.Service.Environment == "uat" {
		Conf.NatsServer.Subject = "xsonboarding.uat"
		Conf.NatsServer.Queue = "accounts-uat"
	} else if Conf.Service.Environment == "production" {
		Conf.NatsServer.Subject = "xsonboarding"
		Conf.NatsServer.Queue = "xsonboarding-accounts"
	}
	//nats server
	Conf.NatsServer.URL = os.Getenv("NATS_URL")
	Conf.NatsServer.Token = os.Getenv("NATS_TOKEN")
	Conf.NatsServer.Username = os.Getenv("NATS_USERNAME")
	Conf.NatsServer.Password = os.Getenv("NATS_PASSWORD")
	Conf.Google.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	Conf.Google.ClientKey = os.Getenv("GOOGLE_CLIENT_KEY")
	Conf.Google.Scopes = os.Getenv("GOOGLE_SCOPES")
	Conf.Google.Redirect = os.Getenv("GOOGLE_REDIRECT")
	Conf.CareerPortal.FrontEndAddress = os.Getenv("CAREER_PORTAL_FRONTENDADDRESS")
	Conf.CareerPortal.Username = os.Getenv("CAREER_PORTAL_USERNAME")
	Conf.CareerPortal.Password = os.Getenv("CAREER_PORTAL_PASSWORD")
}

// ConfigurationWithToml is a method to initialize configuration with toml file
func ConfigurationWithToml(filePath string) error {
	// set varible as file path if configuration is done using toml
	TomlFile = filePath
	log.Println(filePath)
	// parse toml file and save data config structure
	_, err := toml.DecodeFile(filePath, &Conf)
	if err != nil {
		log.Println(err)
		return err
	}

	if Conf.Service.Port == "" {
		Conf.Service.Port = "8000"
	}
	Conf.Database.Ssl = "disable"
	Conf.Service.Build = os.Getenv("BUILD_IMAGE")

	// set constants
	//JWT Token Timeout in minutes
	Conf.JWT.JWTExpireTime = time.Minute * 30

	//service specific configuration
	//default value of is logout other is true
	Conf.Service.IsLogoutOthers = "true"
	if os.Getenv("IS_LOGOUT_OTHER") == falseStr {
		Conf.Service.IsLogoutOthers = falseStr
	}
	// constants
	//Link Expiration time in seconds

	if Conf.Service.VerifyLinkTimeout == 0 {
		Conf.Service.VerifyLinkTimeout = 600
	}
	if Conf.Service.InviteLinkTimeout == 0 {
		Conf.Service.InviteLinkTimeout = 604800
	}
	if Conf.Service.Environment == "development" {
		Conf.NatsServer.Subject = "xsonboarding.develop"
		Conf.NatsServer.Queue = "accounts-develop"
	} else if Conf.Service.Environment == "uat" {
		Conf.NatsServer.Subject = "xsonboarding.uat"
		Conf.NatsServer.Queue = "accounts-uat"
	} else if Conf.Service.Environment == "production" {
		Conf.NatsServer.Subject = "xsonboarding"
		Conf.NatsServer.Queue = "xsonboarding-accounts"
	}
	return nil
}

// SetConfig is a method to re-intialise configuration at runtime
func SetConfig() {
	if TomlFile == "" {
		ConfigurationWithEnv()
	} else {
		ConfigurationWithToml(TomlFile)
	}
}

// DBConfig is a method that return postgres database connection string
func DBConfig() string {
	//again reset the config if any changes in toml file or environment variables
	//	SetConfig()
	// creating postgres connection string
	str := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		Conf.Database.Host,
		Conf.Database.Port,
		Conf.Database.User,
		Conf.Database.Pass,
		Conf.Database.Name,
		Conf.Database.Ssl)
	return str
}
