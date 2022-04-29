package oap_test

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	agollo "github.com/philchia/agollo/v4"
	"github.com/ringsaturn/oap"
)

type DemoConfig struct {
	Foo          string  `apollo:"foo"`
	Hello        string  `apollo:"hello"`
	Float32Field float32 `apollo:"float32Field"`
	Float64Field float32 `apollo:"float64Field"`
	BoolField    bool    `apollo:"boolField"`
	Substruct    struct {
		X string `json:"x"`
		Y int    `json:"y"`
	} `apollo:"substruct,json"`
	SubstructFromYAML struct {
		X string `yaml:"x"`
		Y int    `yaml:"y"`
	} `apollo:"substructFromYAML,yaml"`
}

var testJSONText string = `{"x": "123", "y": 0}`
var yamlText string = `
x: "fffff"
y: 12313212
`

func TestDo(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := NewMockClient(ctrl)

	client.EXPECT().GetString(gomock.Eq("foo")).Return("bar").MaxTimes(1)
	client.EXPECT().GetString(gomock.Eq("hello")).Return("hello").MaxTimes(1)
	client.EXPECT().GetString(gomock.Eq("float32Field")).Return("3.14").MaxTimes(1)
	client.EXPECT().GetString(gomock.Eq("float64Field")).Return("3.14159265").MaxTimes(1)
	client.EXPECT().GetString(gomock.Eq("boolField")).Return("true").MaxTimes(1)
	client.EXPECT().GetString(gomock.Eq("substruct")).Return(testJSONText).MaxTimes(1)
	client.EXPECT().GetString(gomock.Eq("substructFromYAML")).Return(yamlText).MaxTimes(1)

	conf := &DemoConfig{}
	if err := oap.Decode(conf, client, make(map[string][]agollo.OpOption)); err != nil {
		panic(err)
	}
	if conf.Foo != "bar" {
		t.Fatalf("Foo should be bar but got %v", conf.Foo)
	}
	if conf.Hello != "hello" {
		t.Fatalf("Hello should be hello but got %v", conf.Hello)
	}
	if conf.Float32Field != 3.14 {
		t.Fatalf("Float32Field should be `3.14` but got %v", conf.Float32Field)
	}

	if conf.Float64Field != 3.14159265 {
		t.Fatalf("Float64Field should be `3.14159265` but got %v", conf.Float64Field)
	}

	if !conf.BoolField {
		t.Fatalf("BoolField should be `true` but got %v", conf.BoolField)
	}

	if &conf.Substruct == nil {
		t.Fatal("nil")
	}
	if conf.Substruct.X != "123" {
		t.Fatalf("Substruct.X should be `123` but got %v", conf.Substruct.X)
	}

	if conf.SubstructFromYAML.X != "fffff" {
		t.Fatalf("SubstructFromYAML.X should be `'fffff'` but got %v", conf.SubstructFromYAML.X)
	}

	if conf.SubstructFromYAML.Y != 12313212 {
		t.Fatalf("SubstructFromYAML.Y should be `'12313212'` but got %v", conf.SubstructFromYAML.Y)
	}

}
