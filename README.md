# Decode Apollo config to strcut field

Install via:

```bash
go install github.com/ringsaturn/oap
```

Usage like:

```go
import "github.com/ringsaturn/oap"

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
		URLField *url.URL `apollo:"SubstructWithInnerKeyDef.URL,url"`  // `url` is unmarshall type name
	}
}

func main(){
	// init your apollo client here
	// ...

	// add custom unmarshall
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
}
```

Support types:

- [x] String
- [x] Int
- [x] Bool
- [x] Float32
- [x] Float64
- [x] Struct from JSON or YAML
- [x] Struct with inner key def
