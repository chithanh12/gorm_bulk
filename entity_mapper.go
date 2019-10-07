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
		mu         sync.Mutex
		dictionary map[string]EntityColumn
	}

	EntityColumn struct {
		Inserts []string
		Updates []string
	}
)

func init() {
	entityMapper = &colMapper{
		dictionary: map[string]EntityColumn{},
	}
}

func Mapper() *colMapper {
	return entityMapper
}

func (cm *colMapper) getColumns(entity Entity) (EntityColumn, error) {
	val := reflect.ValueOf(entity)
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	result := EntityColumn{
		Inserts: make([]string, 0),
		Updates: make([]string, 0),
	}

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		if value := tag.Get("column"); len(value) > 0 {
			if strings.Contains(value, "insert") {
				result.Inserts = append(result.Inserts, StringHelper().ToSnakeCase(typeField.Name))
			}

			if strings.Contains(value, "update") {
				result.Updates = append(result.Updates, StringHelper().ToSnakeCase(typeField.Name))
			}
		}
	}

	return result, nil
}

func (cm *colMapper) GetInsertColumns(entity Entity) []string {
	structName := structName(entity)
	if mapper, ok := entityMapper.dictionary[structName]; ok {
		return mapper.Inserts
	}

	mapper, _ := cm.getColumns(entity)
	entityMapper.mu.Lock()
	defer entityMapper.mu.Unlock()

	entityMapper.dictionary[structName] = mapper
	return mapper.Inserts
}

func (cm *colMapper) GetUpdateColumns(entity Entity) []string {
	structName := structName(entity)
	if mapper, ok := entityMapper.dictionary[structName]; ok {
		return mapper.Updates
	}

	mapper, _ := cm.getColumns(entity)

	entityMapper.mu.Lock()
	defer entityMapper.mu.Unlock()

	entityMapper.dictionary[structName] = mapper
	return mapper.Updates
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

		if value := tag.Get("column"); len(value) > 0 {
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
