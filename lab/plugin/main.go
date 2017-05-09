package main

import (
	"fmt"
	"reflect"

	yaml "gopkg.in/yaml.v2"
)

type GCRRegistry struct {
	Name        string
	Description string
	Host        string
	Project     string
	Repo        string
	Url         string
	KeyFile     string
}

func (r *GCRRegistry) ListProperties() string {
	return fmt.Sprintf("GCR Properties: %#v", r)
}

type DockerRegistry struct {
	Name        string
	Description string
	Host        string
	Account     string
	Repo        string
	Url         string
}

func (r *DockerRegistry) ListProperties() string {
	return fmt.Sprintf("Docker Properties: %#v", r.ListProperties())
}

type Registrator interface {
	ListProperties() map[string]string
}

type Registry struct {
	Name       string
	Properties map[string]string
}

func (r *Registry) ListProperties() map[string]string {
	fmt.Printf("ListProps: %#v\n", r)
	return r.Properties
}

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
	ddata := `
registry:
  name: docker
  properties:
    description: Docker Hub
    host: docker.io
    account: marktward
    repo: gocloud
    url: docker.io/marktward/gocloud
`
	fmt.Println("GRC Data:", gdata)
	fmt.Println("Docker Data:", ddata)

	data := gdata
	cfg := Config{}
	if err := LoadConfig(data, &cfg); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("Config: %#v\n\n", cfg)

	fmt.Printf("Type Registry: %#v\n\n", typeRegistry)
	var activeRegistry interface{}
	switch cfg.Registry.Name {
	case "gcr":
		activeRegistry = makeInstance("gcr")
		activeRegistry = &cfg.Registry
	case "docker":
		activeRegistry = makeInstance("docker")
		activeRegistry = &cfg.Registry
	default:
		fmt.Println("unknown registry")
	}

	fmt.Printf("Active Registry: %#v %T\n", activeRegistry, activeRegistry)
	ar := activeRegistry.(Registrator)
	// &activeRegistry.ListProperties()

	// ar := r.(Registrator.ListProperties())
	fmt.Printf("AR Registrator: %#v\n\n", ar)

	fmt.Printf("AR Properties: %#v\n", ar.ListProperties())
}
