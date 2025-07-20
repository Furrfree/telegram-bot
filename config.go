package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Token                          string
	AdmissionGroupId               int
	GroupId                        int
	RulesMessageUrl                string
	PresentationTemplateMessageUrl string
}

func getIntEnvVariable(name string) int {
	envVar := os.Getenv(name)
	if envVar == "" {
		log.Fatalf("Env variable %s required", name)
	}

	intValue, err := strconv.Atoi(envVar)
	if err != nil {
		log.Fatalf("Env variable %s must be integer", name)
	}

	return intValue
}

func getStringEnvVariable(name string) string {
	envVar := os.Getenv(name)
	if envVar == "" {
		log.Fatalf("Env variable %s required", name)
	}

	return envVar
}

func getConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		Token:                          getStringEnvVariable("TOKEN"),
		AdmissionGroupId:               getIntEnvVariable("ADMISSION_GROUP_ID"),
		GroupId:                        getIntEnvVariable("GROUP_ID"),
		RulesMessageUrl:                getStringEnvVariable("RULES_MESSAGE_URL"),
		PresentationTemplateMessageUrl: getStringEnvVariable("PRESENTATION_TEMPLATE_MESSAGE_URL"),
	}

}
