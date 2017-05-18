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

	"github.com/urfave/cli"
)

type Helm struct {
	Name      string
	Version   string
	Enabled   bool
	Release   string
	Namespace string
	Chartpath string
	Options   struct {
		Flags  []string
		Values struct {
			Template string
			Output   string
		}
	}
}

func (h *Helm) Deploy(ctx *cli.Context, cfg *Config) (err error) {

	// TODO: release construction should be project specific rule.  config rules?
	release := ctx.String("service") + "-" + ctx.String("branch")

	// helm required flags
	args := []string{"--install", release, "--namespace", ctx.String("namespace")}

	// config file boolean flags
	for _, flag := range cfg.Workflow.CDProvider.Helm.Options.Flags {
		args = append(args, flag)
	}

	// cli flag conversion
	if ctx.GlobalBool("debug") {
		args = append(args, "--debug")
	}

	// convert cicd --dryrun arg to helm dialect
	if ctx.Bool("dryrun") {
		args = append(args, "--dry-run")
	}

	// write runtime helm --values <file> using when available in config  otherwise create/remove a TempFile
	outFile := cfg.Workflow.CDProvider.Helm.Options.Values.Output
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
	err = renderHelmValuesFile(ctx, cfg, valuesFile, ctx.String("repo"), ctx.String("tag"))
	if err != nil {
		return fmt.Errorf("renderHelmValuesFile(): %v", err)
	}

	// join flags and positional args
	args = append(args, "--values", valuesFile.Name())
	args = append(args, ctx.String("chart"))

	var stderr bytes.Buffer
	var cmdOut []byte

	// prepend subcommand deploy to args
	args = append([]string{"upgrade"}, args...)
	cmd := exec.Command("helm", args...)
	log.Println("execute: ", strings.Join(cmd.Args, " "))

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

func renderHelmValuesFile(c *cli.Context, cfg *Config, valuesFile *os.File, repo string, tag string) error {
	type Values struct {
		Repo, Tag, ServiceType string
	}

	// Prepare some data to insert into the template.
	var values = Values{Repo: repo, Tag: tag}

	// initialize the template
	var t *template.Template
	var err error
	if t, err = template.ParseFiles(cfg.Workflow.CDProvider.Helm.Options.Values.Template); err != nil {
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

	log.Println(string(yaml))

	return err
}
