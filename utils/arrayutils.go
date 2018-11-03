package utils

import "reflect"

func InArray(array interface{}, item interface{}) (int, bool) {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(item, s.Index(i).Interface()) == true {
				return i, true
			}
		}
	}
	return -1, false
}

func RemoveFromArray(array []interface{}, idx int) []interface{} {
	if idx >= len(array) {
		return array
	}

	return append(array[:idx], array[idx+1:]...)
}