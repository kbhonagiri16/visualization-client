package config

import (
	"fmt"
)

// ParseError represents error with parsing application configuration
// from config file, env variables and console parameters
type ParseError struct {
	name                string
	configFileProperty  string
	configFileSection   string
	environmentVariable string
	consoleFlag         string
}

// NewParseError initializes ParseError
func NewParseError(name string, configFileProperty string,
	configFileSection string, environmentVariable string,
	consoleFlag string) *ParseError {
	return &ParseError{
		name:                name,
		configFileProperty:  configFileProperty,
		configFileSection:   configFileSection,
		environmentVariable: environmentVariable,
		consoleFlag:         consoleFlag,
	}
}

func (err *ParseError) Error() string {
	errorMessage := fmt.Sprintf("Error in getting config value %s use ",
		err.name)
	if err.configFileProperty != "" && err.configFileSection != "" {
		errorMessage += fmt.Sprintf(
			"property '%s' in [%s] config file section, ",
			err.configFileProperty,
			err.configFileSection,
		)
	}
	if err.environmentVariable != "" {
		errorMessage += fmt.Sprintf("%s env variable, ",
			err.environmentVariable)
	}
	if err.consoleFlag != "" {
		errorMessage += fmt.Sprintf("%s command line flag",
			err.consoleFlag)
	}
	return errorMessage
}
