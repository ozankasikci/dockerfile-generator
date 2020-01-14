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

// Tries to return a *yaml.Node based on the given targetField
func getTargetNode(node *yaml.Node, targetField string) (*yaml.Node, error) {
	type part struct {
		kind   string
		val    string
		intVal int
	}

	last := func(parts []part) *part {
		return &(parts[len(parts)-1])
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
			if len(parts) == 0 || last(parts).kind != "seq" {
				return nil, fmt.Errorf("invalid target-val %s", targetField)
			}

			intVal, err := strconv.Atoi(last(parts).val)
			if err != nil {
				return nil, fmt.Errorf("invalid target-val %s", targetField)
			}

			last(parts).intVal = intVal
			continue
		}

		if len(parts) == 0 {
			return nil, fmt.Errorf("invalid target-val %s", targetField)
		}

		if len(parts) > 0 {
			last(parts).val += string(char)
		}
	}

	curNode := node.Content[0]

	for _, part := range parts {
		if part.kind == "map" {
			if curNode.Kind != yaml.MappingNode {
				return nil, fmt.Errorf("Expected a map with key %s", part.val)
			}

			for pos, curContent := range curNode.Content {
				if curContent.Value == part.val {
					curNode = curNode.Content[pos+1]
					break
				}
			}
		} else if part.kind == "seq" {
			if curNode.Kind != yaml.SequenceNode {
				return nil, fmt.Errorf("Expected a seq with key %s", part.val)
			}

			curNode = curNode.Content[part.intVal]
		}
	}

	return curNode, nil
}

func getStagesDataFromNode(node *yaml.Node) ([]Stage, error) {
	var data DockerfileDataYaml

	stagesInOrder, err := getStagesOrderFromYamlNode(node)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %v", err)
	}

	if err := node.Decode(&data); err != nil {
		return nil, err
	}

	var stages []Stage
	for _, stageName := range stagesInOrder {
		stages = append(stages, data.Stages[stageName])
	}

	return stages, nil
}

// NewDockerFileDataFromYamlField reads a YAML file and tries to extract Dockerfile data
// from the specified targetField option, examples:
// --target-field ".dev.dockerfileConfig"
// --target-field ".serverConfigs[0].docker.server"
func NewDockerFileDataFromYamlField(filename, targetField string) (*DockerfileData, error) {
	node := yaml.Node{}

	err := unmarshallYamlFile(filename, &node)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %v", err)
	}

	targetNode, err := getTargetNode(&node, targetField)
	if err != nil {
		return nil, fmt.Errorf("Can't decode target val: %v", err)
	}

	stages, err := getStagesDataFromNode(targetNode)
	if err != nil {
		return nil, fmt.Errorf("Can't extract stages from node: %v", err)
	}

	return &DockerfileData{Stages: stages}, nil
}

// NewDockerFileDataFromYamlFile reads a file and return a *DockerfileData
func NewDockerFileDataFromYamlFile(filename string) (*DockerfileData, error) {
	node := yaml.Node{}

	err := unmarshallYamlFile(filename, &node)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %v", err)
	}

	// passing node.Content[0] because the file is expected to store solely the dockerfile config
	stages, err := getStagesDataFromNode(node.Content[0])
	if err != nil {
		return nil, fmt.Errorf("Can't extract stages from node: %v", err)
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
