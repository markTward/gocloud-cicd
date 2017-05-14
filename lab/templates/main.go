package main

import (
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
	type Value struct {
		Repo string
		Tag   string
	}

  var repo = "gcr.io/k8s-158622/gocloud"
  var tag = "198273adf"
  var values = []Value{{repo, tag}}

	// Create a new template and parse the letter into it.
	t := template.Must(template.New("yaml").Parse(yaml))

	// Execute the template for each recipient.
	for _, v := range values {
		err := t.Execute(os.Stdout, v)
		if err != nil {
			log.Println("executing template:", err)
		}
	}

}
