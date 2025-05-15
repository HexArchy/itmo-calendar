package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var _safeEnvValue = regexp.MustCompile(`^[A-Za-z0-9.\-,]+$`)

// helpEnv displays configuration help in environment variable format.
func helpEnv(dict Dict) {
	for _, prop := range dict.Sorted() {
		fmt.Fprintf(os.Stdout, "%s=%s\n", envName(prop), envValue(prop))
	}

	os.Exit(0)
}

// envName converts a property name to an environment variable name.
func envName(prop *Prop) string {
	return _envReplacer.Replace(strings.ToUpper(prop.Name))
}

// envValue formats a property value for use in environment variables.
func envValue(prop *Prop) string {
	if prop.CurrentVal == "" {
		return prop.CurrentVal
	}

	if _safeEnvValue.MatchString(prop.CurrentVal) {
		return prop.CurrentVal
	}

	v := strings.ReplaceAll(prop.CurrentVal, "'", "'\"'\"'")

	return "'" + v + "'"
}
