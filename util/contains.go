package util

import (
	"reflect"
	"strings"
)

// Contains 元素包含
func Contains(in interface{}, elem interface{}) bool {
	inValue := reflect.ValueOf(in)
	elemValue := reflect.ValueOf(elem)
	inType := inValue.Type()

	switch inType.Kind() {
	case reflect.String:
		return strings.Contains(inValue.String(), elemValue.String())
	case reflect.Map:
		for _, key := range inValue.MapKeys() {
			if equal(key.Interface(), elem) {
				return true
			}
		}
	case reflect.Slice:
		for i := 0; i < inValue.Len(); i++ {
			if equal(inValue.Index(i).Interface(), elem) {
				return true
			}
		}
	}

	return false
}

func equal(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	return reflect.DeepEqual(expected, actual)

}

// CompareFunc ..
type CompareFunc func(interface{}, interface{}) bool

// IndexOf ..
func IndexOf(in interface{}, e interface{}, cmp CompareFunc) int {
	var (
		i   int
		ins = reflect.ValueOf(in)
		n   = ins.Len()
	)
	for ; i < n; i++ {
		if cmp(e, ins.Index(i).Interface()) {
			return i
		}
	}
	return -1
}
