package maps

import "reflect"

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
			mp[tag] = val.Elem().Interface()
			continue
		}
		mp[tag] = val.Interface()
	}
	return mp
}
