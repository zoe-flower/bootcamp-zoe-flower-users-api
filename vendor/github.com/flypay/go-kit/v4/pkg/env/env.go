package env

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetEnv provides a default override for environment vars as a second argument
func Get(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetOrError will attempt to get the value from the environment variable and
// return an error if not available
func GetOrError(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}
	return "", fmt.Errorf("%s is empty", key)
}

// GetAppEnvironment returns the app environment.
func GetAppEnvironment() string {
	return Get("CONFIG_APP_ENVIRONMENT", "")
}

// IsLocal returns whether or not we are running the app locally.
func IsLocal() bool {
	return GetAppEnvironment() == "local"
}

// IsDevelopment returns whether or not we are running the app in development.
func IsDevelopment() bool {
	return GetAppEnvironment() == "development"
}

// IsStaging returns whether or not we are running the app in staging.
func IsStaging() bool {
	return GetAppEnvironment() == "staging"
}

// IsProduction returns whether or not we are running the app in production.
func IsProduction() bool {
	return GetAppEnvironment() == "production"
}

// ShouldRun determines if a chunk of the app should run
func ShouldRun(API, runOnlyAPI string) bool {
	return runOnlyAPI == "" || runOnlyAPI == API
}

// GetAppName takes the configured app name OR the working directories parent base name
func GetAppName() (appName string) {
	appName = os.Getenv("CONFIG_APP_NAME")

	// If there's no app name set, set it to the parent directory name
	if appName == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		appName = filepath.Base(wd)
	}

	return
}

func GetQueueName() string {
	return Get("CONFIG_EVENT_BUS_SQS_QUEUE_NAME", "")
}

func GetQueueURLByName(qname string) (string, error) {
	var (
		accountID, region string
		err               error
	)
	if accountID, err = GetOrError("CONFIG_AWS_ACCOUNT_ID"); err != nil {
		return "", err
	}
	if region, err = GetOrError("AWS_REGION"); err != nil {
		return "", err
	}

	return fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", region, accountID, qname), nil
}

func GetQueueURL() (string, error) {
	var (
		queueName string
		err       error
	)
	if queueName, err = GetOrError("CONFIG_EVENT_BUS_SQS_QUEUE_NAME"); err != nil {
		return "", err
	}
	return GetQueueURLByName(queueName)
}

func GetTopic() string {
	return Get("CONFIG_EVENT_BUS_SQS_TOPIC_NAME", "")
}

// GetTeam returns the team name configured in the environment variable.
// If the default "not-set" is found (which is the current helm default)
// return an empty string.
func GetTeam() (team string) {
	if team = os.Getenv("CONFIG_APP_TEAM"); team == "not-set" {
		team = ""
	}
	return
}
