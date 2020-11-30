package binding

import (
	"fmt"
	"reflect"
)

func BindWithValuesMap(ptr interface{}, m map[string][]string, tag string) (bool, error) {
	return mapping(reflect.ValueOf(ptr), reflect.StructField{}, valuesMap(m), tag)
}

type valuesMap map[string][]string

func (m valuesMap) TrySet(value reflect.Value, structField reflect.StructField, name string, options options) (bool, error) {
	items, ok := m[name]
	if !ok && !options.hasDefault {
		return false, nil
	}
	switch value.Kind() {
	case reflect.Slice:
		if !ok {
			items = []string{options.defaultValue}
		}
		return true, setSlice(value, items, structField)
	case reflect.Array:
		if !ok {
			items = []string{options.defaultValue}
		}
		if len(items) != value.Len() {
			return false, fmt.Errorf("binding: %q is not valid value for %s", items, value.Type().String())
		}
		return true, setArray(value, items, structField)
	default:
		var item string
		if !ok {
			item = options.defaultValue
		}
		if len(items) > 0 {
			item = items[0]
		}
		return true, setBaseField(value, item, structField)
	}
}
