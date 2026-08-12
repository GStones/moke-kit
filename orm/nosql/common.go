package nosql

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

const (
	// cacheFieldVersion stores the document CAS version inside a HASH cache entry.
	cacheFieldVersion = "__version"
	// cacheFieldEpoch stores the document generation fence inside a HASH cache entry.
	cacheFieldEpoch = "__epoch"
	// cacheFieldData stores the JSON-encoded document payload inside a HASH cache entry.
	cacheFieldData = "__data"
)

// isBasicType checks whether reflect.Kind is a Go basic type.
func isBasicType(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return true
	default:
		return false
	}
}

// marshalAnyMap converts non-basic map values to JSON for Redis HASH storage.
func marshalAnyMap(m map[string]any) (map[string]any, error) {
	if m == nil {
		return make(map[string]any), nil
	}

	res := make(map[string]any, len(m))
	for k, v := range m {
		if v == nil {
			res[k] = nil
			continue
		}
		vType := reflect.TypeOf(v)
		if isBasicType(vType.Kind()) {
			res[k] = v
			continue
		}
		js, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal: %w", err)
		}
		res[k] = js
	}
	return res, nil
}

func structFieldIndex(v reflect.Value) map[string]reflect.Value {
	t := v.Type()
	index := make(map[string]reflect.Value, v.NumField()*2)
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		fv := v.Field(i)
		if !fv.CanSet() {
			continue
		}
		index[field.Name] = fv
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		name := strings.Split(jsonTag, ",")[0]
		if name != "" {
			index[name] = fv
		}
	}
	return index
}

// map2StructShallow copies map values into exported struct fields.
// Keys may be Go field names or json tag names produced by struct2MapShallow.
func map2StructShallow(m map[string]any, obj any) error {
	if obj == nil {
		return nil
	}
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return fmt.Errorf("input obj is a nil pointer")
		}
		if v.Elem().Kind() == reflect.Pointer && v.Elem().IsNil() {
			v.Elem().Set(reflect.New(v.Elem().Type().Elem()))
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("input obj is not a struct, but %T", obj)
	}

	fields := structFieldIndex(v)
	for k1, v1 := range m {
		if v1 == nil {
			continue
		}
		field, ok := fields[k1]
		if !ok {
			continue
		}
		mv1 := reflect.ValueOf(v1)
		if mv1.Type().ConvertibleTo(field.Type()) {
			field.Set(mv1.Convert(field.Type()))
			continue
		}

		var jsonData []byte
		switch typed := v1.(type) {
		case string:
			jsonData = []byte(typed)
		case []byte:
			jsonData = typed
		default:
			var err error
			jsonData, err = json.Marshal(v1)
			if err != nil {
				return fmt.Errorf("failed to marshal value for field %s: %w", k1, err)
			}
		}
		if err := json.Unmarshal(jsonData, field.Addr().Interface()); err != nil {
			return fmt.Errorf("failed to unmarshal field %s: %w", k1, err)
		}
	}
	return nil
}

// struct2MapShallow converts exported top-level struct fields to map[string]any.
func struct2MapShallow(obj any) (map[string]any, error) {
	if obj == nil {
		return nil, fmt.Errorf("input obj is nil")
	}

	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, fmt.Errorf("input obj is a nil pointer")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input obj is not a struct, but %T", obj)
	}

	res := make(map[string]any, v.NumField())
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		bsonTag := field.Tag.Get("bson")
		if jsonTag == "-" || bsonTag == "-" {
			continue
		}
		if !field.IsExported() {
			continue
		}
		fieldName := field.Name
		if jsonTag != "" {
			tagParts := strings.Split(jsonTag, ",")
			if tagParts[0] != "" {
				fieldName = tagParts[0]
			}
		}
		res[fieldName] = v.Field(i).Interface()
	}
	return res, nil
}

// diffMapAny compares two maps and returns added/changed/deleted keys.
func diffMapAny(oldMap, newMap map[string]any) (map[string]any, error) {
	changes := make(map[string]any)
	for key, newValue := range newMap {
		if oldValue, exists := oldMap[key]; !exists {
			changes[key] = newValue
		} else if !reflect.DeepEqual(oldValue, newValue) {
			changes[key] = newValue
		}
	}
	for key := range oldMap {
		if _, exists := newMap[key]; !exists {
			changes[key] = nil
		}
	}
	return changes, nil
}

func parseCacheVersion(v any) (noptions.Version, bool) {
	switch typed := v.(type) {
	case int64:
		return typed, true
	case int32:
		return noptions.Version(typed), true
	case int:
		return noptions.Version(typed), true
	case float64:
		return noptions.Version(typed), true
	case json.Number:
		n, err := typed.Int64()
		if err != nil {
			return 0, false
		}
		return n, true
	case string:
		n, err := strconv.ParseInt(typed, 10, 64)
		if err != nil {
			return 0, false
		}
		return n, true
	case []byte:
		n, err := strconv.ParseInt(string(typed), 10, 64)
		if err != nil {
			return 0, false
		}
		return n, true
	default:
		return 0, false
	}
}

func asJSONBytes(v any) ([]byte, error) {
	switch typed := v.(type) {
	case nil:
		return nil, fmt.Errorf("nil json payload")
	case []byte:
		return typed, nil
	case string:
		return []byte(typed), nil
	case json.RawMessage:
		return typed, nil
	default:
		return json.Marshal(typed)
	}
}
