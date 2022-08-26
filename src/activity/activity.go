package activity

import (
	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
)

// RecordActivity is a method use to record activity of a user in activity table
func RecordActivity(activity database.Activities) {
	//handler panic and Alerts
	defer util.Panic()
	// connecting to db
	db := config.DB

	// recording users activities
	db.Create(&activity)
}

// GetLoginActivities is a method used to get login activities of a user
func GetLoginActivities(email string) ([]database.Activities, error) {
	//handler panic and Alerts
	defer util.Panic()
	// connecting to db
	db := config.DB

	// fetching activities of user from activities table
	var activities []database.Activities
	db.Where("email=?", email).Order("timestamp desc").Limit(5).Find(&activities)
	return activities, nil
}
