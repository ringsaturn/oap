package oap

import (
	"reflect"
)

type MinimalClient interface {
	GetValue(string) string
}

func Do(ptr interface{}, client MinimalClient) error {

	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		tag := structField.Tag
		apolloKey := tag.Get("apollo")
		v.FieldByName(structField.Name).Set(reflect.ValueOf(client.GetValue(apolloKey)).Convert(structField.Type))
	}
	return nil
}
