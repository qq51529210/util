package util

import (
	"fmt"
	"reflect"
)

// CopyStruct 拷贝 src 和 dst 中的相同名称和类型的字段，主要用于数据库结构与其他结构赋值。
// 只复制一层，不支持引用类型和结构类型的字段。小心指针是浅拷贝。
// all 表示是否覆盖 dst 有值的字段。
func CopyStruct(dst, src any, all bool) {
	// dst
	dstValue := reflect.ValueOf(dst)
	dstKind := dstValue.Kind()
	if dstKind != reflect.Pointer {
		panic("dst must be pointer")
	}
	dstValue = dstValue.Elem()
	dstKind = dstValue.Kind()
	if dstKind != reflect.Struct {
		panic("dst must be struct pointer")
	}
	// src
	srcValue := reflect.ValueOf(src)
	srcKind := srcValue.Kind()
	if srcKind == reflect.Pointer {
		srcValue = srcValue.Elem()
		srcKind = srcValue.Kind()
	}
	if srcKind != reflect.Struct {
		panic("src must be struct pointer")
	}
	// 循环 src 字段
	srcStruct := srcValue.Type()
	for i := 0; i < srcStruct.NumField(); i++ {
		srcFieldValue := srcValue.Field(i)
		// src 无效
		if !srcFieldValue.IsValid() {
			continue
		}
		srcField := srcStruct.Field(i)
		// src 同名 dst 字段
		dstFieldValue := dstValue.FieldByName(srcField.Name)
		// dst 无效，不可设置
		if !dstFieldValue.IsValid() || !dstFieldValue.CanSet() {
			continue
		}
		// 值类型
		srcFieldValueKind, dstFieldValueKind := srcFieldValue.Kind(), dstFieldValue.Kind()
		// 数据类型
		var srcFieldKind, dstFieldKind reflect.Kind
		// src 是指针类型
		if srcFieldValueKind == reflect.Pointer {
			// 空指针
			if srcFieldValue.IsNil() {
				continue
			}
			srcFieldKind = srcFieldValue.Type().Elem().Kind()
			if !isSupportType(srcFieldKind) {
				continue
			}
			if dstFieldValueKind == reflect.Pointer {
				if !all && !dstFieldValue.IsNil() {
					continue
				}
				dstFieldKind = dstFieldValue.Type().Elem().Kind()
				if !isSupportType(dstFieldKind) || srcFieldKind != dstFieldKind {
					continue
				}
				// 指针->指针
				dstFieldValue.Set(srcFieldValue)
			} else {
				if !all && !dstFieldValue.IsZero() {
					continue
				}
				dstFieldKind = dstFieldValue.Type().Kind()
				if !isSupportType(dstFieldKind) || srcFieldKind != dstFieldKind {
					continue
				}
				// 指针->值
				dstFieldValue.Set(srcFieldValue.Elem())
			}
			continue
		}
		// src 是值类型
		srcFieldKind = srcFieldValue.Type().Kind()
		if !isSupportType(srcFieldKind) {
			continue
		}
		if dstFieldValueKind == reflect.Pointer {
			if !all && !dstFieldValue.IsNil() {
				continue
			}
			dstFieldKind = dstFieldValue.Type().Elem().Kind()
			if !isSupportType(dstFieldKind) || srcFieldKind != dstFieldKind {
				continue
			}
			// 值->指针，拿地址
			dstFieldValue.Set(srcFieldValue.Addr())
		} else {
			if !all && !dstFieldValue.IsZero() {
				continue
			}
			dstFieldKind = dstFieldValue.Type().Kind()
			if !isSupportType(dstFieldKind) || srcFieldKind != dstFieldKind {
				continue
			}
			// 值->值
			dstFieldValue.Set(srcFieldValue)
		}
	}
}

// isSupportType 返回 CopyStruct 支持的字段数据类型
func isSupportType(k reflect.Kind) bool {
	return (k >= reflect.Bool && k <= reflect.Uint64) ||
		k == reflect.Float32 || k == reflect.Float64 ||
		k == reflect.Pointer || k == reflect.String
}

// IsNilOrEmpty 如果 v 是空指针，或者空值，返回 true
// 指针的值是空值，不算空值，也返回 true
func IsNilOrEmpty(v reflect.Value) bool {
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

func isStructNilOrEmpty(v reflect.Value) bool {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fmt.Println(t.Field(i).Name)
		if !IsNilOrEmpty(v.Field(i)) {
			return false
		}
	}
	return true
}
