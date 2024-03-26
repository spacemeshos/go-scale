package scale

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// GetMaxElementsFromValue returns the max number of elements for the specified field in
// the struct passed as the v argument based on the 'scale' tag. It returns an error if v
// is not a structure, if max is not specified for the field, the field doesn't exist or
// there's a problem parsing the tag.
func GetMaxElementsFromValue(v any, fieldName string) (uint32, error) {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return 0, errors.New("bad value type")
	}
	f, found := typ.FieldByName(fieldName)
	if !found {
		return 0, fmt.Errorf("unknown field %q in %T", fieldName, v)
	}
	maxElements, err := getMaxElements(f.Tag)
	if err != nil {
		return 0, fmt.Errorf("error getting field tag %q in %T: %w", fieldName, v, err)
	}
	if maxElements == 0 {
		return 0, fmt.Errorf("no max in the scale tag for field %q in %T", fieldName, v)
	}
	return maxElements, nil
}

// GetMaxElements is a generic version of GetMaxElementsFromValue that uses the specified
// type instead of a struct value.
func GetMaxElements[T any](fieldName string) (uint32, error) {
	var v T
	return GetMaxElementsFromValue(v, fieldName)
}

// MustGetMaxElements is the same as GetMaxElements, but returns just the max tag value
// and panics in case of an error.
func MustGetMaxElements[T any](fieldName string) uint32 {
	maxElements, err := GetMaxElements[T](fieldName)
	if err != nil {
		panic(err)
	}
	return maxElements
}

func getMaxElements(tag reflect.StructTag) (uint32, error) {
	scaleTagValue, exists := tag.Lookup("scale")
	if !exists || scaleTagValue == "" {
		return 0, nil
	}
	if scaleTagValue == "" {
		return 0, errors.New("scale tag is not defined")
	}
	pairs := strings.Split(scaleTagValue, ",")
	if len(pairs) == 0 {
		return 0, errors.New("no max value found in scale tag")
	}
	var maxElementsStr string
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		data := strings.Split(pair, "=")
		if len(data) < 2 {
			continue
		}
		if data[0] != "max" {
			continue
		}
		maxElementsStr = strings.TrimSpace(data[1])
		break
	}
	if maxElementsStr == "" {
		return 0, errors.New("no max value found in scale tag")
	}
	maxElements, err := strconv.Atoi(maxElementsStr)
	if err != nil {
		return 0, fmt.Errorf("parsing max value: %w", err)
	}
	return uint32(maxElements), nil
}
