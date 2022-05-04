package oap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/philchia/agollo/v4"
	"gopkg.in/yaml.v3"
)

type (
	UnmarshalFunc = func([]byte, interface{}) error
	KindHandler   = func(rawtag string, expectFieldType reflect.Type, client agollo.Client, opts ...agollo.OpOption) (*reflect.Value, error)
)

var (
	registryForUnmarshal   = make(map[string]UnmarshalFunc)
	registryForKindHandler = make(map[reflect.Kind]KindHandler)
)

// You can use custom unmarshal for strcut type filed.
// Support JSON&YAML by default.
// You can add any custom type or overwrite oap's built-in method for JSON or YAML.
func SetUnmarshalFunc(name string, f UnmarshalFunc) {
	registryForUnmarshal[name] = f
}

// You can add your custom process func to relect.Kind type, or support oap not support type.
// For example, support Uint/Uint8/Uint16/Uint32/Uint64 ...
func SetKindHanlderFunc(kind reflect.Kind, f KindHandler) {
	registryForKindHandler[kind] = f
}

// You can access oap's internal registry for unmarshall.
func GetUnmarshalFunc(name string) (UnmarshalFunc, bool) {
	f, ok := registryForUnmarshal[name]
	return f, ok
}

func GetKindHanlderFunc(kind reflect.Kind) (KindHandler, bool) {
	f, ok := registryForKindHandler[kind]
	return f, ok
}

func init() {
	registryForUnmarshal["json"] = json.Unmarshal
	registryForUnmarshal["yaml"] = yaml.Unmarshal

	registryForKindHandler[reflect.Bool] = boolHandler
	registryForKindHandler[reflect.String] = stringHandler
	registryForKindHandler[reflect.Int] = intHandler
	registryForKindHandler[reflect.Float32] = float32Handler
	registryForKindHandler[reflect.Float64] = float64Handler
	registryForKindHandler[reflect.Struct] = structHandler
	registryForKindHandler[reflect.Ptr] = structHandler
}

func boolHandler(rawtag string, expectFieldType reflect.Type, client agollo.Client, opts ...agollo.OpOption) (*reflect.Value, error) {
	confV := client.GetString(rawtag, opts...)
	filedV := false
	if strings.ToLower(confV) == "true" {
		filedV = true
	}
	valueToSet := reflect.ValueOf(filedV)
	return &valueToSet, nil
}

func stringHandler(rawtag string, expectFieldType reflect.Type, client agollo.Client, opts ...agollo.OpOption) (*reflect.Value, error) {
	confV := client.GetString(rawtag, opts...)
	valueToSet := reflect.ValueOf(confV)
	return &valueToSet, nil
}

func float32Handler(rawtag string, expectFieldType reflect.Type, client agollo.Client, opts ...agollo.OpOption) (*reflect.Value, error) {
	confV := client.GetString(rawtag, opts...)

	var filedV float32

	float64V, err := strconv.ParseFloat(confV, 32)
	if err != nil {
		return nil, err
	}
	filedV = float32(float64V)

	valueToSet := reflect.ValueOf(filedV)
	return &valueToSet, nil
}

func float64Handler(rawtag string, expectFieldType reflect.Type, client agollo.Client, opts ...agollo.OpOption) (*reflect.Value, error) {
	confV := client.GetString(rawtag, opts...)

	var filedV float64

	float64V, err := strconv.ParseFloat(confV, 32)
	if err != nil {
		return nil, err
	}
	filedV = float64V

	valueToSet := reflect.ValueOf(filedV)
	return &valueToSet, nil
}

func structWithMarhsallHandler(rawtag string, expectFieldType reflect.Type, client agollo.Client, opts ...agollo.OpOption) (*reflect.Value, error) {
	apolloKeyParts := strings.Split(rawtag, ",")
	apolloKey := apolloKeyParts[0]

	confV := client.GetString(apolloKey, opts...)

	var unmarshalType string
	if len(apolloKeyParts) == 2 {
		unmarshalType = apolloKeyParts[1]
	}
	unmarshalFunc, ok := GetUnmarshalFunc(unmarshalType)
	if !ok {
		return nil, fmt.Errorf("unmarshalType=`%v` from rawtag=`%v` not suported yet", unmarshalType, rawtag)
	}
	v := reflect.New(expectFieldType)
	newP := v.Interface()
	if err := unmarshalFunc([]byte(confV), newP); err != nil {
		return nil, err
	}
	valueToSet := reflect.Indirect(reflect.ValueOf(newP))
	return &valueToSet, nil
}

func structWithEmptyKeyHandler(rawtag string, expectFieldType reflect.Type, client agollo.Client, opts ...agollo.OpOption) (*reflect.Value, error) {
	v := reflect.New(expectFieldType)
	newP := v.Interface()
	if err := Decode(newP, client, make(map[string][]agollo.OpOption)); err != nil {
		return nil, err
	}
	valueToSet := reflect.Indirect(reflect.ValueOf(newP))
	return &valueToSet, nil
}

func structHandler(apolloKey string, expectFieldType reflect.Type, client agollo.Client, opts ...agollo.OpOption) (*reflect.Value, error) {
	if apolloKey == "" {
		return structWithEmptyKeyHandler(apolloKey, expectFieldType, client, opts...)
	}
	return structWithMarhsallHandler(apolloKey, expectFieldType, client, opts...)
}

func intHandler(apolloKey string, expectFieldType reflect.Type, client agollo.Client, opts ...agollo.OpOption) (*reflect.Value, error) {
	confV := client.GetString(apolloKey, opts...)

	var filedV int

	int64V, err := strconv.ParseInt(confV, 10, 64)
	if err != nil {
		return nil, err
	}
	filedV = int(int64V)

	valueToSet := reflect.ValueOf(filedV)
	return &valueToSet, nil
}

func Decode(ptr interface{}, client agollo.Client, keyOpts map[string][]agollo.OpOption) error {
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		tag := structField.Tag
		apolloRawKey := tag.Get("apollo")
		apolloKeyParts := strings.Split(apolloRawKey, ",")
		apolloKey := apolloKeyParts[0]

		// Check struct field type if supported yet
		filedTypeKind := structField.Type.Kind()
		handler, ok := registryForKindHandler[filedTypeKind]
		if !ok {
			continue
		}

		var valueToSetPtr *reflect.Value
		var valueToSetErr error

		// use opts if provieded
		if opts, ok := keyOpts[apolloKey]; ok {
			valueToSetPtr, valueToSetErr = handler(apolloRawKey, structField.Type, client, opts...)
		} else {
			valueToSetPtr, valueToSetErr = handler(apolloRawKey, structField.Type, client)
		}
		if valueToSetErr != nil {
			return valueToSetErr
		}
		// get value from ptr
		valueToSet := valueToSetPtr
		v.FieldByName(structField.Name).Set(valueToSet.Convert(structField.Type))
	}
	return nil
}
