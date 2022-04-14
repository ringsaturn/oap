package oap

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/philchia/agollo/v4"
)

func Decode(ptr interface{}, client agollo.Client, keyOpts map[string][]agollo.OpOption) error {
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		tag := structField.Tag
		apolloKey := tag.Get("apollo")
		if apolloKey == "" {
			// Ignore empty key
			continue
		}
		var confV string
		if opts, ok := keyOpts[apolloKey]; ok {
			confV = client.GetString(apolloKey, opts...)
		} else {
			confV = client.GetString(apolloKey)
		}
		if confV == "" {
			continue
		}

		var valueToSet reflect.Value
		switch structField.Type.Kind() {
		case reflect.String:
			valueToSet = reflect.ValueOf(confV)
		case reflect.Bool:
			filedV := false
			if strings.ToLower(confV) == "true" {
				filedV = true
			}
			valueToSet = reflect.ValueOf(filedV)
		case reflect.Int:
			var filedV int

			int64V, err := strconv.ParseInt(confV, 10, 64)
			if err != nil {
				return err
			}
			filedV = int(int64V)

			valueToSet = reflect.ValueOf(filedV)
		case reflect.Float32:
			var filedV float32

			float64V, err := strconv.ParseFloat(confV, 32)
			if err != nil {
				return err
			}
			filedV = float32(float64V)

			valueToSet = reflect.ValueOf(filedV)
		case reflect.Float64:
			var filedV float64

			float64V, err := strconv.ParseFloat(confV, 64)
			if err != nil {
				return err
			}
			filedV = float64V

			valueToSet = reflect.ValueOf(filedV)
		default:
			// Not support types
			continue
		}
		v.FieldByName(structField.Name).Set(valueToSet.Convert(structField.Type))

	}
	return nil
}
