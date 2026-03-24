package common

import (
	"blogx_server/global"
	"reflect"

	"gorm.io/gorm"
)

type ModelMap interface {
	GetID() uint
}
type ScanMapOptions struct {
	Where *gorm.DB
	Key   string
}

func ScanMap[T ModelMap](model T, options ScanMapOptions) map[uint]T {
	var list []T

	query := global.Db.Where(model)
	if options.Where != nil {
		query.Where(options.Where)
	}
	query.Find(&list)

	var mp = map[uint]T{}
	for _, m := range list {
		mp[m.GetID()] = m
	}
	return mp
}

func ScanMapV2[T any](model T, options ScanMapOptions) map[uint]T {
	var list []T

	query := global.Db.Where(model)
	if options.Where != nil {
		query.Where(options.Where)
	}
	key := "ID"
	if options.Key != "" {
		key = options.Key
	}

	query.Find(&list)

	var mp = map[uint]T{}
	for _, m := range list {
		v := reflect.ValueOf(m)
		idField := v.FieldByName(key)
		id, ok := idField.Interface().(uint)
		if !ok {
			continue
		}
		mp[id] = m
	}
	return mp
}
