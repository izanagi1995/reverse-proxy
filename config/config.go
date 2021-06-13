package config

import "github.com/spf13/viper"

func ReadFromFile() (map[string][]string, error) {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var config = make(map[string][]string)
	err = v.UnmarshalExact(&config)
	return config, err
}
