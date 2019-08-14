package main

import (
	dft "dockerfile-template/pkg/dockerfile-template"
	"os"
)

func main() {
	data := &dft.DockerfileData{
		Stages: []dft.Stage{
			{
				User: "User 1",
				From: dft.From{
					Image: "image:latest",
					As:    "Base image",
				},
				Workdir: "adf",
				Expose:  "80/tcp",
				BuildArgs: []dft.Arg{
					{"test", "vale", true},
				},
				EnvVariables: []dft.EnvVariable{
					{Name: "test"},
				},
				RunCommands: []dft.RunCommand{
					{Command: "echo 1"},
				},
				CopyCommands: []dft.CopyCommand{
					{Command: "tesjt3"},
				},
				Cmd: &dft.Cmd{
					ShellExecForm: dft.ShellExecForm{
						Params: []string{ "cmd1", "cmd2" },
					},
				},
				Entrypoint: &dft.Entrypoint{
					ShellExecForm: dft.ShellExecForm{
						Params: []string{ "entrypoint", "param" },
					},
				},
				Volumes: []dft.Volume{
					{Source: "/App", Destination: "/opt/App"},
				},
			},
			{
				From: dft.From{
					Image: "image:latest",
					As:    "Second image",
				},
			},
		},
	}

	tmpl := dft.NewDockerFileTemplate(data)
	tmpl.Render(os.Stdout, "template/dockerfile.template")
}
