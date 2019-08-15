### Dockerfile Generator

Automated Dockerfile generation library.

### Docs

https://godoc.org/github.com/ozankasikci/dockerfile-generator

### Example

```go
package main

import dfg "github.com/ozankasikci/dockerfile-generator"

func main() {
	data := &dfg.DockerfileData{
		Stages: []dfg.Stage{
    			// Stage 1 - Builder Image
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
	tmpl := dfg.NewDockerFileTemplate(data)
	
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
