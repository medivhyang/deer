package binding

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type setter interface {
	TrySet(value reflect.Value, structField reflect.StructField, name string, options settingOptions) (bool, error)
}

type settingOptions struct {
	hasDefault   bool
	defaultValue string
}

func mapping(value reflect.Value, structField reflect.StructField, setter setter, tag string) (bool, error) {
	if structField.Tag.Get(tag) == "-" {
		return false, nil
	}
	switch value.Kind() {
	case reflect.Ptr:
		var (
			isNew bool
			ptr   = value
		)
		if value.IsNil() {
			isNew = true
			ptr = reflect.New(value.Type().Elem())
		}
		isSetted, err := mapping(ptr.Elem(), structField, setter, tag)
		if err != nil {
			return false, err
		}
		if isNew && isSetted {
			value.Set(ptr)
		}
		return isSetted, nil
	case reflect.Struct:
		var (
			typo     = value.Type()
			isSetted bool
		)
		for i := 0; i < value.NumField(); i++ {
			field := typo.Field(i)
			if field.PkgPath != "" && !field.Anonymous { // unexported
				continue
			}
			ok, err := mapping(value.Field(i), typo.Field(i), setter, tag)
			if err != nil {
				return false, err
			}
			isSetted = isSetted || ok
		}
		return isSetted, nil
	default:
		if structField.Anonymous {
			break
		}
		tagValue, opts := head(structField.Tag.Get(tag), ",")
		if tagValue == "" {
			tagValue = toSnake(structField.Name)
		}
		if tagValue == "" {
			return false, nil
		}
		var (
			setOpts settingOptions
			opt     string
		)
		for len(opts) > 0 {
			opt, opts = head(opts, ",")
			if k, v := head(opt, "="); k == "default" {
				setOpts.hasDefault = true
				setOpts.defaultValue = v
			}
		}
		return setter.TrySet(value, structField, tagValue, setOpts)
	}
	return false, nil
}

func setBaseField(field reflect.Value, value string, structField reflect.StructField) error {
	switch field.Kind() {
	case reflect.Int:
		return setIntField(field, value, 10, 0)
	case reflect.Int8:
		return setIntField(field, value, 10, 8)
	case reflect.Int16:
		return setIntField(field, value, 10, 16)
	case reflect.Int32:
		return setIntField(field, value, 10, 32)
	case reflect.Int64:
		switch field.Interface().(type) {
		case time.Duration:
			return setTimeDuration(field, value)
		}
		return setIntField(field, value, 10, 64)
	case reflect.Uint:
		return setUintField(field, value, 10, 0)
	case reflect.Uint8:
		return setUintField(field, value, 10, 8)
	case reflect.Uint16:
		return setUintField(field, value, 10, 16)
	case reflect.Uint32:
		return setUintField(field, value, 10, 32)
	case reflect.Uint64:
		return setUintField(field, value, 10, 64)
	case reflect.Bool:
		return setBoolField(field, value)
	case reflect.Float32:
		return setFloatField(field, value, 32)
	case reflect.Float64:
		return setFloatField(field, value, 64)
	case reflect.String:
		field.SetString(value)
	case reflect.Struct:
		switch field.Interface().(type) {
		case time.Time:
			return setTimeField(field, value, structField)
		}
		return json.Unmarshal([]byte(value), field.Addr().Interface())
	case reflect.Map:
		return json.Unmarshal([]byte(value), field.Addr().Interface())
	default:
		return errors.New("deer: set base field: unknown type")
	}
	return nil
}

func setIntField(field reflect.Value, value string, base int, bitSize int) error {
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

func setUintField(field reflect.Value, value string, base int, bitSize int) error {
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

func setFloatField(field reflect.Value, value string, bitSize int) error {
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

func setBoolField(field reflect.Value, value string) error {
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

func setTimeField(field reflect.Value, value string, structField reflect.StructField) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		tv, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		d := time.Duration(1)
		if tf == "unixnano" {
			d = time.Second
		}

		t := time.Unix(tv/int64(d), tv%int64(d))
		field.Set(reflect.ValueOf(t))
		return nil

	}

	if value == "" {
		field.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, value, l)
	if err != nil {
		return err
	}

	field.Set(reflect.ValueOf(t))
	return nil
}

func setTimeDuration(field reflect.Value, value string) error {
	d, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(d))
	return nil
}

func setArray(field reflect.Value, values []string, structField reflect.StructField) error {
	for index, value := range values {
		err := setBaseField(field.Index(index), value, structField)
		if err != nil {
			return err
		}
	}
	return nil
}

func setSlice(field reflect.Value, values []string, structField reflect.StructField) error {
	slice := reflect.MakeSlice(field.Type(), len(values), len(values))
	err := setArray(slice, values, structField)
	if err != nil {
		return err
	}
	field.Set(slice)
	return nil
}

func head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}
