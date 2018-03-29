package utils

import "reflect"

func DeleteMapField(m map[string]interface{}, keys ...string) {
	for _, key := range keys {
		delete(m, key)
	}

	for key, val := range m {
		if val == nil {
			delete(m, key)
			continue
		}

		t := reflect.TypeOf(val)

		if t.Kind() == reflect.Map {
			temp := val.(map[string]interface{})
			DeleteMapField(temp, keys...)

			if len(temp) == 0 {
				delete(m, key)
			}

			continue
		}

		if t.Kind() == reflect.Slice {
			s := val.([]interface{})
			for _, sm := range s {
				DeleteMapField(sm.(map[string]interface{}), keys...)
			}

			if len(s) == 0 {
				delete(m, key)
			}

			continue
		}

		if val == reflect.Zero(t).Interface() {
			delete(m, key)
			continue
		}
	}
}
