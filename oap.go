package oap

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/philchia/agollo/v4"
	"gopkg.in/yaml.v3"
)

func Decode(ptr interface{}, client agollo.Client, keyOpts map[string][]agollo.OpOption) error {
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		tag := structField.Tag
		apolloRawKey := tag.Get("apollo")
		if apolloRawKey == "" {
			// Ignore empty key
			continue
		}
		apolloKeyParts := strings.Split(apolloRawKey, ",")
		apolloKey := apolloKeyParts[0]

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
		case reflect.Struct:
			var unmarshallType string
			if len(apolloKeyParts) == 2 {
				unmarshallType = apolloKeyParts[1]
			}
			switch unmarshallType {
			case "json":
				v := reflect.New(structField.Type)
				newP := v.Interface()
				if err := json.Unmarshal([]byte(confV), newP); err != nil {
					return err
				}
				valueToSet = reflect.Indirect(reflect.ValueOf(newP))
			case "yaml":
				v := reflect.New(structField.Type)
				newP := v.Interface()
				if err := yaml.Unmarshal([]byte(confV), newP); err != nil {
					return err
				}
				valueToSet = reflect.Indirect(reflect.ValueOf(newP))
			default:
				continue
			}
		default:
			continue
		}
		v.FieldByName(structField.Name).Set(valueToSet.Convert(structField.Type))

	}
	return nil
}
