package oap_test

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/ringsaturn/oap"
)

type DemoConfig struct {
	Foo       string `apollo:"foo"`
	Hello     string `apollo:"hello"`
	Substruct struct {
		X string `apollo:"x"`
	}
}

func TestDo(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := NewMockClient(ctrl)

	client.EXPECT().GetString(gomock.Eq("foo")).Return("bar")
	client.EXPECT().GetString(gomock.Eq("hello")).Return("hello")

	conf := &DemoConfig{}
	if err := oap.Decode(conf, client); err != nil {
		panic(err)
	}
	if conf.Foo != "bar" {
		t.Fatalf("Foo should be bar but got %v", conf.Foo)
	}

	if conf.Hello != "hello" {
		t.Fatalf("Hello should be hello but got %v", conf.Hello)
	}
}
