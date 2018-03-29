package utils

import (
	"reflect"
	"strings"
)

func GetStructName(stc interface{}) string {
	t := reflect.TypeOf(stc)
	nameArr := strings.Split(t.String(), ".")
	return nameArr[len(nameArr)-1]
}
