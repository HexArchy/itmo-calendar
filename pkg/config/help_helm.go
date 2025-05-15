package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// helpHelm displays configuration help in Helm format.
func helpHelm(dict Dict) {
	for _, prop := range dict.Sorted() {
		fmt.Printf("%s: %s\n", envName(prop), yamlValue(prop))
	}

	os.Exit(0)
}

// yamlValue returns a YAML-formatted string representation of a property's value.
func yamlValue(prop *Prop) string {
	v, _ := yaml.Marshal(prop.CurrentVal)

	return strings.TrimSpace(string(v))
}
