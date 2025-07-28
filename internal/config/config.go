package config

import (
	"log"
	"reflect"

	"github.com/spf13/viper"
)

var Config AppConfig

type AppConfig struct {
	AppPort              string `mapstructure:"APP_PORT"`
	SecretTokenKey       string `mapstructure:"SECRET_TOKEN_KEY"`
	AccessTokenDuration  int    `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration int    `mapstructure:"REFRESH_TOKEN_DURATION"`
	DBHost               string `mapstructure:"DB_HOST"`
	DBName               string `mapstructure:"DB_NAME"`
	DBUsername           string `mapstructure:"DB_USERNAME"`
	DBPassword           string `mapstructure:"DB_PASSWORD"`
	DBMode               string `mapstructure:"DB_MODE"`
}

func InitConfig() {

	v := reflect.ValueOf(&AppConfig{})
	t := v.Elem().Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		env := field.Tag.Get("mapstructure")
		if env == "" {
			continue
		}
		viper.BindEnv(env)
	}

	viper.AutomaticEnv()

	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

}
