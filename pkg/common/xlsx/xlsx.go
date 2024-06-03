package xlsx

import (
	"github.com/openimsdk/tools/errs"
	"reflect"
	"strconv"
)

// GetColumnIndex get header and index
func GetColumnIndex(headers []string) map[string]int {
	colIndex := make(map[string]int)
	for i, header := range headers {
		colIndex[header] = i
	}
	return colIndex
}

// SetStructValues reflect to struct
func SetStructValues(data interface{}, row []string, colIndex map[string]int) error {
	val := reflect.ValueOf(data).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := typ.Field(i).Tag.Get("column")
		if idx, ok := colIndex[tag]; ok && idx < len(row) {
			if err := SetValue(field, row[idx]); err != nil {
				return err
			}
		}
	}
	return nil
}

// SetValue set value
func SetValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			field.SetInt(intValue)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintValue, err := strconv.ParseUint(value, 10, 64); err == nil {
			field.SetUint(uintValue)
		}
	case reflect.Float32, reflect.Float64:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			field.SetFloat(floatValue)
		}
	default:
		return errs.New("not handle error type")
	}
	return nil
}
