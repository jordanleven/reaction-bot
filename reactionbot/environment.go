package reactionbot

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func getEnvironmentalToken(variableName string, environmentalVariablePrefix string) string {
	environmentalVariable := os.Getenv(variableName)

	if !strings.HasPrefix(environmentalVariable, environmentalVariablePrefix) {
		errorMessage := fmt.Sprintf("%s must have the prefix \"%s\".", variableName, environmentalVariablePrefix)
		log.Fatal(errorMessage)
	}
	return environmentalVariable
}

func init() {
	godotenv.Load()
}

// GetSlackTokenApp returns the Slack bot token set in the
// env file
func GetSlackTokenApp(tokenName string) string {
	return getEnvironmentalToken(tokenName, "xapp-")
}

// GetSlackTokenBot returns the Slack app token set in the
// env file
func GetSlackTokenBot(tokenName string) string {
	return getEnvironmentalToken(tokenName, "xoxb-")
}
