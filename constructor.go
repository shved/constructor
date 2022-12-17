package constructor

import (
	"reflect"
	"strings"
)

// Some strings needed to build a select params string.
const (
	DefaultParamKey       = "filter"
	DefaultDelimiter      = ","
	DefaultFieldDelimiter = "*"
)

type Builder struct {
	options Options
}

type Options struct {
	ParamKey       string
	Delimiter      string
	FieldDelimiter string
}

// NewBuilder creates a new query parameter builder with given options.
// Empty builder options will be replaced with default ones.
// You also may use default options values explicitly for whatever reason.
//
// Usage:
// package main
//
// import "github.com/shved/constructor"
//
// type HugeResourceStruct struct {
// 	Name    string  `json:"username"`
// 	Address Address `json:"address"`
// 	// ...
// 	Field108 string `json:""`
// }
//
// b := constructor.NewBuilder(constructor.Options{
// 	ParamKey:       "select",
// 	Delimiter:      constructor.DefaultDelimiter,
// 	FieldDelimiter: constructor.DefaultFieldDelimiter,
// })
//
// queryParam := b.QueryStringFromStruct(HugeResourceStruct{})
func NewBuilder(o Options) *Builder {
	if o.ParamKey == "" {
		o.ParamKey = DefaultParamKey
	}

	if o.Delimiter == "" {
		o.Delimiter = DefaultDelimiter
	}

	if o.FieldDelimiter == "" {
		o.FieldDelimiter = DefaultFieldDelimiter
	}

	return &Builder{
		options: o,
	}
}

// ParamsFromStruct gets the response entity struct instance and returns query parameter string.
func (b *Builder) QueryStringFromStruct(respStruct interface{}) string {
	repr := structRepr(respStruct)

	var res strings.Builder
	res.WriteString(b.options.ParamKey)
	res.WriteRune('=')

	var cnt int

	for k, v := range repr {
		cnt += 1

		switch {
		case len(v) == 0:
			res.WriteString(k)
			if cnt < len(repr) {
				res.WriteString(b.options.Delimiter)
			}
		case len(v) > 0:
			for i, f := range v {
				res.WriteString(k)
				res.WriteString(b.options.FieldDelimiter)
				res.WriteString(f)
				if cnt < len(repr) {
					res.WriteString(b.options.Delimiter)
				}

				if cnt == len(repr) && i < len(v)-1 {
					res.WriteString(b.options.Delimiter)
				}
			}
		}
	}

	if res.Len() == len(b.options.ParamKey)+1 {
		return ""
	}

	return res.String()
}

// structRepr makes a special map - intermediate structure that is used to build the result query
// parameterstring. Keys in map are root structure field names and values are names of nested sturct fields if
// it is a struct or a slice of structs. If it is not struct - value is just empty string slice.
func structRepr(respStruct interface{}) map[string][]string {
	repr := map[string][]string{}

	st := reflect.TypeOf(respStruct)
	for i := 0; i < st.NumField(); i++ { // Iterate root model fields.
		field := st.Field(i)

		if field.PkgPath != "" { // Skip unexported fields.
			continue
		}

		objName := fieldName(field)
		if objName == "" {
			continue
		}

		repr[objName] = []string{}

		fKind := field.Type.Kind()
		switch fKind {
		case reflect.Struct:
			fieldType := field.Type

			repr[objName] = collectFields(fieldType)
		case reflect.Slice:
			elemType := field.Type.Elem()
			if elemType.Kind() != reflect.Struct { // Skip slice if it is not a struct slice.
				continue
			}

			repr[objName] = collectFields(elemType)
		}
	}

	return repr
}

// collectFields gets a particular struct type, iterates its fields and returns a list
// of their names representation for query string.
func collectFields(t reflect.Type) []string {
	res := []string{}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.PkgPath != "" { // Skip unexported fields.
			continue
		}

		fieldName := fieldName(f)
		if fieldName == "" {
			continue
		}

		res = append(res, fieldName)
	}
	return res
}

// fieldName gets struct field and returns it's name according to rules applied to select query string.
// It tries to take it's json tag. If json tag is explicitly set to empty string it is ignored. If there is
// no json tag, the name will be the same as name of the field itself. You may also explicitly ignore the field
// setting constructor tag to omit. The field will also be omited if it is unexported.
func fieldName(field reflect.StructField) string {
	if constructorTag, ok := field.Tag.Lookup("constructor"); ok {
		if constructorTag == "omit" {
			return ""
		}
	}

	if jsonTag, ok := field.Tag.Lookup("json"); ok {
		return jsonTag
	}

	return field.Name
}
