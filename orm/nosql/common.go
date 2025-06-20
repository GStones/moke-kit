package nosql

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// isBasicType 检查 reflect.Kind 是否为 Go 的基本类型（bool, numeric, string）。
func isBasicType(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128, // 注意：Complex 类型 json.Marshal 不直接支持，但这里我们仍将其视为基本类型，后续 Marshal 会处理（或报错）
		reflect.String:
		return true
	default:
		return false
	}
}

// marshalAnyMap 将map中的非基本类型转换为json字符串
// 用于redis HASH 存储
func marshalAnyMap(m map[string]any) (map[string]any, error) {
	res := make(map[string]any)
	for k, v := range m {
		if !isBasicType(reflect.TypeOf(v).Kind()) {
			if js, err := json.Marshal(v); err != nil {
				return nil, fmt.Errorf("failed to marshal: %w", err)
			} else {
				res[k] = js
			}
		} else {
			res[k] = v
		}
	}
	return res, nil
}

// map2Struct 将 map[string]any 转换为结构体。
// 该函数会将 map 中的值转换为对应结构体字段的类型。
func map2StructShallow(m map[string]any, obj any) error {
	if obj == nil {
		return nil
	}
	v := reflect.ValueOf(obj)
	// 循环解引用，直到获取到非指针类型
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// 确保输入是结构体类型
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("input obj is not a struct, but %T", obj)
	}
	for k1, v1 := range m {
		if v1 == nil {
			continue
		}
		field := v.FieldByName(k1)
		if !field.IsValid() || !field.CanSet() {
			continue
		}
		mv1 := reflect.ValueOf(v1)
		if mv1.Type().ConvertibleTo(field.Type()) {
			field.Set(mv1.Convert(field.Type()))
			continue
		}
		
		// 处理需要 JSON 反序列化的情况
		var jsonData []byte
		switch v := v1.(type) {
		case string:
			jsonData = []byte(v)
		case []byte:
			jsonData = v
		default:
			// 如果不是字符串或字节数组，尝试序列化再反序列化
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

// struct2MapShallow 将结构体的顶层字段转换为 map[string]any。
// 参数 obj 必须是一个结构体或者指向结构体的指针。
func struct2MapShallow(obj any) (map[string]any, error) {
	if obj == nil {
		return nil, fmt.Errorf("input obj is nil")
	}

	v := reflect.ValueOf(obj)
	// 循环解引用，直到获取到非指针类型
	for v.Kind() == reflect.Ptr {
		// 处理 nil 指针的情况
		if v.IsNil() {
			return nil, fmt.Errorf("input obj is a nil pointer")
		}
		v = v.Elem()
	}
	// 确保输入是结构体类型
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input obj is not a struct, but %T", obj)
	}
	res := make(map[string]any, v.NumField())
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		// 判断json
		jsonTag := field.Tag.Get("json")
		bsonTag := field.Tag.Get("bson")
		if jsonTag == "-" || bsonTag == "-" {
			// 跳过该字段
			continue
		}
		// 只处理导出的字段 (首字母大写)
		if field.IsExported() {
			fieldName := field.Name
			// 使用 json 标签名称（如果存在）
			if jsonTag != "" && jsonTag != "-" {
				// 处理 json 标签中的选项（如 omitempty）
				tagParts := strings.Split(jsonTag, ",")
				if tagParts[0] != "" {
					fieldName = tagParts[0]
				}
			}
			fieldValue := v.Field(i)
			res[fieldName] = fieldValue.Interface()
		}
	}
	return res, nil
}

// diffMapAny 比较两个 map 并返回差异
func diffMapAny(oldMap, newMap map[string]any) (map[string]any, error) {
	changes := make(map[string]any)

	// 检查新增和修改的字段
	for key, newValue := range newMap {
		if oldValue, exists := oldMap[key]; !exists {
			// 新增字段
			changes[key] = newValue
		} else if !reflect.DeepEqual(oldValue, newValue) {
			// 修改字段
			changes[key] = newValue
		}
	}

	// 检查删除的字段
	for key := range oldMap {
		if _, exists := newMap[key]; !exists {
			// 标记为已删除
			changes[key] = nil
		}
	}

	return changes, nil
}
