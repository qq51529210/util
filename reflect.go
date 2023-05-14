package util

import (
	"fmt"
	"reflect"
)

// IsNilOrEmpty 如果 v 是空指针，或者空值，返回 true
// 指针的值是空值，不算空值，也返回 true
func IsNilOrEmpty(v any) bool {
	return isNilOrEmpty(reflect.ValueOf(v))
}

// isNilOrEmpty 如果 v 是空指针，或者空值，返回 true
// 指针的值是空值，不算空值，也返回 true
func isNilOrEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return true
		}
		v = v.Elem()
		switch v.Kind() {
		case reflect.Struct:
			return isStructNilOrEmpty(v)
		default:
			return false
		}
	case reflect.Func:
		return v.IsNil()
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Struct:
		return isStructNilOrEmpty(v)
	case reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.IsZero()
	}
	return false
}

// isStructNilOrEmpty 封装 IsNilOrEmpty 中判断 struct 的代码
func isStructNilOrEmpty(v reflect.Value) bool {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fmt.Println(t.Field(i).Name)
		if !isNilOrEmpty(v.Field(i)) {
			return false
		}
	}
	return true
}

// CopyStruct 拷贝 src 和 dst 中的相同名称和类型的字段，主要用于数据库结构与其他结构赋值。
func CopyStruct(dst, src any) {
	// dst
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Pointer {
		panic("dst must be pointer")
	}
	dstElem := dstVal.Elem()
	if dstElem.Kind() != reflect.Struct {
		panic("dst must be struct pointer")
	}
	// src
	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() != reflect.Pointer {
		panic("src must be pointer")
	}
	srcElem := srcVal.Elem()
	if srcElem.Kind() != reflect.Struct {
		panic("src must be struct pointer")
	}
	copyStruct(dstElem, srcElem)
}

// copyStruct 封装 CopyStruct 代码
func copyStruct(dst, src reflect.Value) {
	// type
	srcType, dstType := src.Type(), dst.Type()
	for i := 0; i < srcType.NumField(); i++ {
		// 相同名称
		srcTypeField := srcType.Field(i)
		dstTypeField, ok := dstType.FieldByName(srcTypeField.Name)
		if !ok {
			continue
		}
		// 检查类型
		srcField := src.Field(i)
		dstField := dst.Field(i)
		srcFieldKind := srcField.Kind()
		dstFieldKind := dstField.Kind()
		if srcTypeField.Type != dstTypeField.Type {
			// 看看是不是结构体
			if srcFieldKind == reflect.Pointer {
				if srcField.IsNil() {
					continue
				}
				srcField = srcField.Elem()
				srcFieldKind = srcField.Kind()
			}
			if dstFieldKind == reflect.Pointer {
				dstField = dstField.Elem()
				dstFieldKind = dstField.Kind()
			}
			if srcFieldKind == reflect.Struct && dstFieldKind == reflect.Struct {
				if !dst.IsValid() {
					continue
				}
				copyStruct(dstField, srcField)
			}
			continue
		}
		if !dstField.CanSet() {
			continue
		}
		// 赋值
		dstField.Set(srcField)
	}
}

// CopyStructIgnore 拷贝 src 和 dst 中的相同名称和类型的字段，主要用于数据库结构与其他结构赋值。
// 忽略 dst 中不为空的字段
func CopyStructIgnore(dst, src any) {
	// dst
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Pointer {
		panic("dst must be pointer")
	}
	dstElem := dstVal.Elem()
	if dstElem.Kind() != reflect.Struct {
		panic("dst must be struct pointer")
	}
	// src
	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() != reflect.Pointer {
		panic("src must be pointer")
	}
	srcElem := srcVal.Elem()
	if srcElem.Kind() != reflect.Struct {
		panic("src must be struct pointer")
	}
	copyStructIgnore(dstElem, srcElem)
}

// copyStructIgnore 封装 CopyStructIgnore 代码
func copyStructIgnore(dst, src reflect.Value) {
	// type
	srcType, dstType := src.Type(), dst.Type()
	for i := 0; i < srcType.NumField(); i++ {
		// 相同名称
		srcTypeField := srcType.Field(i)
		dstTypeField, ok := dstType.FieldByName(srcTypeField.Name)
		if !ok {
			continue
		}
		dstField := dst.Field(i)
		if !dstField.IsValid() || dstField.IsZero() {
			continue
		}
		srcField := src.Field(i)
		srcFieldKind := srcField.Kind()
		dstFieldKind := dstField.Kind()
		// 检查类型
		if srcTypeField.Type != dstTypeField.Type {
			// 看看是不是结构体
			if srcFieldKind == reflect.Pointer {
				if srcField.IsNil() {
					continue
				}
				srcField = srcField.Elem()
				srcFieldKind = srcField.Kind()
			}
			if dstFieldKind == reflect.Pointer {
				dstField = dstField.Elem()
				dstFieldKind = dstField.Kind()
			}
			if srcFieldKind == reflect.Struct && dstFieldKind == reflect.Struct {
				if !dst.IsValid() {
					continue
				}
				copyStruct(dstField, srcField)
			}
			continue
		}
		if !dstField.CanSet() {
			continue
		}
		// 赋值
		dstField.Set(srcField)
	}
}

// StructToMap 将 v 转换为 map，v 必须是结构体
func StructToMap(v any) map[string]any {
	return structToMap(reflect.ValueOf(v))
}

// structToMap 封装 StructToMap 的代码
func structToMap(vVal reflect.Value) map[string]any {
	if vVal.Kind() == reflect.Pointer {
		vVal = vVal.Elem()
	}
	vType := vVal.Type()
	result := make(map[string]any)
	for i := 0; i < vType.NumField(); i++ {
		fieldName := vType.Field(i).Name
		field := vVal.Field(i)
		if !field.IsValid() {
			fmt.Println("field name", fieldName)
			continue
		}
		fieldKind := field.Kind()
		if fieldKind == reflect.Pointer {
			field = field.Elem()
			fieldKind = field.Kind()
			if fieldKind == reflect.Invalid {
				result[fieldName] = nil
				continue
			}
		}
		if fieldKind == reflect.Struct {
			result[fieldName] = structToMap(field)
		} else {
			result[fieldName] = field.Interface()
		}
	}
	return result
}
