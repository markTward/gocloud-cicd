FROM golang:1.8.1

ADD . /go/src/github.com/markTward/gocloud-cicd

RUN go get -v -d ./...
RUN go install -v github.com/markTward/gocloud-cicd