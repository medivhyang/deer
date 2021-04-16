package binding

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type ValuesMap map[string][]string

func set(value reflect.Value, vm ValuesMap) error {
	switch value.Kind() {
	case reflect.Ptr:
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		if err := set(value.Elem(), vm); err != nil {
			return err
		}
	case reflect.Struct:
		var (
			t = value.Type()
		)
		for i := 0; i < value.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" && !f.Anonymous { // unexported
				continue
			}
			if err := setField(value.Field(i), f, vm); err != nil {
				return err
			}
		}
	default:
		return errors.New("binding: unsupported type")
	}
	return nil
}

func setField(rv reflect.Value, field reflect.StructField, vm ValuesMap) error {
	options := parseFieldOptions(field)
	if options.Ignore {
		return nil
	}
	values := vm[options.Name]
	if len(values) == 0 && options.DefaultValue == nil {
		return nil
	}
	switch rv.Kind() {
	case reflect.Slice:
		if len(values) == 0 && options.DefaultValue != nil {
			values = []string{*options.DefaultValue}
		}
		return setSlice(rv, values, options)
	case reflect.Array:
		if len(values) == 0 && options.DefaultValue != nil {
			values = []string{*options.DefaultValue}
		}
		if len(values) != rv.Len() {
			return fmt.Errorf("binding: \"%+v\" is not valid value for %s", values, rv.Type().String())
		}
		return setArray(rv, values, options)

	default:
		if field.Anonymous {
			return nil
		}
		var value string
		if len(values) > 0 {
			value = values[0]
		} else {
			value = *options.DefaultValue
		}
		return setBase(rv, value, options)
	}
}

func setArray(field reflect.Value, values []string, options fieldOptions) error {
	for index, value := range values {
		err := setBase(field.Index(index), value, options)
		if err != nil {
			return err
		}
	}
	return nil
}

func setSlice(field reflect.Value, values []string, options fieldOptions) error {
	slice := reflect.MakeSlice(field.Type(), len(values), len(values))
	err := setArray(slice, values, options)
	if err != nil {
		return err
	}
	field.Set(slice)
	return nil
}

func setBase(field reflect.Value, value string, options fieldOptions) error {
	switch field.Kind() {
	case reflect.Int:
		return setInt(field, value, 10, 0, options)
	case reflect.Int8:
		return setInt(field, value, 10, 8, options)
	case reflect.Int16:
		return setInt(field, value, 10, 16, options)
	case reflect.Int32:
		return setInt(field, value, 10, 32, options)
	case reflect.Int64:
		switch field.Interface().(type) {
		case time.Duration:
			return setTimeDuration(field, value, options)
		}
		return setInt(field, value, 10, 64, options)
	case reflect.Uint:
		return setUint(field, value, 10, 0, options)
	case reflect.Uint8:
		return setUint(field, value, 10, 8, options)
	case reflect.Uint16:
		return setUint(field, value, 10, 16, options)
	case reflect.Uint32:
		return setUint(field, value, 10, 32, options)
	case reflect.Uint64:
		return setUint(field, value, 10, 64, options)
	case reflect.Bool:
		return setBool(field, value, options)
	case reflect.Float32:
		return setFloat(field, value, 32, options)
	case reflect.Float64:
		return setFloat(field, value, 64, options)
	case reflect.String:
		field.SetString(value)
	case reflect.Struct:
		switch field.Interface().(type) {
		case time.Time:
			return setTime(field, value, options)
		}
		return json.Unmarshal([]byte(value), field.Addr().Interface())
	case reflect.Map:
		return json.Unmarshal([]byte(value), field.Addr().Interface())
	default:
		return errors.New("deer: set base field: unknown type")
	}
	return nil
}

func setInt(field reflect.Value, value string, base int, bitSize int, options fieldOptions) error {
	if value == "" {
		value = "0"
	}
	intValue, err := strconv.ParseInt(value, base, bitSize)
	if err != nil {
		return err
	}
	field.SetInt(intValue)
	return nil
}

func setUint(field reflect.Value, value string, base int, bitSize int, options fieldOptions) error {
	if value == "" {
		value = "0"
	}
	uintValue, err := strconv.ParseUint(value, base, bitSize)
	if err != nil {
		return err
	}
	field.SetUint(uintValue)
	return nil
}

func setFloat(field reflect.Value, value string, bitSize int, options fieldOptions) error {
	if value == "" {
		value = "0"
	}
	floatValue, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return err
	}
	field.SetFloat(floatValue)
	return nil
}

func setBool(field reflect.Value, value string, options fieldOptions) error {
	if value == "" {
		value = "0"
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	field.SetBool(boolValue)
	return nil
}

func setTime(field reflect.Value, value string, options fieldOptions) error {
	format := time.RFC3339
	if options.TimeFormat != nil {
		format = *options.TimeFormat
	}
	switch format {
	case TagValueTimeUnix, TagValueTimeUnixNano:
		total, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		unit := time.Duration(1)
		if format == TagValueTimeUnixNano {
			unit = time.Second
		}
		t := time.Unix(total/int64(unit), total%int64(unit))
		field.Set(reflect.ValueOf(t))
	default:
		finalLocation := time.Local
		if options.TimeZone != nil {
			l, err := time.LoadLocation(*options.TimeZone)
			if err != nil {
				return err
			}
			finalLocation = l
		}
		t, err := time.ParseInLocation(format, value, finalLocation)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(t))
	}
	return nil
}

func setTimeDuration(field reflect.Value, value string, options fieldOptions) error {
	d, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(d))
	return nil
}
