## dfg - Dockerfile Generator

`dfg` is both a go library and an executable that produces valid Dockerfiles using various input channels.

[![Build Status](https://travis-ci.org/ozankasikci/dockerfile-generator.svg?branch=master)](https://travis-ci.org/ozankasikci/dockerfile-generator)
[![GoDoc](https://godoc.org/github.com/ozankasikci/dockerfile-generator?status.svg)](https://godoc.org/github.com/ozankasikci/dockerfile-generator)
[![Coverage Status](https://coveralls.io/repos/github/ozankasikci/dockerfile-generator/badge.svg?branch=master)](https://coveralls.io/github/ozankasikci/dockerfile-generator?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ozankasikci/dockerfile-generator)](https://goreportcard.com/report/github.com/ozankasikci/dockerfile-generator)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) 

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
  * [Installing as an Executable](#installing-as-an-executable)
  * [Installing as a Library](#installing-as-a-library)
- [Getting Started](#getting-started)
  * [Using dfg as an Executable](#using-dfg-as-an-executable)
  * [Using dfg as a Library](#using-dfg-as-a-library)
- [Examples](#examples)
  * [YAML File Example](#yaml-file-example)
  * [Library Usage Example](#library-usage-example)
- [TODO](#todo)

## Overview

`dfg` is a Dockerfile generator that accepts input data from various sources, produces and redirects the generated Dockerfile to an output target such as a file or stdout.

It is especially useful for generating Dockerfile instructions conditionally since Dockerfile language has no control flow logic.

## Installation
#### Installing as an Executable

* MacOS

```shell
curl -o dfg -L https://github.com/ozankasikci/dockerfile-generator/releases/download/v0.1.0/dfg_v0.1.0_darwin_amd64
chmod +x dfg && sudo mv dfg /usr/local/bin
```

* Linux

```shell
curl -o dfg -L https://github.com/ozankasikci/dockerfile-generator/releases/download/v0.1.0/dfg_v0.1.0_linux_amd64
chmod +x dfg && sudo mv dfg /usr/local/bin
```

* Windows

```shell
curl -o dfg.exe -L https://github.com/ozankasikci/dockerfile-generator/releases/download/v0.1.0/dfg_v0.1.0_windows_amd64.exe
```

#### Installing as a Library

`go get -u github.com/ozankasikci/dockerfile-generator`

## Getting Started

### Using dfg as an Executable

Available commands:

`dfg generate --input path/to/yaml --out Dockerfile` generates a file named `Dockerfile`

`dfg generate --help` lists available flags

### Using dfg as a Library

When using `dfg` as a go library, you need to pass a `[]dfg.Stage` slice as data.
This approach enables and encourages multi staged Dockerfiles.
Dockerfile instructions will be generated in the same order as in the `[]dfg.Instruction` slice.

Some `Instruction`s accept a `runForm` field which specifies if the `Instruction` should be run in the `shell form` or the `exec form`.
If the `runForm` is not specified, it will be chosen based on [Dockerfile best practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/). 

For detailed usage example please see [Library Usage Example](#library-usage-example)

## Examples

#### YAML File Example
```yaml
stages:
  final:
    - from:
        image: kstaken/apache2
    - run:
        runForm: shell
        params:
          - apt-get update &&
          - apt-get install -y
          - php5
          - libapache2-mod-php5 &&
          - apt-get clean &&
          - rm -rf /var/lib/apt/lists/*
    - cmd:
        params:
          - /usr/sbin/apache2
          - -D
          - FOREGROUND
```
Use dfg as binary:
```shell
dfg generate -i ./example-input-files/apache-php.yaml --stdout
```
Or as a library
```go
data, err := dfg.NewDockerFileDataFromYamlFile("./example-input-files/apache-php.yaml")
tmpl := dfg.NewDockerfileTemplate(data)
err = tmpl.Render(output)
```

#### Output

```dockerfile
FROM kstaken/apache2
RUN apt-get update && apt-get install -y php5 libapache2-mod-php5 && apt-get clean && rm -rf /var/lib/apt/lists/*
CMD ["/usr/sbin/apache2", "-D", "FOREGROUND"]
```

#### Library Usage Example

```go
package main

import dfg "github.com/ozankasikci/dockerfile-generator"

func main() {
	data := &dfg.DockerfileData{
		Stages: []dfg.Stage{
			// Stage 1 - Builder Image
			// An instruction is just an interface, so you can pass custom structs as well
			[]dfg.Instruction{
				dfg.From{
					Image: "golang:1.7.3", As: "builder",
				},
				dfg.User{
					User: "ozan",
				},
				dfg.Workdir{
					Dir: "/go/src/github.com/alexellis/href-counter/",
				},
				dfg.RunCommand{
					Params: []string{"go", "get", "-d", "-v", "golang.org/x/net/html"},
				},
				dfg.CopyCommand{
					Sources: []string{"app.go"}, Destination: ".",
				},
				dfg.RunCommand{
					Params: []string{"CGO_ENABLED=0", "GOOS=linux", "go", "build", "-a", "-installsuffix", "cgo", "-o", "app", "."},
				},
			},
			// Stage 2 - Final Image
			[]dfg.Instruction{
				dfg.From{
					Image: "alpine:latest", As: "final",
				},
				dfg.RunCommand{
					Params: []string{"apk", "--no-cache", "add", "ca-certificates"},
				},
				dfg.User{
					User: "root", Group: "admin",
				},
				dfg.Workdir{
					Dir: "/root/",
				},
				dfg.CopyCommand{
					From: "builder", Sources: []string{"/go/src/github.com/alexellis/href-counter/app"}, Destination: ".",
				},
				dfg.Cmd{
					Params: []string{"./app"},
				},
			},
		},
	}
	tmpl := dfg.NewDockerfileTemplate(data)
	
	// write to a file
	file, err := os.Create("Dockerfile")
	err = tmpl.Render(file)
	
	// or write to stdout
	err = tmpl.Render(os.Stdout)
}
``` 

#### Output
```Dockerfile
FROM golang:1.7.3 as builder
USER ozan
WORKDIR /go/src/github.com/alexellis/href-counter/
RUN go get -d -v golang.org/x/net/html
COPY app.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest as final
RUN apk --no-cache add ca-certificates
USER root:admin
WORKDIR /root/
COPY --from=builder /go/src/github.com/alexellis/href-counter/app .
CMD ["./app"]
```

## TODO
- [ ] Add reading Dockerfile data from an existing yaml file support
- [ ] Implement json file input channel
- [ ] Implement stdin input channel
- [ ] Implement toml file input channel

