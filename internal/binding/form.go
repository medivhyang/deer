package binding

import "net/url"

func BindWithForm(ptr interface{}, values url.Values) error {
	_, err := BindWithValuesMap(ptr, values, "form")
	return err
}
