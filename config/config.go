package config

import (
	"fmt"

	db "github.com/aswa2ds/as-cron-db"
	"github.com/spf13/viper"
)

var Config = config{}

type config struct {
	DatabaseConfig db.Config
}

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	viper.Unmarshal(&Config)
}
