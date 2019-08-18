/*
Package dockerfilegenerator is a Dockerfile generation library. It receives any kind of Dockerfile instructions
and spits out a generated Dockerfile.
*/
package dockerfilegenerator

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"text/template"
)

// DockerfileTemplate defines the template struct that generates Dockerfile output
type DockerfileTemplate struct {
	Data *DockerfileData
}

// NewDockerfileTemplate return a new NewDockerfileTemplate instance
func NewDockerfileTemplate(data *DockerfileData) *DockerfileTemplate {
	return &DockerfileTemplate{Data: data}
}

// NewDockerFileDataFromYamlFile reads a file and return a *DockerfileData
func NewDockerFileDataFromYamlFile(filename string) (*DockerfileData, error) {
	d := &DockerfileDataYaml{}
	node := &yaml.Node{}

	err := unmarshallYamlFile(filename, node, d)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %v", err)
	}

	stagesInOrder, err := getStagesOrderFromYamlNode(node)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %v", err)
	}

	var stages []Stage
	for _, stageName := range stagesInOrder {
		stages = append(stages, d.Stages[stageName])
	}

	return &DockerfileData{Stages: stages}, nil
}

// Render iterates through the given dockerfile instruction instances and executes the template.
// The output would be a generated Dockerfile.
func (d *DockerfileTemplate) Render(writer io.Writer) error {
	templateString := "{{- range .Stages -}}" +
		"{{- range $i, $instruction := . }}" +
		"{{- if gt $i 0 }}\n{{ end }}" +
		"{{ $instruction.Render }}\n" +
		"{{- end }}\n\n" +
		"{{ end }}"

	tmpl, err := template.New("dockerfile.template").Parse(templateString)
	if err != nil {
		return err
	}

	err = tmpl.Execute(writer, d.Data)
	if err != nil {
		return err
	}

	return nil
}
