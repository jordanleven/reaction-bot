package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type error interface {
	Error() string
}

func getEnvironmentalVariable(variableName string, environmentalVariablePrefix string) string {
	environmentalVariable := os.Getenv(variableName)

	if !strings.HasPrefix(environmentalVariable, environmentalVariablePrefix) {
		errorMessage := fmt.Sprintf("%s must have the prefix \"%s\".", variableName, environmentalVariablePrefix)
		fmt.Fprintf(os.Stderr, errorMessage)
	}
	return environmentalVariable
}

func init() {
	godotenv.Load()
}

// GetSlackTokenApp returns the Slack bot token set in the
// env file
func GetSlackTokenApp() string {
	return getEnvironmentalVariable("SLACK_TOKEN_APP", "xapp-")
}

// GetSlackTokenBot returns the Slack app token set in the
// env file
func GetSlackTokenBot() string {
	return getEnvironmentalVariable("SLACK_TOKEN_BOT", "xoxb-")
}
