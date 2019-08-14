package dockerfile_template

import (
	"fmt"
	"strings"
)

type DockerfileData struct {
	Stages []Stage
}

type Stage struct {
	BuildArgs    []Arg
	From         From
	Labels       []Label
	User         string
	RunCommands  []RunCommand
	Workdir      string
	EnvVariables []EnvVariable
	CopyCommands []CopyCommand
	Volumes      []Volume
	Cmd          *Cmd
	Entrypoint   *Entrypoint
	Expose       string
}

type Arg struct {
	Name  string
	Value string
	Test  bool
}

type From struct {
	Image string
	As    string
}

type Label struct {
	Name  string
	Value string
}

type Volume struct {
	Source      string
	Destination string
}

type RunCommand struct {
	Command string
}

type EnvVariable struct {
	Name  string
	Value string
}

type CopyCommand struct {
	Command string
}

type Cmd struct {
	ShellExecForm
}

type Entrypoint struct {
	ShellExecForm
}

type ShellExecForm struct {
	Params []string
}

func (e ShellExecForm) MapParams(f func(string) string) []string {
	res := make([]string, len(e.Params))
	for i, v := range e.Params {
		res[i] = f(v)
	}
	return res
}

func (e ShellExecForm) ExecForm() string {
	params := e.MapParams(func(s string) string {
		return fmt.Sprintf("\"%s\"", s)
	})

	paramsString := strings.Join(params, ", ")
	execFormString := fmt.Sprintf("[%s]", paramsString)

	return execFormString
}

func (e ShellExecForm) ShellForm() string {
	return strings.Join(e.Params, " ")
}

