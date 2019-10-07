package gorm_bulk

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

type stringHelper struct{}

func StringHelper() *stringHelper {
	return &stringHelper{}
}

func (sh *stringHelper) ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func (sh *stringHelper) StructToMap(item interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(item)

	t := reflect.TypeOf(item)
	if t.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	result := make(map[string]interface{}, 0)

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		//fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
		if len(tag.Get("gorm")) > 0 && strings.Contains(tag.Get("gorm"), "json") {
			jsonValue, err := json.Marshal(valueField.Interface())
			if err != nil {
				return nil, err
			}
			result[typeField.Name] = string(jsonValue)
		} else {
			result[typeField.Name] = valueField.Interface()
		}
	}

	return result, nil
}

func (sh *stringHelper) SliceContains(items []string, value string) bool {
	for _, val := range items {
		if val == value {
			return true
		}
	}
	return false
}
