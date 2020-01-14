package dockerfilegenerator

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestGetStagesOrderFromYamlNode(t *testing.T) {
	node := yaml.Node{}

	err := unmarshallYamlFile("./example-input-files/test-input.yaml", &node)
	assert.NoError(t, err)

	stages, err := getStagesOrderFromYamlNode(node.Content[0])
	assert.NoError(t, err)
	assert.Equal(t, stages, []string{"builder", "final"})
}
