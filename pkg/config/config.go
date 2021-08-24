package config

import "github.com/spf13/viper"

const (
	Production  = "production"
	Development = "development"
	Test        = "test"
)

const (
	BucketName = "fileapi"
	AwsRegion  = "sa-east-1"
)

func init() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}

	viper.AutomaticEnv()

	// Project defaults
	viper.SetDefault("http_port", 5555)
	viper.SetDefault("base_url", "https://rubbioli.com/fileapi/graphql")
}

func Environment() string {
	switch viper.GetString("environment") {
	case "development":
		return Development
	case "test":
		return Test
	default:
		return Production
	}
}

func Port() int {
	return viper.GetInt("http_port")
}

func BaseURL() string {
	return viper.GetString("base_url")
}

func AWSKey() string {
	return viper.GetString("aws_key")
}

func AWSSecret() string {
	return viper.GetString("aws_secret")
}
