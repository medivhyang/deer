package binding

import (
	"net/url"
	"reflect"
)

func Bind(vm ValuesMap, i interface{}) error {
	return set(reflect.ValueOf(i), vm)
}

func BindUrlValues(values url.Values, i interface{}) error {
	return set(reflect.ValueOf(i), ValuesMap(values))
}
