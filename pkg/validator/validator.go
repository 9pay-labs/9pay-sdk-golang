package validator

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type Validator struct {
	rules map[string]func(v reflect.Value, param string) bool
}

func New() *Validator {
	v := &Validator{rules: make(map[string]func(reflect.Value, string) bool)}
	v.rules["required"] = func(val reflect.Value, _ string) bool {
		if !val.IsValid() {
			return false
		}
		return !val.IsZero()
	}

	v.rules["url"] = func(val reflect.Value, _ string) bool {
		u, err := url.ParseRequestURI(val.String())
		return err == nil && u.Scheme != "" && u.Host != ""
	}

	v.rules["min"] = func(val reflect.Value, p string) bool {
		min, _ := strconv.Atoi(p)
		if val.Kind() == reflect.String || val.Kind() == reflect.Slice {
			return val.Len() >= min
		}
		if val.Kind() == reflect.Float64 {
			return val.Float() >= float64(min)
		}
		return val.Int() >= int64(min)
	}

	v.rules["max"] = func(val reflect.Value, p string) bool {
		max, _ := strconv.Atoi(p)
		if val.Kind() == reflect.String || val.Kind() == reflect.Slice {
			return val.Len() <= max
		}
		return val.Int() <= int64(max)
	}

	v.rules["len"] = func(val reflect.Value, p string) bool {
		l, _ := strconv.Atoi(p)
		return val.Len() == l
	}

	v.rules["numeric"] = func(val reflect.Value, _ string) bool {
		_, err := strconv.ParseFloat(val.String(), 64)
		return err == nil
	}

	v.rules["is_vnd"] = func(val reflect.Value, _ string) bool {
		return val.String() == "VND"
	}

	return v
}

func (v *Validator) Validate(s interface{}) error {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("validate")

		if (field.Kind() == reflect.Ptr || field.Kind() == reflect.Struct) && !field.IsZero() {
			if err := v.Validate(field.Interface()); err != nil {
				return err
			}
		}

		if tag == "" || tag == "-" {
			continue
		}

		for _, r := range strings.Split(tag, ",") {
			name, param := r, ""
			if strings.Contains(r, "=") {
				parts := strings.Split(r, "=")
				name, param = parts[0], parts[1]
			}
			if fn, ok := v.rules[name]; ok {
				if !fn(field, param) {
					return fmt.Errorf("field '%s' failed validation: %s", fieldType.Name, r)
				}
			}
		}
	}
	return nil
}

func ToMap(s interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		ft := t.Field(i)
		tag := ft.Tag.Get("json")

		if tag == "" || tag == "-" || (strings.Contains(tag, "omitempty") && field.IsZero()) {
			continue
		}

		name := strings.Split(tag, ",")[0]

		if field.Kind() == reflect.Ptr && !field.IsNil() {
			out[name] = field.Elem().Interface()
		} else if field.Kind() == reflect.Struct && ft.Name != "BaseRequest" {
			out[name] = ToMap(field.Interface())
		} else if ft.Name != "BaseRequest" {
			out[name] = field.Interface()
		}
	}
	if baseField := v.FieldByName("BaseRequest"); baseField.IsValid() {
		if extra := baseField.FieldByName("Extra"); extra.IsValid() && !extra.IsNil() {
			for _, key := range extra.MapKeys() {
				out[key.String()] = extra.MapIndex(key).Interface()
			}
		}
	}
	return out
}
