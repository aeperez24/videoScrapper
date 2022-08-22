package service

import "github.com/spf13/viper"

func LoadConfig(path string) (*AppConfiguration, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app.yaml")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &AppConfiguration{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, nil
	}

	return config, nil
}
