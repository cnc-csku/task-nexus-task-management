package config

import (
	"log"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	// Server
	PORT string `mapstructure:"PORT"`

	// Database
	MongoURI string `mapstructure:"MONGO_URI"`
}

func getMapstructureTags(v interface{}) []string {
	typ := reflect.TypeOf(v)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	var tags []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if tag, ok := field.Tag.Lookup("mapstructure"); ok {
			tags = append(tags, tag)
		}
	}
	return tags
}

func NewConfig() *Config {
	config := Config{}

	// Set the .env file and read environment variables
	viper.SetConfigFile(".env")

	// Attempt to read the .env file
	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️  .env file not found or cannot be read, using environment variables")

		// Bind environment variables
		envs := getMapstructureTags(config)
		for _, env := range envs {
			viper.MustBindEnv(env)
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Println("❌ Unable to decode into struct:", err)
	}

	return &config
}
