package lang

import (
	"reflect"
	"strings"
)

// IndexOf ..
func IndexOf(in interface{}, elem interface{}) int {
	inValue := reflect.ValueOf(in)
	elemValue := reflect.ValueOf(elem)

	t := inValue.Type().Kind()

	if t == reflect.String {
		return strings.Index(inValue.String(), elemValue.String())
	}

	if t == reflect.Slice {
		for i := 0; i < inValue.Len(); i++ {
			if equal(inValue.Index(i).Interface(), elem) {
				return i
			}
		}
	}

	return -1
}

func equal(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	return reflect.DeepEqual(expected, actual)

}
