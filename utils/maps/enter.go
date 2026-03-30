package maps

import (
	"encoding/json"
	"reflect"
)

func StructToMap(data interface{}, t string) map[string]interface{} {
	var mp = make(map[string]interface{})
	v := reflect.ValueOf(data)
	for i := 0; i < v.NumField(); i++ {
		val := v.Field(i)
		tag := v.Type().Field(i).Tag.Get(t)

		if tag == "" || tag == "-" {
			continue
		}

		if val.IsNil() {
			continue
		}

		if val.Kind() == reflect.Ptr {
			v1 := val.Elem().Interface()
			if val.Elem().Kind() == reflect.Slice {
				byteData, _ := json.Marshal(v1)
				mp[tag] = string(byteData)
			} else {
				mp[tag] = v1
			}
			continue
		}
		mp[tag] = val.Interface()
	}
	return mp
}
