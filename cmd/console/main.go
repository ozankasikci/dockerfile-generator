package main

import (
	dfg "github.com/ozankasikci/dockerfile-generator"
	"os"
)

func main() {
	data := &dfg.DockerfileData{
		Stages: []dfg.Stage{
			[]dfg.Instruction{
				dfg.From{Image: "golang:1.7.3", As: "builder"},
				dfg.Workdir{Dir: "/go/src/github.com/alexellis/href-counter/"},
				dfg.RunCommand{Params: []string{"go", "get", "-d", "-v", "golang.org/x/net/html"}},
				dfg.CopyCommand{Sources: []string{"app.go"}, Destination: "."},
				dfg.RunCommand{ Params: []string{"CGO_ENABLED=0", "GOOS=linux", "go", "build", "-a", "-installsuffix", "cgo", "-o", "app", "."}, },
			},
			[]dfg.Instruction{
				dfg.From{Image: "alpine:latest", As: "final"},
				dfg.RunCommand{Params: []string{"apk", "--no-cache", "add", "ca-certificates"}},
				dfg.Workdir{Dir: "/root/"},
				dfg.CopyCommand{From: "builder", Sources: []string{"/go/src/github.com/alexellis/href-counter/app"}, Destination: "."},
				dfg.Cmd{Params: []string{"./app"}},
			},
		},
	}
	tmpl := dfg.NewDockerFileTemplate(data)
	err := tmpl.Render(os.Stdout)

	if err != nil {
		println(err.Error())
	}
}
