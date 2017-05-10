//Example From blog: http://merbist.com/2011/06/27/golang-reflection-example
package main

import (
	"fmt"
	"reflect"
)

func main() {
	// iterate through the attributes of a Data Model instance
	d := Dish{Id_Val: "111", Name: "ABC", Origin: "point of origin", Query: "please work?"}

	attr, attrmap := attributes(&d)
	fmt.Printf("attr: %#v\n", attr)
	fmt.Printf("attrmap: %#v\n", attrmap)
	for name, mtype := range attr { // Instantiate Data Model and get attributes
		fmt.Printf("Name: %s, Type: %s, Value: %s\n", name, mtype, attrmap[name])
	}
}

// Data Model
type Dish struct {
	Id_Val string
	Name   string
	Origin string
	Query  string
}

// Example of how to use Go's reflection
// Print the attributes of a Data Model
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
