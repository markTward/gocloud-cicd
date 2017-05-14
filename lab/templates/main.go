package main

import (
	"log"
	"os"
	"text/template"
)

func main() {

	const (
		rtTemplate = "./runtime_values.tpl"
		rtYaml     = "./runtime_values.yaml"
	)

	repo, tag := "gcr.io/k8s-158622", "a90dsf809a8d"

	err := renderHelmValuesFile(rtTemplate, rtYaml, repo, tag)
	if err != nil {
		log.Println("error: ", err)
		os.Exit(1)
	}

}

func renderHelmValuesFile(tf string, of string, repo string, tag string) error {

	// Prepare some data to insert into the template.
	type Values struct {
		Repo, Tag string
	}
	var values = Values{Repo: repo, Tag: tag}

	// initialize the template
	var t *template.Template
	var err error
	if t, err = template.ParseFiles(tf); err != nil {
		return err
	}

	// TODO: best practice to write to random temp file
	// create target file for output
	f, err := os.Create(of)
	if err != nil {
		return err
	}
	defer f.Close()

	// render the template
	err = t.Execute(f, values)

	return err
}
