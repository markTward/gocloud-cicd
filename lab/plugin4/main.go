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
	// ListProperties()
	Copy(Registry)
}

type GCR struct {
	Name        string
	Description string
	Host        string
	Project     string
	Repo        string
	Url         string
	KeyFile     string
}

// func (r *GCR) Copy(src Registry) {
//
// 	// r = reflect.New(reflect.TypeOf(GCR{})).Elem().Interface().(GCR)
// 	// r.Name = src[name]
//
// 	s := reflect.ValueOf(r).Elem()
// 	typeOfT := s.Type()
//
// 	for i := 0; i < s.NumField(); i++ {
// 		n := typeOfT.Field(i).Name
// 		v := src[n]
// 		reflect.ValueOf(r).Elem().FieldByName(n).SetString(v)
// 		// fmt.Printf("field name: %d %v :: %v\n", i, n, v)
// 	}
//
// }

func (gcr *GCR) Copy(r Registry) {

	gcrType := reflect.ValueOf(gcr).Elem().Type()
	for i := 0; i < gcrType.NumField(); i++ {
		n := gcrType.Field(i).Name
		reflect.ValueOf(gcr).Elem().FieldByName(n).SetString(r[n])
		// fmt.Printf("GCR1 field name: %d: %#v == %v\n", i, n, v)
	}
}

var typeRegistry = make(map[string]reflect.Type)

func init() {
	typeRegistry["gcr"] = reflect.TypeOf(GCR{})
	// typeRegistry["docker"] = reflect.TypeOf(DockerRegistry{})
}

var data = `
registry:
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
		newReg := makeInstance("gcr").(GCR)
		activeRegistry = &newReg
	// case "docker":
	//   newReg := makeInstance("docker").(DockerRegistry)
	//   activeRegistry = &newReg
	default:
		fmt.Println("unknown registry")
	}
	fmt.Printf("Active Registry: %#v %T\n\n", activeRegistry, activeRegistry)

	ar := activeRegistry.(Registrator)
	fmt.Printf("AR asserted: %#v %T\n\n", ar, ar)

	ar.Copy(cfg.Registry)
	fmt.Printf("Active Registry after COPY: %#v\n\n", activeRegistry)

}

func LoadConfig(data string, cfg *Config) error {
	err := yaml.Unmarshal([]byte(data), &cfg)
	return err
}

func makeInstance(name string) interface{} {
	v := reflect.New(typeRegistry[name]).Elem()
	// Maybe fill in fields here if necessary
	return v.Interface()
}

func attributes(m interface{}) (map[string]reflect.Type, map[string]string) {
	// create an attribute data structure as a map of types keyed by a string.
	attrs := make(map[string]reflect.Type)
	attrsmap := make(map[string]string)

	typ := reflect.TypeOf(m)

	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	// Only structs are supported so return an empty result if the passed object is not a struct
	if typ.Kind() != reflect.Struct {
		fmt.Printf("%v type can't have attributes inspected\n", typ.Kind())
		return attrs, attrsmap
	}
	// loop through the struct's fields and set the map
	v := reflect.ValueOf(m).Elem()

	for i := 0; i < typ.NumField(); i++ {
		p := typ.Field(i)
		if !p.Anonymous {
			attrs[p.Name] = p.Type
			attrsmap[p.Name] = v.Field(i).Interface().(string)
		}
	}

	// fmt.Println(attrmap)
	return attrs, attrsmap
}
