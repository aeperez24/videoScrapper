package service

import "github.com/spf13/viper"

func LoadConfig(path string) (config AppConfiguration, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app.yaml")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
