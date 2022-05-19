package config

import (
	"fmt"
	"github.com/NubeIO/rubix-assist/pkg/helpers/homedir"
	"github.com/NubeIO/rubix-assist/pkg/logger"
	"github.com/spf13/viper"
)

var Config *Configuration

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
	Path     PathConfiguration
}

// Setup initialize configuration
func Setup() error {
	var configuration *Configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
	}
	viper.AddConfigPath(home + "/.updater")

	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Error reading config file, %s", err)
		fmt.Println(err)
		//return err
	}

	err = viper.Unmarshal(&configuration)
	if err != nil {
		logger.Errorf("Unable to decode into struct, %v", err)
		fmt.Println(err)
	}
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("database.driver", "sqlite")
	viper.SetDefault("database.name", "/tmp/updater.db")

	Config = configuration
	return nil
}

// GetConfig helps you to get configuration data
func GetConfig() *Configuration {
	return Config
}
