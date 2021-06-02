package tbam

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var config *Config

// Configs exported
type Config struct {
	Log       LogConfig
	Webserver WebserverConfig
	Cookie    CookieConfig
	AuthFile  string
	Validity  ValidityConfig
}

type ValidityConfig struct {
	Session    int
	Credential int
}

// LogConfig exported
type LogConfig struct {
	Level string
}

// Recaptcha exported
type SmtpConfig struct {
	Server           string
	From             string
	DefaultRecipient string
	Login            string
	Password         string
}

// WebserverConfig exported
type WebserverConfig struct {
	Protocol string
	Hostname string
	Port     int
	Debug    bool
}

// CookieConfig exported
type CookieConfig struct {
	Name   string
	Secret string
	Domain string
}

func readConfig(configName, configType string) (*Config, error) {
	//log.Info("Reading configuration")
	// Set the file name of the configurations file
	viper.SetConfigName(configName)

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType(configType)
	//var configuration c.Configs

	if err := viper.ReadInConfig(); err != nil {
		//log.Error("Error reading config file, %s", err)
		return config, err
	}

	// Set undefined variables
	//viper.SetDefault("database.dbname", "test_db")

	err := viper.Unmarshal(&config)
	if err != nil {
		//log.Error("Unable to decode into struct, %v", err)
		return config, err
	}

	return config, nil
}

// NewGlobalConfig creates a new global config from file or arguments
func NewGlobalConfig() *Config {

	var err error
	// Init config
	config, err := readConfig("config", "yml")
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	return config
}
