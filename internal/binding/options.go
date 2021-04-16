package binding

import (
	"reflect"
	"strings"
)

const (
	TagBinding           = "binding"
	TagItemSep           = ","
	TagKVSep             = "="
	TagKeyIgnore         = "-"
	TagKeyName           = "name"
	TagKeyDefault        = "default"
	TagKeyTime           = "time"
	TagKeyTimeZone       = "tz"
	TagValueTimeUnix     = "unix"
	TagValueTimeUnixNano = "unix_nano"
)

type fieldOptions struct {
	Ignore       bool
	Name         string
	DefaultValue *string
	TimeFormat   *string
	TimeZone     *string
}

func parseFieldOptions(sf reflect.StructField) fieldOptions {
	m := parseTagToMap(sf.Tag.Get(TagBinding), TagItemSep, TagKVSep)
	f := fieldOptions{}
	if _, ok := m[TagKeyIgnore]; ok {
		f.Ignore = true
	}
	if name, ok := m[TagKeyName]; ok && name != "" {
		f.Name = name
	} else {
		f.Name = toSnake(sf.Name)
	}
	if defaultValue, ok := m[TagKeyDefault]; ok {
		f.DefaultValue = &defaultValue
	}
	if timeFormat, ok := m[TagKeyTime]; ok && timeFormat != "" {
		f.TimeFormat = &timeFormat
	}
	if timeZone, ok := m[TagKeyTimeZone]; ok && timeZone != "" {
		f.TimeZone = &timeZone
	}
	return f
}

func parseTagToMap(tag string, itemSep, kvSep string) map[string]string {
	r := map[string]string{}
	items := strings.Split(tag, itemSep)
	for _, item := range items {
		var (
			kv   = strings.Split(item, kvSep)
			k, v string
		)
		if len(kv) == 1 {
			k = kv[0]
		} else {
			k, v = kv[0], kv[1]
		}
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)
		if k != "" {
			r[k] = v
		}
	}
	return r
}
