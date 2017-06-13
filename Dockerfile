FROM golang:1.8.1

# get docker client; minikube running on 1.11.1; gke ?
RUN wget -qO- https://get.docker.com/builds/Linux/x86_64/docker-1.11.1.tgz | tar xvz && mv docker/docker /usr/local/bin

ADD . /go/src/github.com/markTward/gocloud-cicd

RUN go get -v -d ./...
RUN go install -v github.com/markTward/gocloud-cicd