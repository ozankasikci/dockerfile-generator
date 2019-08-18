package dockerfilegenerator

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestGetStagesOrderFromYamlNode(t *testing.T) {
	d := &DockerfileDataYaml{}
	node := &yaml.Node{}

	err := unmarshallYamlFile("./example-input-files/test-input.yaml", node, d)
	assert.NoError(t, err)

	stages, err := getStagesOrderFromYamlNode(node)
	assert.NoError(t, err)
	assert.Equal(t, stages, []string{"builder", "final"})
}
