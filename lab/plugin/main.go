package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Registry struct {
	Name       string
	Properties map[string]string
}

type Registrator interface {
	ListProperties()
	Copy(Registry)
}

type GCRRegistry struct {
	Name       string
	Properties struct {
		Description string
		Host        string
		Project     string
		Repo        string
		Url         string
		KeyFile     string
	}
}

func (r *GCRRegistry) Copy(src Registry) {

	r.Name = src.Name
	s := reflect.ValueOf(&r.Properties).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		n := typeOfT.Field(i).Name
		v := src.Properties[strings.ToLower(n)]
		reflect.ValueOf(&r.Properties).Elem().FieldByName(n).SetString(v)
	}

}

func (r *GCRRegistry) ListProperties() {
	log.Println("GCR ListProperties", r.Properties)
}

type DockerRegistry struct {
	Name       string
	Properties struct {
		Description string
		Host        string
		Account     string
		Repo        string
		Url         string
	}
}

func (r *DockerRegistry) Copy(src Registry) {

	r.Name = src.Name
	s := reflect.ValueOf(&r.Properties).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		n := typeOfT.Field(i).Name
		v := src.Properties[strings.ToLower(n)]
		reflect.ValueOf(&r.Properties).Elem().FieldByName(n).SetString(v)
	}

}

func (r *DockerRegistry) ListProperties() {
	log.Println("GCR ListProperties", r.Properties)
}

// func (r *Registry) ListProperties() map[string]string {
// 	fmt.Printf("ListProps: %#v\n", r)
// 	return r.Properties
// }

type Config struct {
	Registry
}

func LoadConfig(data string, cfg *Config) error {
	err := yaml.Unmarshal([]byte(data), &cfg)
	return err
}

var typeRegistry = make(map[string]reflect.Type)

func init() {
	typeRegistry["gcr"] = reflect.TypeOf(GCRRegistry{})
	typeRegistry["docker"] = reflect.TypeOf(DockerRegistry{})
}

func makeInstance(name string) interface{} {
	v := reflect.New(typeRegistry[name]).Elem()
	// Maybe fill in fields here if necessary
	return v.Interface()
}

func main() {
	gdata := `
registry:
  name: gcr
  properties:
    description: Google Container Registry
    host: gcr.io
    project: k8sdemo-159622
    repo: gocloud
    url: gcr.io/k8sdemo-159622/gocloud
    keyfile: ./client-secret.json
`
	// ddata := `
	// registry:
	//   name: docker
	//   properties:
	//     description: Docker Hub
	//     host: docker.io
	//     account: marktward
	//     repo: gocloud
	//     url: docker.io/marktward/gocloud
	// `
	fmt.Printf("Type Registry: %#v\n\n", typeRegistry)

	data := gdata
	cfg := Config{}
	if err := LoadConfig(data, &cfg); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("Config: %#v\n\n", cfg)

	var activeRegistry interface{}
	switch cfg.Registry.Name {
	case "gcr":
		// x := makeInstance("gcr").(GCRRegistry)
		// activeRegistry = &x
		activeRegistry = makeInstance("gcr").(GCRRegistry)
	case "docker":
		// x := makeInstance("docker").(DockerRegistry)
		// activeRegistry = &x
		activeRegistry = makeInstance("docker").(DockerRegistry)
	default:
		fmt.Println("unknown registry")
	}
	fmt.Printf("Active Registry: %#v %T\n\n", activeRegistry, activeRegistry)

	ar := activeRegistry.(Registrator)
	fmt.Printf("AR asserted: %#v %T\n\n", ar, ar)

	ar.Copy(cfg.Registry)
	fmt.Printf("Active Registry after COPY: %#v\n\n", activeRegistry)
	ar.ListProperties()
}
