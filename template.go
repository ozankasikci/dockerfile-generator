/*
Package dockerfilegenerator is a Dockerfile generation library. It receives any kind of Dockerfile instructions
and spits out a generated Dockerfile.
*/
package dockerfilegenerator

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"strconv"
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

func getTargetNode(node *yaml.Node, targetField string) (*yaml.Node, error) {
	type part struct{
		kind   string
		val    string
		intVal int
	}

	var parts []part
	var curPart part

	for _, char := range targetField {
		if char == '.' {
            curPart = part{kind: "map"}
			parts = append(parts, curPart)
			continue
		} else if char == '[' {
			curPart = part{kind: "seq"}
			parts = append(parts, curPart)
			continue
		} else if char == ']' {
			if len(parts) == 0 || parts[len(parts)-1].kind != "seq" {
				return nil, fmt.Errorf("invalid target-val %s", targetField)
			}
			intVal, err := strconv.Atoi(parts[len(parts)-1].val)
			if err != nil {
				return nil, fmt.Errorf("invalid target-val %s", targetField)
			}

			parts[len(parts)-1].intVal = intVal
			continue
		}

		if len(parts) == 0 {
			return nil, fmt.Errorf("invalid target-val %s", targetField)
		}

		if len(parts) > 0 {
			parts[len(parts)-1].val += string(char)
		}
	}

    curNode := node.Content[0]

	for _, part := range parts {
		if	part.kind == "map" {
			if curNode.Kind != yaml.MappingNode {
				return nil, fmt.Errorf("Expected a map with key %s", part.val)
			}

			for pos, curContent := range curNode.Content {
				println(curContent.Value)
				if curContent.Value == part.val {
					println(part.val)
					println(pos)
                    curNode = curNode.Content[pos + 1]
                    break
				}
			}
		}
	}

	return curNode, nil
}

func NewDockerFileDataFromYamlField(filename, targetField string) (*DockerfileData, error) {
	d := DockerfileDataYaml{}
	node := yaml.Node{}

	err := unmarshallYamlFile(filename, &node, &d)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %v", err)
	}

	targetNode, err := getTargetNode(&node, targetField)
	if err != nil {
		return nil, fmt.Errorf("Can't decode target val: %v", err)
	}

	stagesInOrder, err := getStagesOrderFromYamlNode(targetNode)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %v", err)
	}

	if err := targetNode.Decode(&d); err != nil {
		return nil, err
	}

	var stages []Stage
	for _, stageName := range stagesInOrder {
		stages = append(stages, d.Stages[stageName])
	}

	return &DockerfileData{Stages: stages}, nil
}

// NewDockerFileDataFromYamlFile reads a file and return a *DockerfileData
func NewDockerFileDataFromYamlFile(filename string) (*DockerfileData, error) {
	d := DockerfileDataYaml{}
	node := yaml.Node{}

	err := unmarshallYamlFile(filename, &node, &d)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %v", err)
	}

	stagesInOrder, err := getStagesOrderFromYamlNode(node.Content[0])
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
