package main

import (
	"github.com/ozankasikci/dockerfile-generator"
	"os"
)

func main() {
	dataMap := &dockerfile_generator.DockerfileData{
		Stages: []dockerfile_generator.Stage{
			[]dockerfile_generator.Instruction{
				dockerfile_generator.From{Image: "node:8.15.0-alpine", As: "builder"},
				dockerfile_generator.Arg{Name: "WORKDIR", Value: "/app"},
				dockerfile_generator.Arg{Name: "VIEWER_DIR", Value: "$WORKDIR/client"},
				dockerfile_generator.Workdir{Dir: "$WORKDIR"},
				dockerfile_generator.RunCommand{Params: []string{"apk add --no-cache python py-pip git curl openssh"}},
				dockerfile_generator.Arg{Name: "TARGET_GIT_BRANCH", Test: true},
				dockerfile_generator.RunCommand{Params: []string{"mkdir", "/root/.ssh/"}},
				dockerfile_generator.RunCommand{Params: []string{"echo", "\"${SSH_KEY}\"", ">", "/root/.ssh"}},
				dockerfile_generator.RunCommand{Params: []string{"chmod", "600", ">", "/root/.ssh"}},
			},
			[]dockerfile_generator.Instruction{
				dockerfile_generator.From{Image: "nginx:1.14-pearl", As: "final"},
				dockerfile_generator.Arg{Name: "WORKDIR", Value: "/app"},
				dockerfile_generator.Arg{Name: "INSTALL_DIR", Value: "/opt/company"},
				dockerfile_generator.CopyCommand{From: "builder", Sources: []string{"$CMD/build"}, Destination: "$BUILD_DIR"},
				dockerfile_generator.Cmd{Params: []string{"nginx", "-g", "daemon off;"}},
			},
		},
	}
	tmpl := dockerfile_generator.NewDockerFileTemplate(dataMap)
	err := tmpl.Render(os.Stdout)

	if err != nil {
		println(err.Error())
	}
}
