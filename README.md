## dfg - Dockerfile Generator
[![Build Status](https://travis-ci.org/ozankasikci/dockerfile-generator.svg?branch=master)](https://travis-ci.org/ozankasikci/dockerfile-generator)

Automated Dockerfile generation library.

### Installation

`go get -u github.com/ozankasikci/dockerfile-generator`

### Docs

https://godoc.org/github.com/ozankasikci/dockerfile-generator

### Library Usage Example

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

### Output
```Dockerfile
FROM golang:1.7.3 as builder
WORKDIR /go/src/github.com/alexellis/href-counter/
RUN go get -d -v golang.org/x/net/html
COPY app.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest as final
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/alexellis/href-counter/app .
CMD ["./app"]
```

### YAML File Example
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
          - apt-get clean &&
          - rm -rf /var/lib/apt/lists/*
    - cmd:
        params:
          - /usr/sbin/apache2
          - -D
          - FOREGROUND
```
```
data, err := NewDockerFileDataFromYamlFile("./example-input-files/test-input.yaml")
tmpl := NewDockerfileTemplate(data)
err = tmpl.Render(output)
```

### Output

```dockerfile
FROM kstaken/apache2
RUN apt-get update && apt-get install -y php5 apt-get clean && rm -rf /var/lib/apt/lists/*
CMD ["/usr/sbin/apache2", "-D", "FOREGROUND"]
```