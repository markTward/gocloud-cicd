package main

import (
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
	const rtTemplate = "./runtime_values.tpl"

	// Prepare some data to insert into the template.
	type Values struct {
		Repo, Tag string
	}
	var values = []Values{
		{Repo: "gcr.io/k8s-158622", Tag: "98ads7fa"},
	}

	// Create a new template and parse the letter into it.
	// t := template.Must(template.New("yaml").Parse(yaml))
	var t *template.Template
	var err error
	if t, err = template.ParseFiles(rtTemplate); err != nil {
		log.Println("error:", err)
		return
	}

	// write to file
	fmt.Println("write to file", rtTemplate)
	f, err := os.Create("./runtime_values.yaml")
	if err != nil {
		log.Println("create file: ", err)
		return
	}
	defer f.Close()

	for _, r := range values {
		err = t.Execute(f, r)
		if err != nil {
			log.Println("executing template:", err)
		}
	}

	// write to var
	// var rt bytes.Buffer
	// for _, r := range values {
	// 	err := t.Execute(&rt, r)
	// 	if err != nil {
	// 		log.Println("executing template:", err)
	// 	}
	// }
	// fmt.Println("write to var:", rt.String())

	// write to Stdout
	// fmt.Println("stdout")
	// for _, r := range values {
	// 	err := t.Execute(os.Stdout, r)
	// 	if err != nil {
	// 		log.Println("executing template:", err)
	// 	}
	// }

}
