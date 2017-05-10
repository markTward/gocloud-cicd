package main

import (
	"fmt"
	"reflect"

	yaml "gopkg.in/yaml.v2"
)

type Registry map[string]string

type Config struct {
	Registry
}

type Registrator interface {
	Copy(Registry)
}

type GCRegistry struct {
	Name        string
	Description string
	Host        string
	Project     string
	Repo        string
	Url         string
	Keyfile     string
}

func (gcr *GCRegistry) Copy(r Registry) {
	gcrType := reflect.ValueOf(gcr).Elem().Type()
	for i := 0; i < gcrType.NumField(); i++ {
		n := gcrType.Field(i).Name
		reflect.ValueOf(gcr).Elem().FieldByName(n).SetString(r[n])
	}
}

var providerRegistry = make(map[string]reflect.Type)

func init() {
	providerRegistry["gcr"] = reflect.TypeOf(GCRegistry{})
}

func makeInstance(name string) interface{} {
	v := reflect.New(providerRegistry[name]).Elem()
	return v.Interface()
}

func LoadConfig(data string, cfg *Config) error {
	err := yaml.Unmarshal([]byte(data), &cfg)
	return err
}

var data = `
registry:
  # element keys must match case of fields in destination struct.  ex: GCR{Name: string, ...}
  Name: gcr
  Description: Google Container Registry
  Host: gcr.io
  Project: k8sdemo-159622
  Repo: gocloud
  Url: gcr.io/k8sdemo-159622/gocloud
  Keyfile: ./client-secret.json
`

func main() {
	cfg := Config{}
	if err := LoadConfig(data, &cfg); err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("Config: %#v\n\n", cfg.Registry)

	var activeRegistry interface{}
	switch cfg.Registry["Name"] {
	case "gcr":
		newReg := makeInstance("gcr").(GCRegistry)
		activeRegistry = &newReg
	default:
		fmt.Println("unknown registry")
	}
	fmt.Printf("Active Registry: %#v\n\n", activeRegistry)

	ar := activeRegistry.(Registrator)

	ar.Copy(cfg.Registry)
	fmt.Printf("Active Registry after COPY: %#v\n\n", activeRegistry)

}
