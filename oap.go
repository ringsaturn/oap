package oap

import (
	"reflect"

	"github.com/philchia/agollo/v4"
)

func Decode(ptr interface{}, client agollo.Client, opts ...agollo.OpOption) error {
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		// Ignore not support type
		if structField.Type.Kind() != reflect.String {
			continue
		}
		tag := structField.Tag
		apolloKey := tag.Get("apollo")
		confV := client.GetString(apolloKey, opts...)
		v.FieldByName(structField.Name).Set(reflect.ValueOf(confV).Convert(structField.Type))
	}
	return nil
}
