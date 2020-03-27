package util

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// StructToFlatMap converts input to flat map
func StructToFlatMap(input interface{}, prefix, sep string) map[string]interface{} {
	output := make(map[string]interface{})
	valueIntoFlatMap(reflect.ValueOf(input), output, prefix, sep)
	return output
}

func valueIntoFlatMap(val reflect.Value, output map[string]interface{}, name, sep string) {
	switch val.Kind() {
	case reflect.Ptr:
		if !val.IsNil() {
			valueIntoFlatMap(val.Elem(), output, name, sep)
		}
	case reflect.Interface:
		valueIntoFlatMap(val.Elem(), output, name, sep)
		return
	case reflect.Struct:
		typ := val.Type()
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			if f.PkgPath != "" { // ignore unexported field
				continue
			}
			n := f.Name
			jt := f.Tag.Get("json")
			if jt != "" {
				arr := strings.Split(jt, ",")
				if arr[0] != "" {
					n = arr[0]
				}
			}
			if name != "" {
				n = name + sep + n
			}
			valueIntoFlatMap(val.Field(i), output, n, sep)
		}
		return
	default:
		if val.IsValid() {
			var v interface{}
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v = float64(val.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				v = float64(val.Uint())
			case reflect.Float32, reflect.Float64:
				v = float64(val.Float())
			default:
				v = val.Interface()
			}
			output[name] = v
		}
	}
}

// FormatMap formats the string with the map ( the key name of map must be combination of a-zA-z0-9 and . and _ )
// eg:  format: "this namespace is %namespace"
// 		map: 	{"namespace": "ns"}
// 		return: "this namespace is ns"
func FormatMap(format string, m map[string]interface{}) string {
	var buf bytes.Buffer
	end := len(format)
	for i := 0; i < end; {
		lasti := i
		for i < end && format[i] != '%' {
			i++
		}
		if i > lasti {
			buf.WriteString(format[lasti:i])
		}
		i++
		lasti = i
		for i < end {
			c := format[i]
			if (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') ||
				(c >= 'a' && c <= 'z') || c == '_' || c == '.' {
				i++
			} else {
				break
			}
		}
		if i > lasti {
			key := format[lasti:i]
			buf.WriteString(fmt.Sprint(m[key]))
		}
	}
	return buf.String()
}

