package config

import (
	"os"
	"testing"
)

func TestConfigurationWithEnv(t *testing.T) {

	// test case 1-> call with out set any variable
	SetConfig()
	DBConfig()
	ConfigurationWithEnv()
	ConfigurationWithToml("../example.toml")

	//test case 2 -> set some variable
	os.Setenv("AUTH_DB_PORT", "8080")

	ConfigurationWithEnv()

	ConfigurationWithToml("../example.toml")
	SetConfig()
}
