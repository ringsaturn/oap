package oap_test

import (
	"testing"

	"github.com/ringsaturn/oap"
)

type FakeClient struct {
}

func (c *FakeClient) GetValue(key string) string {
	if key == "foo" {
		return "bar"
	}
	return key
}

type DemoConfig struct {
	Foo   string `apollo:"foo"`
	Hello string `apollo:"hello"`
}

func TestDo(t *testing.T) {
	client := &FakeClient{}
	conf := &DemoConfig{}
	if err := oap.Do(conf, client); err != nil {
		panic(err)
	}
	if conf.Foo != "bar" {
		t.Fatalf("Foo should be bar but got %v", conf.Foo)
	}

	if conf.Hello != "hello" {
		t.Fatalf("Hello should be hello but got %v", conf.Hello)
	}
}
