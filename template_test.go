package dockerfilegenerator

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCodeRendering(t *testing.T) {
	data := &DockerfileData{
		Stages: []Stage{
			// Stage 1 - Builder
			[]Instruction{
				From{
					Image: "golang:1.7.3", As: "builder",
				},
				Arg{
					Name: "arg-name", Test: true, EnvVariable: true,
				},
				Workdir{
					Dir: "/go/src/github.com/alexellis/href-counter/",
				},
				User{
					User: "ozan",
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

	tmpl := NewDockerfileTemplate(data)
	output := &bytes.Buffer{}
	err := tmpl.Render(output)
	assert.NoError(t, err)

	expectedOutput := `FROM golang:1.7.3 as builder
ARG arg-name
RUN test -n "${arg-name}"
ENV arg-name="${arg-name}"
WORKDIR /go/src/github.com/alexellis/href-counter/
USER ozan
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

func TestCodeRenderingWithoutUser(t *testing.T) {
	data := &DockerfileData{
		Stages: []Stage{
			// Stage 1 - Builder
			[]Instruction{
				From{
					Image: "golang:1.7.3", As: "builder",
				},
			},
			[]Instruction{
				From{
					Image: "alpine:latest", As: "final",
				},
				Cmd{
					Params: []string{"./app"},
				},
			},
		},
	}

	tmpl := NewDockerfileTemplate(data)
	output := &bytes.Buffer{}
	err := tmpl.Render(output)
	assert.NoError(t, err)

	expectedOutput := `FROM golang:1.7.3 as builder

FROM alpine:latest as final
CMD ["./app"]

`

	assert.Equal(t, expectedOutput, output.String())
}
func TestCodeRenderingWithUserMap(t *testing.T) {
	data := &DockerfileData{
		Stages: []Stage{
			// Stage 1 - Builder
			[]Instruction{
				From{
					Image: "golang:1.7.3", As: "builder",
				},
				User{
					User: "ozan", Group: "admin",
				},
			},
		},
	}

	tmpl := NewDockerfileTemplate(data)
	output := &bytes.Buffer{}
	err := tmpl.Render(output)
	assert.NoError(t, err)

	expectedOutput := `FROM golang:1.7.3 as builder
USER ozan:admin

`

	assert.Equal(t, expectedOutput, output.String())
}

func TestYamlRendering(t *testing.T) {
	data, err := NewDockerFileDataFromYamlFile("./example-input-files/test-input.yaml")
	tmpl := NewDockerfileTemplate(data)
	assert.NoError(t, err)

	output := &bytes.Buffer{}
	err = tmpl.Render(output)
	assert.NoError(t, err)

	expectedOutput := `FROM alpine:latest as builder
WORKDIR /app
USER ozan
ARG test-arg=arg-value
RUN test -n "${test-arg}"
ENV test-arg="${test-arg}"
VOLUME some/source ./some/destination
RUN echo "test" 1
ENV env=dev
COPY --chown=me:me /etc/conf /opt/app/conf
ONBUILD echo test

FROM alpine:latest as final
ARG test-arg=arg-value
RUN test -n "${test-arg}"
ENV test-arg="${test-arg}"
LABEL label1=label-value
ENV DB_PASSWORD=password
CMD echo test
ENTRYPOINT ["echo", "test"]
HEALTHCHECK --interval=DURATION --timeout=3s CMD curl -f http://localhost/
SHELL ["powershell", "-command"]
WORKDIR test dir

`

	assert.Equal(t, expectedOutput, output.String())
}

func TestYamlRenderingNoUser(t *testing.T) {
	data, err := NewDockerFileDataFromYamlFile("./example-input-files/test-input-no-user.yaml")
	tmpl := NewDockerfileTemplate(data)
	assert.NoError(t, err)

	output := &bytes.Buffer{}
	err = tmpl.Render(output)
	assert.NoError(t, err)

	expectedOutput := `FROM alpine:latest as builder

FROM alpine:latest as final

`

	assert.Equal(t, expectedOutput, output.String())
}

func TestYamlRenderingUserAndGroup(t *testing.T) {
	data, err := NewDockerFileDataFromYamlFile("./example-input-files/test-input-user-group.yaml")
	tmpl := NewDockerfileTemplate(data)
	assert.NoError(t, err)

	output := &bytes.Buffer{}
	err = tmpl.Render(output)
	assert.NoError(t, err)

	expectedOutput := `FROM alpine:latest as builder
USER ozan:admin

FROM alpine:latest as final
USER 1000:1000

`

	assert.Equal(t, expectedOutput, output.String())
}

func TestYamlRenderingFail(t *testing.T) {
	data, err := NewDockerFileDataFromYamlFile("./example-input-files/invalid-input.yaml")
	tmpl := NewDockerfileTemplate(data)
	assert.Error(t, err)

	output := &bytes.Buffer{}
	err = tmpl.Render(output)
	assert.Error(t, err)
}

func TestInvalidYamlFilePath(t *testing.T) {
	_, err := NewDockerFileDataFromYamlFile("non-existent.yaml")
	assert.EqualError(t, err, "Unmarshal: yamlFile.Get err #open non-existent.yaml: no such file or directory")
}
