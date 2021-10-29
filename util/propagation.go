package util

import (
	"context"
	"reflect"
	"unsafe"

	"github.com/gin-gonic/gin"
)

func CopyContextValue(dest context.Context, src context.Context) context.Context {
	res := dest
	keys := GetContextKeys(src)
	for _, key := range keys {
		res = context.WithValue(dest, key, src.Value(key))
	}
	return res
}

func CopyGinContextValue(dest context.Context, src context.Context) context.Context {
	var newCtx context.Context

	// ginContextType := reflect.TypeOf((*gin.Context)(nil)).Elem()

	if srcGin, ok := src.(*gin.Context); ok {
		// srcCtx, _ := reflect.ValueOf(src).Interface().(*gin.Context)
		newCtx = CopyContextValue(dest, srcGin.Request.Context())
	} else {
		newCtx = dest
	}
	return newCtx
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
