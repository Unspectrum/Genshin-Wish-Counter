package utils

import "github.com/fatih/structs"

func StructToMapInterface(obj interface{}) map[string]interface{} {
	return structs.Map(obj)
}
