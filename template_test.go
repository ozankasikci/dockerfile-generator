package dockerfile_generator

import (
	"bytes"
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestCodeRendering(t *testing.T) {
	data := &DockerfileData{
		Stages: []Stage{
			// Stage 1 - Builder
			[]Instruction{
				From{
					Image: "golang:1.7.3", As: "builder",
				},
				Workdir{
					Dir: "/go/src/github.com/alexellis/href-counter/",
				},
				RunCommand{
					Params: []string{"go", "get", "-d", "-v", "golang.org/x/net/html"},
				},
				CopyCommand{
					Sources: []string{"app.go"}, Destination: ".",
				},
				RunCommand{
					Params: []string{"CGO_ENABLED=0", "GOOS=linux", "go", "build", "-a", "-installsuffix", "cgo", "-o", "app", "."},
				},
			},
			// Stage 2 - Final
			[]Instruction{
				From{
					Image: "alpine:latest", As: "final",
				},
				RunCommand{
					Params: []string{"apk", "--no-cache", "add", "ca-certificates"},
				},
				Workdir{
					Dir: "/root/",
				},
				CopyCommand{
					From: "builder", Sources: []string{"/go/src/github.com/alexellis/href-counter/app"}, Destination: ".",
				},
				Cmd{
					Params: []string{"./app"},
				},
			},
		},
	}

	tmpl := NewDockerFileTemplate(data)
	output := &bytes.Buffer{}
	err := tmpl.Render(output)
	assert.NoError(t, err)

    expectedOutput := `FROM golang:1.7.3 as builder
WORKDIR /go/src/github.com/alexellis/href-counter/
RUN go get -d -v golang.org/x/net/html
COPY app.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest as final
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/alexellis/href-counter/app .
CMD ["./app"]

`

	assert.Equal(t, expectedOutput, output.String())
}

