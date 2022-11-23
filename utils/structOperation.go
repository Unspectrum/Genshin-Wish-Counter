package utils

import (
	"errors"
	"github.com/fatih/structs"
	"reflect"
)

func GetJSONTag(obj interface{}) ([]string, error) {
	// Check if any obj is type struct or not
	rt := reflect.TypeOf(obj)
	if rt.Kind() != reflect.Struct {
		return make([]string, 0), errors.New("parameter is not struct")
	}
	jsonTags := make([]string, 0)

	for index := 0; index < rt.NumField(); index++ {
		field := rt.Field(index)
		jsonTags = append(jsonTags, field.Tag.Get("json"))
	}
	return jsonTags, nil
}

func ParseStructToJsonMap(obj interface{}) (map[string]string, error) {
	jsonTags, err := GetJSONTag(obj)
	if err != nil {
		return make(map[string]string), errors.New("unable to get struct json tag")
	}
	newStructMap := make(map[string]string)
	oldStructMap := structs.Map(obj)

	jsonTagIndex := 0
	for _, val := range oldStructMap {
		strVal, ok := val.(string)
		if ok {
			newStructMap[jsonTags[jsonTagIndex]] = strVal
		}
	}
	return newStructMap, nil
}
