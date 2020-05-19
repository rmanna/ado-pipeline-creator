package fileutils

import (
	"fmt"

	"github.com/rmanna/ado-pipeline-creator/internal/config"
	"github.com/spf13/viper"
)

// ReadYamlConfig exported
func ReadYamlConfig(configPath string, configName string) (configuration config.Configurations) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	//var configuration config.Configurations
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	fmt.Printf("%#v\n", configuration)

	return configuration
}
