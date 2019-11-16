package gorm_bulk

import (
	"encoding/json"
	"reflect"
	"strings"
	"sync"
)

var (
	entityMapper *colMapper
)

type Entity interface{}

type (
	colMapper struct {
		mu              sync.Mutex
		tableColumnDict map[string][]string
	}
)

func init() {
	entityMapper = &colMapper{
		tableColumnDict: make(map[string][]string, 0),
	}
}

func Mapper() *colMapper {
	return entityMapper
}

func (cm *colMapper) getColumns(entity Entity) ([]string, error) {
	val := reflect.ValueOf(entity)
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	result := make([]string, 0)

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		result = append(result, StringHelper().ToSnakeCase(typeField.Name))
	}

	return result, nil
}

func (cm *colMapper) GetColumns(entity Entity) ([]string, error) {
	structName := structName(entity)
	rs := make([]string, 0)
	if columns, ok := entityMapper.tableColumnDict[structName]; ok {
		copy(rs, columns)
		return rs, nil
	}

	columns, err := cm.getColumns(entity)
	if err != nil {
		return nil, err
	}
	entityMapper.mu.Lock()
	defer entityMapper.mu.Unlock()

	entityMapper.tableColumnDict[structName] = columns
	copy(rs, columns)

	return rs, nil
}

func (cm *colMapper) GetValues(entity Entity) (map[string]interface{}, error) {
	val := reflect.ValueOf(entity)

	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	result := make(map[string]interface{}, 0)

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		column := StringHelper().ToSnakeCase(typeField.Name)

		if value := tag.Get("gorm"); len(value) > 0 && strings.Contains(value, "json") {
			jsonValue, err := json.Marshal(valueField.Interface())
			if err != nil {
				return nil, err
			}

			result[column] = string(jsonValue)
		} else {
			result[column] = valueField.Interface()
		}
	}

	return result, nil
}

func structName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
