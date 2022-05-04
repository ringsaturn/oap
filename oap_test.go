package oap_test

import (
	"fmt"
	"net/url"
	"strconv"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/philchia/agollo/v4"
	"github.com/ringsaturn/oap"
	"github.com/stretchr/testify/assert"
)

func unmarshalForURL(b []byte, i interface{}) error {
	u, err := url.Parse(string(b))
	if err != nil {
		return err
	}
	urlV := i.(**url.URL)
	*urlV = &*u
	return nil
}

func unmarshalForTimeTime(b []byte, i interface{}) error {
	intV, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	t := time.Unix(intV, 0)

	timeV := i.(*time.Time)
	*timeV = t
	return nil
}

func ExampleSetUnmarshalFunc() {
	oap.SetUnmarshalFunc("url", unmarshalForURL)

	f, ok := oap.GetUnmarshalFunc("url")
	if !ok {
		panic(fmt.Errorf("not get expect func"))
	}
	urlVale := &url.URL{}
	if err := f([]byte("http://example.com"), &urlVale); err != nil {
		panic(err)
	}
	fmt.Println(urlVale.Host)
	// Output: example.com
}

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
	TimeTimeField time.Time `apollo:"TimeTimeField,time.Time"`
	// TimeDuratonField time.Duration `apollo:"TimeDuratonField,time.Duration"`
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
	client.EXPECT().GetString(gomock.Eq("TimeTimeField")).Return("1651673967").MaxTimes(1)

	oap.SetUnmarshalFunc("url", unmarshalForURL)
	oap.SetUnmarshalFunc("time.Time", unmarshalForTimeTime)

	// TODO(ringsaturn): support time.Duration.
	// client.EXPECT().GetString(gomock.Eq("TimeDuratonField")).Return("3600").MaxTimes(1).MinTimes(1)
	// oap.SetUnmarshalFunc("time.Duration", unmarshalForTimeDuration)

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

	assert.Equal(t, int64(1651673967), conf.TimeTimeField.Unix())
	// log.Println("conf.TimeDuratonField", conf.TimeDuratonField)
	// assert.Equal(t, int64(3600), conf.TimeDuratonField.Seconds())
}
