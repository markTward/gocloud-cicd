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
	Release   string
	Namespace string
	ChartPath string
}

// helm upgrade \
// $DRYRUN_OPTION \
// --debug \
// --install $RELEASE_NAME \
// --namespace=$NAMESPACE \
// --set service.gocloudAPI.image.repository=$DOCKER_REPO \
// --set service.gocloudAPI.image.tag=":$COMMIT_TAG" \
// --set service.gocloudGrpc.image.repository=$DOCKER_REPO \
// --set service.gocloudGrpc.image.tag=":$COMMIT_TAG" \
// $CHARTPATH

func (h *Helm) Deploy(args []string) (err error) {
	var stderr bytes.Buffer
	var cmdOut []byte

	//TODO: add args to command
	cmd := exec.Command("helm", "help")
	cmd.Stderr = &stderr
	log.Println(strings.Join(cmd.Args, " "))

	if cmdOut, err = cmd.Output(); err != nil {
		logCmdOutput(stderr.Bytes())
		err = fmt.Errorf("%v", stderr.String())
	} else {
		logCmdOutput(cmdOut)
	}

	return err
}
