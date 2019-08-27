package dockerfilegenerator

import (
	"fmt"
	"strings"
)

// RunForm specifies in which form the instruction string should be constructed.
// Check Dockerfile exec and shell forms for more information
type RunForm string

const (
	// ExecForm is essentially a json array of string, e.g. ["echo", "1"]
	ExecForm RunForm = "ExecForm"

	// ShellForm is the form of a usual terminal command, e.g. echo 1
	ShellForm RunForm = "ShellForm"

	// RunCommandDefaultRunForm is the default RunForm for RunCommand
	RunCommandDefaultRunForm = ShellForm

	// CmdDefaultRunForm is the default RunForm for Cmd
	CmdDefaultRunForm = ExecForm

	// EntrypointDefaultRunForm is the default RunForm for Entrypoint
	EntrypointDefaultRunForm = ExecForm
)

// Instruction represents a Dockerfile instruction, e.g. FROM alpine:latest
type Instruction interface {
	Render() string
}

// DockerfileData struct can hold multiple stages for a multi-staged Dockerfile
// Check https://docs.docker.com/develop/develop-images/multistage-build/ for more information
type DockerfileData struct {
	Stages []Stage `yaml:"stages,omitempty"`
}

// Stage is a set of instructions, the purpose is to keep the order of the given instructions
// and generate a Dockerfile using the output of these instructions
type Stage []Instruction

// UnmarshalYAML implements an interface to let go-yaml be able to decode Stages in to Stage struct
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

// Arg represents a Dockerfile instruction, see https://docs.docker.com/engine/reference/builder/#arg
type Arg struct {
	Name        string `yaml:"name"`
	Value       string `yaml:"value"`
	Test        bool   `yaml:"test,omitempty"`
	EnvVariable bool   `yaml:"envVariable,omitempty"`
}

// Render returns a string in the form of ARG <name>[=<default value>]
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

// From represents a Dockerfile instruction, see https://docs.docker.com/engine/reference/builder/#from
type From struct {
	Image string `yaml:"image"`
	As    string `yaml:"as"`
}

// Render returns a string in the form of FROM <image> [AS <name>]
func (f From) Render() string {
	res := fmt.Sprintf("FROM %s", f.Image)

	if f.As != "" {
		res = fmt.Sprintf("%s as %s", res, f.As)
	}

	return res
}

// Label represents a Dockerfile instruction, see https://docs.docker.com/engine/reference/builder/#from
type Label struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Render returns a string in the form of LABEL <key>=<value>
func (l Label) Render() string {
	return fmt.Sprintf("LABEL %s=%s", l.Name, l.Value)
}

// Volume represents a Dockerfile instruction, see https://docs.docker.com/engine/reference/builder/#volume
type Volume struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

// Render returns a string in the form of VOLUME <source> <destination>
func (v Volume) Render() string {
	return fmt.Sprintf("VOLUME %s %s", v.Source, v.Destination)
}

// RunCommand represents a Dockerfile instruction, see https://docs.docker.com/engine/reference/builder/#run
type RunCommand struct {
	Params  `yaml:"params"`
	RunForm `yaml:"runForm"`
}

// Render returns a string in the form of RUN <command>
func (r RunCommand) Render() string {
	if r.RunForm == "" {
		r.RunForm = RunCommandDefaultRunForm
	}

	if r.RunForm == ExecForm {
		return fmt.Sprintf("RUN %s", r.ExecForm())
	}

	return fmt.Sprintf("RUN %s", r.ShellForm())
}

// EnvVariable represents a Dockerfile instruction, see https://docs.docker.com/engine/reference/builder/#env
type EnvVariable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Render returns a string in the form of ENV <key> <value>
func (e EnvVariable) Render() string {
	return fmt.Sprintf("ENV %s=%s", e.Name, e.Value)
}

// CopyCommand represents a Dockerfile instruction, see https://docs.docker.com/engine/reference/builder/#copy
type CopyCommand struct {
	Sources     []string `yaml:"sources"`
	Destination string   `yaml:"destination"`
	Chown       string   `yaml:"chown"`
	From        string   `yaml:"from"`
}

// Render returns a string in the form of COPY [--chown=<user>:<group>] <src>... <dest>
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

// Cmd represents a Dockerfile instruction, see https://docs.docker.com/engine/reference/builder/#cmd
type Cmd struct {
	Params  `yaml:"params"`
	RunForm `yaml:"runForm"`
}

// Render returns a string in the form of CMD ["executable","param1","param2"]
func (c Cmd) Render() string {
	if c.RunForm == "" {
		c.RunForm = ExecForm
	}

	if c.RunForm == ExecForm {
		return fmt.Sprintf("CMD %s", c.ExecForm())
	}

	return fmt.Sprintf("CMD %s", c.ShellForm())
}

// Entrypoint represents a Dockerfile instruction, see https://docs.docker.com/engine/reference/builder/#entrypoint
type Entrypoint struct {
	Params  `yaml:"params"`
	RunForm `yaml:"runForm"`
}

// Render returns a string in the form of ENTRYPOINT ["executable", "param1", "param2"]
func (e Entrypoint) Render() string {
	if e.RunForm == "" {
		e.RunForm = ExecForm
	}

	if e.RunForm == ExecForm {
		return fmt.Sprintf("ENTRYPOINT %s", e.ExecForm())
	}

	return fmt.Sprintf("ENTRYPOINT %s", e.ShellForm())
}

// Onbuild represents a Dockerfile instruction, see https://docs.docker.com/engine/reference/builder/#onbuuild
type Onbuild struct {
	Params `yaml:"params"`
}

// Render returns a string in the form of ONBUILD [INSTRUCTION]
func (o Onbuild) Render() string {
	return fmt.Sprintf("ONBUILD %s", o.ShellForm())
}

// HealthCheck represents a Dockerfile instruction, see https://docs.docker.com/engine/reference/builder/#healthcheck
type HealthCheck struct {
	Params `yaml:"params"`
}

// Render returns a string in the form of HEALTHCHECK [OPTIONS] CMD command
func (h HealthCheck) Render() string {
	return fmt.Sprintf("HEALTHCHECK %s", h.ShellForm())
}

// Shell represents a Dockerfile instruction, see https://docs.docker.com/engine/reference/builder/#shell
type Shell struct {
	Params `yaml:"params"`
}

// Render returns a string in the form of SHELL ["executable", "parameters"]
func (s Shell) Render() string {
	return fmt.Sprintf("SHELL %s", s.ExecForm())
}

// Workdir represents a Dockerfile instruction, see https://docs.docker.com/engine/reference/builder/#workdir
type Workdir struct {
	Dir string `yaml:"dir"`
}

// Render returns a string in the form of WORKDIR /path/to/workdir
func (w Workdir) Render() string {
	return fmt.Sprintf("WORKDIR %s", w.Dir)
}

// User represents a Dockerfile instruction, see https://docs.docker.com/engine/reference/builder/#user
type User struct {
	User  string `yaml:user`
	Group string `yaml:group`
}

// Render returns a string in the form of WORKDIR /path/to/workdir
func (u User) Render() string {
	res := "USER"

	if u.User == "" && u.Group == "" {
		return ""
	}

	res = fmt.Sprintf("%s %s", res, u.User)

	if u.Group != "" {
		res = fmt.Sprintf("%s:%s", res, u.Group)
	}

	return res
}

// Params is struct that supports rendering exec and shell form string
type Params []string

func (p Params) mapParams(f func(string) string) []string {
	res := make([]string, len(p))
	for i, v := range p {
		res[i] = f(v)
	}
	return res
}

// ExecForm joins params slice in exec form
func (p Params) ExecForm() string {
	params := p.mapParams(func(s string) string {
		return fmt.Sprintf("\"%s\"", s)
	})

	paramsString := strings.Join(params, ", ")
	execFormString := fmt.Sprintf("[%s]", paramsString)

	return execFormString
}

// ShellForm joins params slice in shell form
func (p Params) ShellForm() string {
	return strings.Join(p, " ")
}
