package binding

import "net/url"

func BindWithQuery(ptr interface{}, values url.Values) error {
	_, err := BindWithValuesMap(ptr, values, "query")
	return err
}
