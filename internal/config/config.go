package config

import "github.com/spf13/viper"

var Cfg Config

type Config struct {
	D            string
	Token        string
	TouchChannel string
}

func LoadConfig() error {
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&Cfg)
	if err != nil {
		return err
	}

	return nil

}
