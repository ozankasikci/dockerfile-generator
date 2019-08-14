package main

import (
	"github.com/ozankasikci/dockerfile-generator/pkg"
	"os"
)

func main() {
	dataMap := &pkg.DockerfileData{
		Stages: []pkg.Stage{
			[]pkg.Instruction{
				pkg.From{Image: "node:8.15.0-alpine", As: "builder"},
				pkg.Arg{Name: "WORKDIR", Value: "/app"},
				pkg.Arg{Name: "VIEWER_DIR", Value: "$WORKDIR/client"},
				pkg.Workdir{Dir: "$WORKDIR"},
				pkg.RunCommand{Params: []string{"apk add --no-cache python py-pip git curl openssh"}},
				pkg.Arg{Name: "TARGET_GIT_BRANCH", Test: true},
				pkg.RunCommand{Params: []string{"mkdir", "/root/.ssh/"}},
				pkg.RunCommand{Params: []string{"echo", "\"${SSH_KEY}\"", ">", "/root/.ssh"}},
				pkg.RunCommand{Params: []string{"chmod", "600", ">", "/root/.ssh"}},
			},
			[]pkg.Instruction{
				pkg.From{Image: "nginx:1.14-pearl", As: "final"},
				pkg.Arg{Name: "WORKDIR", Value: "/app"},
				pkg.Arg{Name: "INSTALL_DIR", Value: "/opt/company"},
				pkg.CopyCommand{From: "builder", Sources: []string{"$CMD/build"}, Destination: "$BUILD_DIR"},
				pkg.Cmd{Params: []string{"nginx", "-g", "daemon off;"}},
			},
		},
	}
	tmpl := pkg.NewDockerFileTemplate(dataMap)
	err := tmpl.Render(os.Stdout)

	if err != nil {
		println(err.Error())
	}
}
