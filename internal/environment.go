package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const environmentDevelopment = "development"

func init() {
	godotenv.Load()
}

func getEnvironmentalVariable(variableName string) string {
	return os.Getenv(variableName)
}

func getEnvironmentalToken(variableName string, environmentalVariablePrefix string) string {
	environmentalVariable := getEnvironmentalVariable(variableName)

	if !strings.HasPrefix(environmentalVariable, environmentalVariablePrefix) {
		errorMessage := fmt.Sprintf("%s must have the prefix \"%s\".", variableName, environmentalVariablePrefix)
		fmt.Fprintf(os.Stderr, errorMessage)
	}
	return environmentalVariable
}

func getEnvironment() string {
	return getEnvironmentalVariable("APP_ENV")
}

func isDevelopmentEnvironment() bool {
	return getEnvironment() == environmentDevelopment
}

// GetSlackTokenApp returns the Slack bot token set in the
// env file
func GetSlackTokenApp() string {
	return getEnvironmentalToken("SLACK_TOKEN_APP", "xapp-")
}

// GetSlackTokenBot returns the Slack app token set in the
// env file
func GetSlackTokenBot() string {
	return getEnvironmentalToken("SLACK_TOKEN_BOT", "xoxb-")
}
