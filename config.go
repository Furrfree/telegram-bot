package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Token                          string
	AdmissionGroupId               string
	GroupId                        string
	RulesMessageUrl                string
	PresentationTemplateMessageUrl string
}

func getEnvVariable(name string) string {
	envVar := os.Getenv(name)
	if envVar == "" {
		log.Fatalf("Env variable %s required", name)
	}

	return os.Getenv(name)
}

func getConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		Token:                          getEnvVariable("TOKEN"),
		AdmissionGroupId:               getEnvVariable("ADMISSION_GROUP_ID"),
		GroupId:                        getEnvVariable("GROUP_ID"),
		RulesMessageUrl:                getEnvVariable("RULES_MESSAGE_URL"),
		PresentationTemplateMessageUrl: getEnvVariable("PRESENTATION_TEMPLATE_MESSAGE_URL"),
	}

}
