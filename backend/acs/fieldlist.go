package acs

import (
	"reflect"
	"strings"
	"sync"
)

var cache = newFieldListCache()

func FieldList(doc any) string {
	ty := reflect.TypeOf(doc)
	if ty.Kind() != reflect.Struct {
		return ""
	}

	fieldList, ok := cache.Get(ty)
	if !ok {
		fields := make([]string, 0, ty.NumField())
		for i := 0; i < ty.NumField(); i++ {
			field := ty.Field(i)
			// if field.PkgPath == "" {
			// 	continue
			// }

			var fieldName string
			if tag, ok := field.Tag.Lookup("solr"); ok {
				fieldName = tag
			} else {
				fieldName = field.Name
			}
			fields = append(fields, fieldName)
		}

		fieldList = strings.Join(fields, ",")
		cache.Set(ty, fieldList)
	}

	return fieldList

}

// type FieldLister struct {
// 	cache *cache
// }

// func NewFieldLister() *FieldLister {
// 	return &FieldLister{
// 		cache: newCache(),
// 	}
// }

// func (f *FieldLister) FieldList(doc any) string {
// 	ty := reflect.TypeOf(doc)
// 	if ty.Kind() != reflect.Struct {
// 		return ""
// 	}

// 	fieldList, ok := f.cache.Get(ty)
// 	if !ok {
// 		fields := make([]string, 0, ty.NumField())
// 		for i := 0; i < ty.NumField(); i++ {
// 			field := ty.Field(i)
// 			// if field.PkgPath == "" {
// 			// 	continue
// 			// }

// 			var fieldName string
// 			if tag, ok := field.Tag.Lookup("solr"); ok {
// 				fieldName = tag
// 			} else {
// 				fieldName = field.Name
// 			}
// 			fields = append(fields, fieldName)
// 		}

// 		fieldList = strings.Join(fields, ",")
// 		f.cache.Set(ty, fieldList)
// 	}

// 	return fieldList
// }

type fieldListCache struct {
	vault map[reflect.Type]string
	mu    sync.RWMutex
}

func newFieldListCache() fieldListCache {
	return fieldListCache{
		vault: make(map[reflect.Type]string),
	}
}

func (c *fieldListCache) Get(key reflect.Type) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.vault[key]
	return value, ok
}

func (c *fieldListCache) Set(key reflect.Type, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.vault[key] = value
}