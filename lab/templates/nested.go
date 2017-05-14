package main

import (
	"log"
	"os"
	"text/template"
)

type A struct {
	B
}
type B struct {
	C string
}

type Dynamic struct {
	Repo, Tag string
}

// text/template forces rendering with only one struct as schema
type Values struct {
	Dynamic
	A
}

func main() {
	// template location
	const rtTemplate = "./nested_values.tpl"

	// Prepare some data to insert into the template.
	var values = Values{
		Dynamic: Dynamic{
			Repo: "gcr.io/k8s-158622", Tag: "values",
		},
		A: A{B: B{C: "abc"}},
	}

	// Create a new template and parse the letter into it.
	var t *template.Template
	var err error
	if t, err = template.ParseFiles(rtTemplate); err != nil {
		log.Println("error:", err)
		return
	}

	// create output file target  TODO: create a tmp file
	log.Println("write to file", rtTemplate)
	f, err := os.Create("./runtime_values.yaml")
	if err != nil {
		log.Println("error ", err)
		return
	}
	defer f.Close()

	// write rendered template to file
	err = t.Execute(f, values)

	if err != nil {
		log.Println("error:", err)
	}

}
