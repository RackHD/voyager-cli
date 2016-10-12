ORGANIZATION = RackHD
PROJECT = voyager-cli
BINARYNAME = mcc
GOOUT = ./bin

default: deps build test

deps:
	go get github.com/onsi/ginkgo/ginkgo

	go get github.com/fatih/color
	go get github.com/onsi/gomega
	go get github.com/onsi/gomega/ghttp
	go get ./...
	env GOOS=windows go get ./...

integration-test: build
	ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=30 --focus="\bINTEGRATION\b"

unit-test: build
	ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=30 --focus="\bUNIT\b"

test: build
	ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=30

cover-cmd: test
	go tool cover -html=cmd/cmd.coverprofile

build: build-Linux build-Mac build-Windows

build-Linux:
	env GOOS=linux go build -o $(GOOUT)/linux/$(BINARYNAME)

build-Mac:
	env GOOS=darwin go build -o $(GOOUT)/darwin/$(BINARYNAME)

build-Windows:
	env GOOS=windows go build -o $(GOOUT)/windows/$(BINARYNAME).exe
