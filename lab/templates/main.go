package main

import (
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/davecgh/go-spew/spew"
	"github.com/markTward/gocloud-cicd/config"
)

func main() {
	// initialize configuration object
	cfg := config.New()
	if err := config.Load("../../../gocloud/cicd.yaml", &cfg); err != nil {
		log.Println("error:", err)
	}

	spew.Dump(cfg)

	rtTemplate := cfg.CDProvider.Helm.Options.Values.Template
	rtYaml := cfg.CDProvider.Helm.Options.Values.Output

	repo, tag := "gcr.io/k8s-158622", "a90dsf809a8d"

	err := renderHelmValuesFile(rtTemplate, rtYaml, repo, tag)
	if err != nil {
		log.Println("error: ", err)
		os.Exit(1)
	}

}

func renderHelmValuesFile(tf string, of string, repo string, tag string) error {
	type Values struct {
		Repo, Tag, ServiceType string
	}

	// Prepare some data to insert into the template.
	var values = Values{Repo: repo, Tag: tag}
	spew.Dump(values)

	// initialize the template
	var t *template.Template
	var err error
	if t, err = template.ParseFiles(tf); err != nil {
		return err
	}

	var f *os.File
	switch {
	case of == "":
		f, err = ioutil.TempFile("", "runtime_values.yaml")
		log.Println("tmp output file:", f.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(f.Name()) // clean up
	default:
		// create target file for output
		f, err = os.Create(of)
		log.Println("cicd output file:", f.Name())
		if err != nil {
			return err
		}
		defer f.Close()
	}

	// render the template
	log.Println("output file before exec:", f.Name())
	err = t.Execute(f, values)

	// verify rendered file contents
	yaml, err := ioutil.ReadFile(f.Name())
	if err != nil {
		log.Println("error read yaml:", err)
		return err
	}

	log.Println(string(yaml))

	return err
}
