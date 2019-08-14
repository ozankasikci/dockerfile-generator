package dockerfile_template

import (
	"fmt"
	"strings"
)

type Instruction interface {
	Render() string
}

type DockerfileDataSlice struct {
	Stages []StageSlice
}

type StageSlice []Instruction

//type DockerfileData struct {
//	Stages []Stage
//}

//type Stage struct {
//	BuildArgs    []Arg
//	From         From
//	Labels       []Label
//	User         string
//	RunCommands  []RunCommand
//	Workdir      string
//	EnvVariables []EnvVariable
//	CopyCommands []CopyCommand
//	Volumes      []Volume
//	Cmd          *Cmd
//	Entrypoint   *Entrypoint
//	Onbuild      *Onbuild
//	HealthCheck  *HealthCheck
//	Shell        *Shell
//	StopSignal   string
//	Expose       string
//}

type Arg struct {
	Name  string
	Value string
	Test  bool
}

type From struct {
	Image string
	As    string
}

func (f From) Render() string {
	res := fmt.Sprintf("FROM %s", f.Image)

	if f.As != "" {
		res = fmt.Sprintf("%s as %s", res, f.As)
	}

	return res
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
	Params
}

func (r RunCommand) Render() string {
	return fmt.Sprintf("RUN %s", r.ExecForm())
}

type EnvVariable struct {
	Name  string
	Value string
}

type CopyCommand struct {
	Command string
}

type Cmd struct {
	Params
}

type Entrypoint struct {
	Params
}

type Onbuild struct {
	Params
}

type HealthCheck struct {
	Params
}

type Shell struct {
	Params
	RunCommands []RunCommand
}

type Params []string

func (p Params) mapParams(f func(string) string) []string {
	res := make([]string, len(p))
	for i, v := range p {
		res[i] = f(v)
	}
	return res
}

func (p Params) ExecForm() string {
	params := p.mapParams(func(s string) string {
		return fmt.Sprintf("\"%s\"", s)
	})

	paramsString := strings.Join(params, ", ")
	execFormString := fmt.Sprintf("[%s]", paramsString)

	return execFormString
}

func (p Params) ShellForm() string {
	return strings.Join(p, " ")
}
