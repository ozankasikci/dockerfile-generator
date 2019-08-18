package dockerfile_generator

import (
	"fmt"
	"strings"
)

type RunForm string

const (
	ExecForm                 RunForm = "ExecForm"
	ShellForm                RunForm = "ShellForm"
	RunCommandDefaultRunForm         = ShellForm
	CmdDefaultRunForm                = ExecForm
	EntrypointDefaultRunForm         = ExecForm
)

type Instruction interface {
	Render() string
}

// DockerfileData struct can hold multiple stages for a multi-staged Dockerfile
// Check https://docs.docker.com/develop/develop-images/multistage-build/ for more information
type DockerfileData struct {
	Stages []Stage `yaml:"stages,omitempty"`
}

type Stage []Instruction

func (s *Stage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data []interface{}
	var result []Instruction
	err := unmarshal(&data)
	if err != nil {
		return err
	}

	*s = append(result, cleanUpInterfaceArray(data)...)
	return nil
}

type Arg struct {
	Name        string `yaml:"name"`
	Value       string `yaml:"value"`
	Test        bool   `yaml:"test,omitempty"`
	EnvVariable bool   `yaml:"envVariable,omitempty"`
}

func (a Arg) Render() string {
	res := fmt.Sprintf("ARG %s", a.Name)

	if a.Value != "" {
		res = fmt.Sprintf("%s=%s", res, a.Value)
	}

	if a.Test {
		res = fmt.Sprintf("%s\nRUN test -n \"${%s}\"", res, a.Name)
	}

	if a.EnvVariable {
		res = fmt.Sprintf("%s\nENV %s=\"${%s}\"", res, a.Name, a.Name)
	}

	return res
}

type From struct {
	Image string `yaml:"image"`
	As    string `yaml:"as"`
}

func (f From) Render() string {
	res := fmt.Sprintf("FROM %s", f.Image)

	if f.As != "" {
		res = fmt.Sprintf("%s as %s", res, f.As)
	}

	return res
}

type Label struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func (l Label) Render() string {
	return fmt.Sprintf("LABEL %s=%s", l.Name, l.Value)
}

type Volume struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

func (v Volume) Render() string {
	return fmt.Sprintf("VOLUME %s %s", v.Source, v.Destination)
}

type RunCommand struct {
	Params  `yaml:"params"`
	RunForm `yaml:"runForm"`
}

func (r RunCommand) Render() string {
	if r.RunForm == "" {
		r.RunForm = RunCommandDefaultRunForm
	}

	if r.RunForm == ExecForm {
		return fmt.Sprintf("RUN %s", r.ExecForm())
	} else {
		return fmt.Sprintf("RUN %s", r.ShellForm())
	}
}

type EnvVariable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func (e EnvVariable) Render() string {
	return fmt.Sprintf("ENV %s=%s", e.Name, e.Value)
}

type CopyCommand struct {
	Sources     []string `yaml:"sources"`
	Destination string   `yaml:"destination"`
	Chown       string   `yaml:"chown"`
	From        string   `yaml:"from"`
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
	Params  `yaml:"params"`
	RunForm `yaml:"runForm"`
}

func (c Cmd) Render() string {
	if c.RunForm == "" {
		c.RunForm = ExecForm
	}

	if c.RunForm == ExecForm {
		return fmt.Sprintf("CMD %s", c.ExecForm())
	} else {
		return fmt.Sprintf("CMD %s", c.ShellForm())
	}
}

type Entrypoint struct {
	Params  `yaml:"params"`
	RunForm `yaml:"runForm"`
}

func (e Entrypoint) Render() string {
	if e.RunForm == "" {
		e.RunForm = ExecForm
	}

	if e.RunForm == ExecForm {
		return fmt.Sprintf("ENTRYPOINT %s", e.ExecForm())
	} else {
		return fmt.Sprintf("ENTRYPOINT %s", e.ShellForm())
	}
}

type Onbuild struct {
	Params `yaml:"params"`
}

func (o Onbuild) Render() string {
	return fmt.Sprintf("ONBUILD %s", o.ShellForm())
}

type HealthCheck struct {
	Params `yaml:"params"`
}

func (h HealthCheck) Render() string {
	return fmt.Sprintf("HEALTHCHECK %s", h.ShellForm())
}

type Shell struct {
	Params `yaml:"params"`
}

func (s Shell) Render() string {
	return fmt.Sprintf("SHELL %s", s.ExecForm())
}

type Workdir struct {
	Dir string `yaml:"dir"`
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
