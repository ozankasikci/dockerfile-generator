package dockerfile_generator

import (
	"fmt"
	"strings"
)

type RunForm string
const (
	ExecForm RunForm = "ExecForm"
	ShellForm RunForm = "ShellForm"
)

type Instruction interface {
	Render() string
}

type DockerfileData struct {
	Stages []Stage
}

type Stage []Instruction

type Arg struct {
	Name  string
	Value string
	Test  bool
	EnvVariable bool
}

func (a Arg) Render() string {
	res := fmt.Sprintf("ARG %s", a.Name)

	if	a.Value != "" {
		res = fmt.Sprintf("%s=%s", res, a.Value)
	}

	if a.Test {
		res = fmt.Sprintf("%s\nRUN test -n \"${%s}\"", res, a.Name)
	}

	if	a.EnvVariable {
		res = fmt.Sprintf("%s\nENV %s=\"${%s}\"", a.Name)
	}

	return res
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

func (l Label) Render() string {
	return fmt.Sprintf("LABEL %s=%s", l.Name, l.Value)
}

type Volume struct {
	Source      string
	Destination string
}

func (v Volume) Render() string {
	return fmt.Sprintf("VOLUME %s %s", v.Source, v.Destination)
}

type RunCommand struct {
	Params
	RunForm
}

func (r RunCommand) Render() string {
	if r.RunForm == "" || r.RunForm == ShellForm {
		return fmt.Sprintf("RUN %s", r.ShellForm())
	} else {
		return fmt.Sprintf("RUN %s", r.ExecForm())
	}
}

type EnvVariable struct {
	Name  string
	Value string
}

func (e EnvVariable) Render() string {
	return fmt.Sprintf("ENV %s=%s", e.Name, e.Value)
}

type CopyCommand struct {
	Sources []string
	Destination string
	Chown string
	From string
}

func (c CopyCommand) Render() string {
	res := "COPY"

	if c.From != "" {
		res = fmt.Sprintf("%s --from=%s", res, c.From)
	}

	if c.Chown != "" {
		res = fmt.Sprintf("%s --chown=%s", res, c.Chown)
	}

	sources := strings.Join(c.Sources, " ")
	res = fmt.Sprintf("%s %s %s", res, sources, c.Destination)

	return res
}

type Cmd struct {
	Params
}

func (c Cmd) Render() string {
	return fmt.Sprintf("CMD %s", c.ExecForm())
}

type Entrypoint struct {
	Params
}

func (e Entrypoint) Render() string {
	return fmt.Sprintf("ENTRYPOINT %s", e.ExecForm())
}

type Onbuild struct {
	Params
}

func (o Onbuild) Render() string {
	return fmt.Sprintf("ONBUILD %s", o.ShellForm())
}

type HealthCheck struct {
	Params
}

func (h HealthCheck) Render() string {
	return fmt.Sprintf("HEALTHCHECK %s", h.ShellForm())
}

type Shell struct {
	Params
}

func (s Shell) Render() string {
	return fmt.Sprintf("SHELL %s", s.ExecForm())
}

type Workdir struct {
	Dir string
}

func (w Workdir) Render() string {
	return fmt.Sprintf("WORKDIR %s", w.Dir)
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
