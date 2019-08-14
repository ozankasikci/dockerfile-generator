package dockerfile

type DockerfileData struct {
	Stages []Stage
}

type Arg struct {
	Name  string
	Value string
	Test  bool
}

type Stage struct {
	Args         []Arg
	From         string
	As           string
	RunCommands  []RunCommand
	Workdir      string
	EnvVariables []EnvVariable
	CopyCommands []CopyCommand
	Cmd          Cmd
	Expose       string
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
	Command string
}
