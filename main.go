package main

import (
	"os"
	"text/template"
	"dockerfile-template/pkg"
)

func main() {
	data := dockerfile.DockerfileData{
		Stages: []dockerfile.Stage{
			{
				From:    "Task 1",
				As:      "test",
				Workdir: "adf",
				Expose: "80/tcp",
				Args: []dockerfile.Arg{
					{"test", "vale", true},
				},
				EnvVariables: []dockerfile.EnvVariable{
					{Name: "test"},
				},
				RunCommands: []dockerfile.RunCommand{
					{Command: "echo 1"},
				},
				CopyCommands: []dockerfile.CopyCommand{
					{Command: "tesjt3"},
				},
				Cmd: dockerfile.Cmd{
					Command: "some cmd",
				},
			},
			{From: "Task 1", As: "test2"},
		},
	}

	tmpl := template.Must(template.ParseFiles("template/dockerfile.template"))
	err := tmpl.Execute(os.Stdout, data)
	if err != nil {
		println(err)
	}
}
