//+build !test

package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/src/nats"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/routes"
	"git.xenonstack.com/xs-onboarding/accounts/src/scheduler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func main() {
	//conf flag is used to set the way user want configuration of service to be set
	// It can be via TOML or ENVIRONMENT
	conf := flag.String("conf", "environment", "set configuration from toml file or environment variables")
	//file flag is used to give path to toml file if conf is set with TOML
	file := flag.String("file", "", "set path of toml file")
	//Below function parses flags used as command line arguement
	flag.Parse()

	//check if conf is set to environment if yes call function to set configuration using environment
	if *conf == "environment" {
		log.Println("environment")
		config.ConfigurationWithEnv()
	} else if *conf == "toml" { //if conf is set to TOML set configuration using file
		log.Println("toml")
		if *file == "" {
			log.Println("Please pass toml file path")
			os.Exit(1)
		} else { // set configuration by calling below function
			err := config.ConfigurationWithToml(*file)
			if err != nil {
				log.Println("Error in parsing toml file")
				log.Println(err)
				os.Exit(1)
			}
		}
	} else {
		log.Println("Please pass valid arguments, conf can be set as toml or environment")
		os.Exit(1)
	}
	//If environement is set to production
	if config.Conf.Service.Environment != "production" {
		// removing info file if any.
		_ = os.Remove("info.txt")

		// creating and opening info.txt file for writting logs
		file, err := os.OpenFile("info.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		// changing default writer of gin to file and std output
		gin.DefaultWriter = io.MultiWriter(file, os.Stdout)

		// setting output for logs this will writes all logs to file
		log.SetOutput(gin.DefaultWriter)
		// writing log to check all in working
		log.Print("Logging to a file in Go!")
	}
	//-------------------------------- DATABASE CONNECTION ------------------------------------------
	//Set Database Connection
	//create database
	database.CreateDatabase()

	//database config
	dbConfig := config.DBConfig()
	//number of ideal connections
	var ideal int
	idealStr := config.Conf.Database.Ideal
	if idealStr == "" {
		ideal = 50
	} else {
		ideal, _ = strconv.Atoi(idealStr)
	}
	// connecting db using connection string
	db, err := gorm.Open("postgres", dbConfig)
	if err != nil {
		log.Println(err)
		return
	}
	// close db instance whenever whole work completed
	defer db.Close()
	//set maximum number of connections that can be made with Database
	db.DB().SetMaxIdleConns(ideal)
	//set maximum lifetime of a connection
	db.DB().SetConnMaxLifetime(1 * time.Hour)
	//set global db connection
	config.DB = db
	//create auth-team database tables
	go database.CreateDatabaseTables()

	//--------------------------------- CORS CONFIGURATION ----------------------------------------------
	//allowing CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AddAllowMethods("DELETE")
	// initialize gin router
	router := gin.Default()
	router.Use(cors.New(corsConfig))

	//-----------------------------------NATS Connection------------------------------------------
	nats.InitConnection()
	go nats.Subscribe()
	//start scheduler
	go scheduler.Start()
	//---------------------------------- VIEW ALL THE REGISTERED ROUTES --------------------------------------------
	// index handler to view all registered routes
	router.GET("/", func(c *gin.Context) {
		type finalPath struct {
			Method string
			Path   string
		}

		data := router.Routes()
		finalPaths := make([]finalPath, 0)

		for i := 0; i < len(data); i++ {
			finalPaths = append(finalPaths, finalPath{
				Path:   data[i].Path,
				Method: data[i].Method,
			})
		}
		c.JSON(200, gin.H{
			"routes": finalPaths,
		})
	})

	// service specific routes
	routes.Routes(router)

	// run the service
	router.Run(":" + config.Conf.Service.Port)
}
