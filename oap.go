package oap

import (
	"reflect"
)

type MinimalClient interface {
	GetString(string) string
}

func Decode(ptr interface{}, client MinimalClient) error {

	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		// Ignore not support type
		if structField.Type.Kind() != reflect.String {
			continue
		}
		tag := structField.Tag
		apolloKey := tag.Get("apollo")
		v.FieldByName(structField.Name).Set(reflect.ValueOf(client.GetString(apolloKey)).Convert(structField.Type))
	}
	return nil
}
