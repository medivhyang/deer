package binding

import "net/url"

func BindWithPostForm(ptr interface{}, values url.Values) error {
	_, err := BindWithValuesMap(ptr, values, "post_form")
	return err
}
