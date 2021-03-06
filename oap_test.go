package oap_test

import (
	"net/url"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/philchia/agollo/v4"
	"github.com/ringsaturn/oap"
	"github.com/stretchr/testify/assert"
)

type DemoConfig struct {
	Foo          string  `apollo:"foo"`
	Hello        string  `apollo:"hello"`
	Float32Field float32 `apollo:"float32Field"`
	Float64Field float64 `apollo:"float64Field"`
	BoolField    bool    `apollo:"boolField"`
	Substruct    struct {
		X string `json:"x"`
		Y int    `json:"y"`
	} `apollo:"substruct,json"`
	SubstructFromYAML struct {
		X string `yaml:"x"`
		Y int    `yaml:"y"`
	} `apollo:"substructFromYAML,yaml"`
	SubstructWithInnerKeyDef struct {
		X        string   `apollo:"SubstructWithInnerKeyDef.X"`
		Y        string   `apollo:"SubstructWithInnerKeyDef.Y"`
		URLField *url.URL `apollo:"SubstructWithInnerKeyDef.URL,url"`
	}
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
	client.EXPECT().GetString(gomock.Eq("SubstructWithInnerKeyDef.X")).Return("balabala").MaxTimes(1)
	client.EXPECT().GetString(gomock.Eq("SubstructWithInnerKeyDef.Y")).Return("habahaba").MaxTimes(1)
	client.EXPECT().GetString(gomock.Eq("SubstructWithInnerKeyDef.URL")).Return("http://example.com").MaxTimes(1)

	oap.SetUnmarshalFunc("url", func(b []byte, i interface{}) error {
		u, err := url.Parse(string(b))
		if err != nil {
			return err
		}
		urlV := i.(**url.URL)
		*urlV = &*u
		return nil
	})

	conf := &DemoConfig{}
	if err := oap.Decode(conf, client, make(map[string][]agollo.OpOption)); err != nil {
		panic(err)
	}
	assert.Equal(t, "bar", conf.Foo)
	assert.Equal(t, "hello", conf.Hello)
	assert.Equal(t, float32(3.14), conf.Float32Field)
	assert.InDelta(t, float64(3.14159265), conf.Float64Field, 0.0000001)
	assert.Equal(t, true, conf.BoolField)

	assert.Equal(t, "123", conf.Substruct.X)

	assert.Equal(t, "fffff", conf.SubstructFromYAML.X)
	assert.Equal(t, 12313212, conf.SubstructFromYAML.Y)

	assert.Equal(t, "balabala", conf.SubstructWithInnerKeyDef.X)
	assert.Equal(t, "habahaba", conf.SubstructWithInnerKeyDef.Y)
	assert.Equal(t, "example.com", conf.SubstructWithInnerKeyDef.URLField.Host)
}
