/*
Package dockerfile-generator is a Dockerfile generation library. It receives any kind of Dockerfile instructions
and spits out a generated Dockerfile.
 */
package dockerfile_generator
import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"text/template"
)

type DockerfileTemplate struct {
	Data *DockerfileData
}

func NewDockerfileTemplate(data *DockerfileData) *DockerfileTemplate {
    return &DockerfileTemplate{ Data: data }
}

func NewDockerFileDataFromYamlFile(filename string) (*DockerfileData, error) {
	d := &DockerfileDataYaml{}
	node := &yaml.Node{}

	err := unmarshallYamlFile(filename, node, d)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unmarshal: %v", err))
	}

	stagesInOrder, err := getStagesOrderFromYamlNode(node)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unmarshal: %v", err))
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
