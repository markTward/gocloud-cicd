package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"
)

func main() {
	// Define a template.
	const yaml = `
  service:
    gocloudAPI:
      image:
        repository: {{.Repo}}
        tag: {{.Tag}}
    gocloudGrpc:
      image:
        repository: {{.Repo}}
        tag: {{.Tag}}
`
	// Prepare some data to insert into the template.
	type Values struct {
		Repo, Tag string
	}
	var values = []Values{
		{"gcr.io/k8s-158622", "ada89f7da8"},
	}

	// Create a new template and parse the letter into it.
	t := template.Must(template.New("yaml").Parse(yaml))

	// write to Stdout
	fmt.Println("stdout")
	for _, r := range values {
		err := t.Execute(os.Stdout, r)
		if err != nil {
			log.Println("executing template:", err)
		}
	}

	// write to var
	var rt bytes.Buffer
	for _, r := range values {
		err := t.Execute(&rt, r)
		if err != nil {
			log.Println("executing template:", err)
		}
	}
	fmt.Println("bytes.Buffer:", rt.String())

	// write to file
	fmt.Println("write to file")
	f, err := os.Create("./runtime_values.yaml")
	if err != nil {
		log.Println("create file: ", err)
		return
	}
	defer f.Close()
	for _, r := range values {
		err := t.Execute(f, r)
		if err != nil {
			log.Println("executing template:", err)
		}
	}

}
