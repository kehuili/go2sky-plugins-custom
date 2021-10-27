package propagation

import (
	"context"
	"reflect"
	"unsafe"
)

func CopyContextValue(dest context.Context, src context.Context) context.Context {
	keys := GetContextKeys(src)
	for _, key := range keys {
		dest = context.WithValue(dest, key, src.Value(key))
	}
	return dest
}

func GetContextKeys(ctx context.Context) []interface{} {
	contextValues := reflect.ValueOf(ctx).Elem()
	contextKeys := reflect.TypeOf(ctx).Elem()

	var keys []interface{}
	if contextKeys.Kind() == reflect.Struct {
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

			reflectField := contextKeys.Field(i)

			if reflectField.Name == "Context" {
				innerKeys := GetContextKeys(reflectValue.Interface().(context.Context))
				keys = append(keys, innerKeys...)
			} else {
				if reflectField.Name == "key" {
					keys = append(keys, reflectValue.Interface())
				}
			}
		}
	}

	return keys
}
