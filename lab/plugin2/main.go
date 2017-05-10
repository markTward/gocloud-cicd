package main

import (
	"fmt"
	"reflect"
)

type Registry1 struct {
	Name    string
	Host    string
	Url     string
	Keyfile string
}

type Registry2 struct {
	Name       string
	Properties map[string]string
}

type GCR1 struct {
	Name    string
	Host    string
	Url     string
	Keyfile string
}

type GCR2 struct {
	Name       string
	Properties struct {
		Host    string
		Url     string
		Keyfile string
	}
}

//this is the Registry1 of types by name
var registry = map[string]reflect.Type{}

// add a type to the registry
func registerType(t reflect.Type) {
	name := t.Name()
	registry[name] = t
}

// create a new object by name, returning it as interface{}
func newByName(name string) interface{} {

	t, found := registry[name]
	if !found {
		panic("name not found!")
	}

	return reflect.New(t).Elem().Interface()
}

// func init() {
// 	registerType(reflect.TypeOf(GCR2{}))
// 	log.Printf("Type Registry: %#v / %T\n\n", registry, registry)
// }

func main() {
	// registry1 nested with properties
	r1 := Registry1{Name: "GCR1", Host: "gcr.io", Url: "gcr.io/k8s-123/gocloud", Keyfile: "client.json"}
	fmt.Printf("Registry1: %#v\n\n", r1)

	// make map of Registry1
	_, attrmap := attributes(&r1)

	// flat init of GCR
	gcr1 := reflect.New(reflect.TypeOf(GCR1{})).Elem().Interface().(GCR1)
	gcr1Type := reflect.ValueOf(&gcr1).Elem().Type()

	for i := 0; i < gcr1Type.NumField(); i++ {
		n := gcr1Type.Field(i).Name
		reflect.ValueOf(&gcr1).Elem().FieldByName(n).SetString(attrmap[n])
		// fmt.Printf("GCR1 field name: %d: %#v == %v\n", i, n, v)
	}
	fmt.Printf("gcr1: %#v / %T\n\n", gcr1, gcr1)

	// nested struct
	r2 := Registry2{Name: "GCR2", Properties: map[string]string{"Host": "gcr.io", "Url": "gcr.io/k8s-123/gocloud", "Keyfile": "client.json"}}
	fmt.Printf("Registry2: %#v\n\n", r2)
	gcr2 := reflect.New(reflect.TypeOf(GCR2{})).Elem().Interface().(GCR2)
	gcr2.Name = r2.Name

	s := reflect.ValueOf(&gcr2.Properties).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		n := typeOfT.Field(i).Name
		v := r2.Properties[n]
		reflect.ValueOf(&gcr2.Properties).Elem().FieldByName(n).SetString(v)
		// fmt.Printf("field name: %d %v :: %v\n", i, n, v)
	}

	fmt.Printf("GCR Two: %#v (%T)\n\n", gcr2, gcr2)

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
