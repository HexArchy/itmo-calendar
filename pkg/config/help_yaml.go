package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
)

// helpYaml displays configuration help in YAML format.
func helpYaml(dict Dict) {
	res := map[string]interface{}{}

	for _, prop := range dict.Sorted() {
		path := strings.Split(prop.Name, ".")
		deepSet(res, path, typedYamlValue(prop.Type, prop.CurrentVal))
	}

	r, err := yaml.Marshal(res)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(r))
	os.Exit(0)
}

// deepSet sets a value in a nested map following the specified path.
func deepSet(to map[string]interface{}, path []string, value interface{}) {
	if len(path) == 1 {
		to[path[0]] = value

		return
	}

	cur := path[0]
	next := path[1:]

	if _, ok := to[cur]; !ok {
		to[cur] = map[string]interface{}{}
	}

	if m, ok := to[cur].(map[string]interface{}); ok {
		deepSet(m, next, value)

		return
	}
}

// typedYamlValue converts a string value to an appropriate type for YAML output.
func typedYamlValue(t string, val string) interface{} {
	if strings.HasSuffix(t, "Slice") {
		items := strings.Split(val, ",")
		res := []interface{}{}

		if strings.Contains(t, "int") {
			for _, v := range items {
				res = append(res, cast.ToInt64(v))
			}

			return res
		}

		if strings.Contains(t, "float") {
			for _, v := range items {
				res = append(res, cast.ToFloat64(v))
			}

			return res
		}

		return items
	}

	if strings.Contains(t, "int") {
		return cast.ToInt64(val)
	}

	if strings.Contains(t, "float") {
		return cast.ToFloat64(val)
	}

	return val
}
