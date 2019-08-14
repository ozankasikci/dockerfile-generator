package main

import (
	dft "github.com/ozankasikci/dockerfile-generator/pkg/dockerfile-generator"
	"os"
)

func main() {
	dataMap := &dft.DockerfileData{
		Stages: []dft.Stage{
			[]dft.Instruction{
				dft.From{Image: "node:8.15.0-alpine", As: "builder"},
				dft.Arg{Name: "WORKDIR", Value: "/app"},
				dft.Arg{Name: "VIEWER_DIR", Value: "$WORKDIR/client"},
				dft.Workdir{Dir: "$WORKDIR"},
				dft.RunCommand{Params: []string{"apk add --no-cache python py-pip git curl openssh"}},
				dft.Arg{Name: "TARGET_GIT_BRANCH", Test: true},
				dft.RunCommand{Params: []string{"mkdir", "/root/.ssh/"}},
				dft.RunCommand{Params: []string{"echo", "\"${SSH_KEY}\"", ">", "/root/.ssh"}},
				dft.RunCommand{Params: []string{"chmod", "600", ">", "/root/.ssh"}},
			},
			[]dft.Instruction{
				dft.From{Image: "nginx:1.14-pearl", As: "final"},
				dft.Arg{Name: "WORKDIR", Value: "/app"},
				dft.Arg{Name: "INSTALL_DIR", Value: "/opt/company"},
				dft.CopyCommand{From: "builder", Sources: []string{"$CMD/build"}, Destination: "$BUILD_DIR"},
				dft.Cmd{Params: []string{"nginx", "-g", "daemon off;"}},
			},
		},
	}
	tmpl := dft.NewDockerFileTemplate(dataMap)
	err := tmpl.Render(os.Stdout, "template/dockerfile.template")

	if err != nil {
		println(err.Error())
	}
}
