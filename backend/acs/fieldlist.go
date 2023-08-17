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

			var fieldName string
			if tag, ok := field.Tag.Lookup("json"); ok {
				if tag == "-" {
					continue
				}
				f, _, _ := strings.Cut(tag, ",")
				fieldName = f
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
