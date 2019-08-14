package main

import (
	dft "dockerfile-template/pkg/dockerfile-template"
	"os"
)

func main() {
	//_ = &dft.DockerfileData{
	//	Stages: []dft.Stage{
	//		{
	//			User: "User 1",
	//			From: dft.From{
	//				Image: "image:latest",
	//				As:    "Base image",
	//			},
	//			Workdir:    "adf",
	//			Expose:     "80/tcp",
	//			StopSignal: "TERM",
	//			BuildArgs: []dft.Arg{
	//				{"test", "vale", true},
	//			},
	//			EnvVariables: []dft.EnvVariable{
	//				{Name: "test"},
	//			},
	//			RunCommands: []dft.RunCommand{
	//				{Params: dft.Params{Params: []string{"cmd1", "cmd2"}}},
	//			},
	//			CopyCommands: []dft.CopyCommand{
	//				{Command: "tesjt3"},
	//			},
	//			Cmd: &dft.Cmd{
	//				Params: dft.Params{Params: []string{"cmd1", "cmd2"}},
	//			},
	//			Entrypoint: &dft.Entrypoint{
	//				Params: dft.Params{Params: []string{"entrypoint", "param"}},
	//			},
	//			Onbuild: &dft.Onbuild{
	//				Params: dft.Params{Params: []string{"RUN", "echo", "1"}},
	//			},
	//			HealthCheck: &dft.HealthCheck{
	//				Params: dft.Params{Params: []string{"--interval=DURATION", "CMD", "command"}},
	//			},
	//			Shell: &dft.Shell{
	//				Params: dft.Params{Params: []string{"powershell", "command"}},
	//				RunCommands: []dft.RunCommand{
	//					{Params: dft.Params{Params: []string{"shell cmd 1", "shell cmd 2"}}},
	//				},
	//			},
	//			Volumes: []dft.Volume{
	//				{Source: "/App", Destination: "/opt/App"},
	//			},
	//		},
	//	},
	//}

	//tmpl := dft.NewDockerFileTemplate(data)
	//err := tmpl.Render(os.Stdout, "template/dockerfile.template")

	dataMap := &dft.DockerfileDataSlice{
		Stages: []dft.StageSlice{
			[]dft.Instruction{
				dft.From{Image: "debian", As: "deb"},
				dft.RunCommand{
					Params: []string{"dsf"}},
				},
			[]dft.Instruction{
				dft.From{Image: "debian", As: "real"},
				dft.RunCommand{
					Params: []string{"dsf"}},
			},
			},
		}
	tmpl := dft.NewDockerFileMapTemplate(dataMap)
	err := tmpl.Render(os.Stdout, "template/dockerfile.map.template")

	if err != nil {
		println(err.Error())
	}
}
