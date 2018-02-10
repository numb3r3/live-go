package config

import (
	"github.com/spf13/viper"
)

func readConfig(filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	//v.AddConfigPath("$HOME/config")

	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}
