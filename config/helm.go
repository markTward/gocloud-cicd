package config

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
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

func (h *Helm) Deploy(cfg *Config, args []string) (err error) {

	var stderr bytes.Buffer
	var cmdOut []byte

	// prepend subcommand deploy to args
	args = append([]string{"upgrade"}, args...)
	cmd := exec.Command("helm", args...)
	log.Println("execute: ", strings.Join(cmd.Args, " "))

	cmd.Stderr = &stderr
	if cmdOut, err = cmd.Output(); err != nil {
		logCmdOutput(stderr.Bytes())
		err = fmt.Errorf("%v", stderr.String())
	} else {
		logCmdOutput(cmdOut)
	}

	return err
}
