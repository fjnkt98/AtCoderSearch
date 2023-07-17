package common

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type FieldLister struct {
	cache *cache
}

func NewFieldLister() *FieldLister {
	return &FieldLister{
		cache: newCache(),
	}
}

func (f *FieldLister) FieldList(doc any) (string, error) {
	ty := reflect.TypeOf(doc)
	if ty.Kind() != reflect.Struct {
		return "", fmt.Errorf("%T is not a struct", doc)
	}

	fieldList, ok := f.cache.Get(ty)
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
		f.cache.Set(ty, fieldList)
	}

	return fieldList, nil
}

type cache struct {
	vault map[reflect.Type]string
	mu    sync.RWMutex
}

func newCache() *cache {
	return &cache{
		vault: make(map[reflect.Type]string),
	}
}

func (c *cache) Get(key reflect.Type) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.vault[key]
	return value, ok
}

func (c *cache) Set(key reflect.Type, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.vault[key] = value
}
