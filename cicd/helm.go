package cicd

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

type Helm struct {
	Name      string
	Version   string
	Release   string
	Namespace string
	Chartpath string
	Values    struct {
		Template string
		Output   string
	}
}

func (h *Helm) Deploy(wf *Workflow) (err error) {

	// create helm release name
	release := viper.GetString("service") + "-" + viper.GetString("branch")

	// helm required flags
	args := []string{"--install", release, "--namespace", viper.GetString("namespace")}

	// cli flag conversion
	if IsDebug() {
		args = append(args, "--debug")
	}

	// convert cicd --dryrun arg to helm dialect
	if IsDryRun() {
		args = append(args, "--dry-run")
	}

	// write runtime helm --values <file> using when available in config  otherwise create/remove a TempFile
	outFile := wf.Provider.CD.Helm.Values.Output
	var valuesFile *os.File
	switch {
	case outFile == "":
		valuesFile, err = ioutil.TempFile("", "runtime_values.yaml.")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(valuesFile.Name())
	default:
		valuesFile, err = os.Create(outFile)
		if err != nil {
			return err
		}
		defer valuesFile.Close()
	}

	// render values file from template
	err = renderHelmValuesFile(valuesFile, viper.GetString("repo"), viper.GetString("tag"))
	if err != nil {
		return fmt.Errorf("renderHelmValuesFile(): %v", err)
	}

	// join flags and positional args
	args = append(args, "--values", valuesFile.Name())
	args = append(args, viper.GetString("chart"))

	// init command vars
	var stderr bytes.Buffer
	var cmdOut []byte

	// prepend subcommand deploy to args
	args = append([]string{"upgrade"}, args...)
	cmd := exec.Command("helm", args...)

	log.Println(viper.GetString("cmdMode"), strings.Join(cmd.Args, " "))

	// execute helm command
	cmd.Stderr = &stderr
	if cmdOut, err = cmd.Output(); err != nil {
		logCmdOutput(stderr.Bytes())
		err = fmt.Errorf("%v", stderr.String())
	} else {
		logCmdOutput(cmdOut)
	}

	return err
}

func renderHelmValuesFile(valuesFile *os.File, repo string, tag string) error {
	type Values struct {
		Repo, Tag, ServiceType string
	}

	// Prepare some data to insert into the template.
	var values = Values{Repo: repo, Tag: tag}

	// initialize the template
	var t *template.Template
	var err error
	if t, err = template.ParseFiles(viper.GetString("template")); err != nil {
		return err
	}

	// render the template
	log.Println("helm runtime values filename: ", valuesFile.Name())
	err = t.Execute(valuesFile, values)

	// verify rendered file contents
	yaml, err := ioutil.ReadFile(valuesFile.Name())
	if err != nil {
		log.Println("error read yaml:", err)
		return err
	}

	LogDebug(fmt.Sprintf("helm runtime values: \n%v", string(yaml)))

	return err
}
