// +build !test

package health

import (
	"errors"

	"git.xenonstack.com/xs-onboarding/accounts/config"
)

// ServiceHealth is a method to check service each components health
func ServiceHealth() error {

	// connecting to db
	db := config.DB
	var ok string
	err := db.DB().QueryRow("SELECT 1").Scan(&ok)
	if err != nil {
		return errors.New("postgres db is not working")
	}
	return nil
}
