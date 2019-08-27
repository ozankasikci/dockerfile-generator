package dockerfilegenerator

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// DockerfileDataYaml is used to decode yaml data. It has a Stages map instead of a slice for that purpose.
type DockerfileDataYaml struct {
	Stages map[string]Stage `yaml:"stages"`
}

type yamlMapStringInterface map[string]interface{}
type yamlMapInterfaceInterface map[interface{}]interface{}

func ensureMapInterfaceInterface(value interface{}) map[interface{}]interface{} {
	v, ok := value.(map[interface{}]interface{})
	if !ok {
		panic(fmt.Sprintf("Yaml contains un expected data, caused by %v", value))
	}

	return v
}

func ensureMapStringInterface(value interface{}) map[string]interface{} {
	v, ok := value.(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("Yaml contains un expected data, caused by %[1]v, type: %[1]T", value))
	}

	return v
}

func ensureMapString(value interface{}) string {
	v, ok := value.(string)
	if !ok {
		panic(fmt.Sprintf("Yaml contains un expected data, caused by %[1]v, type: %[1]T", value))
	}

	return v
}

func convertMapIIToMapSS(mapInterface map[interface{}]interface{}) map[string]string {
	mapString := make(map[string]string)

	for key, value := range mapInterface {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		mapString[strKey] = strValue
	}

	return mapString
}

func convertMapSIToMapSS(mapInterface yamlMapStringInterface) map[string]string {
	mapString := make(map[string]string)

	for key, value := range mapInterface {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		mapString[strKey] = strValue
	}

	return mapString
}

func convertSliceInterfaceToString(s interface{}) ([]string, error) {
	slice, ok := s.([]interface{})
	if !ok {
		return nil, errors.New("Invalid type, can't cast interface{} to []interface{}")
	}

	res := make([]string, len(slice))

	for i, value := range slice {
		res[i] = fmt.Sprintf("%v", value)
	}

	return res, nil
}

func cleanUpFrom(value yamlMapStringInterface) From {
	v := convertMapSIToMapSS(value)
	var from From

	if v["image"] != "" {
		from.Image = v["image"]
	}

	if v["as"] != "" {
		from.As = v["as"]
	}

	return from
}

func cleanUpArg(value yamlMapInterfaceInterface) Arg {
	v := convertMapIIToMapSS(value)
	var arg Arg

	if v["name"] != "" {
		arg.Name = v["name"]
	}

	if v["value"] != "" {
		arg.Value = v["value"]
	}

	if v["test"] == "true" || v["test"] == "yes" {
		arg.Test = true
	}

	if v["envVariable"] == "true" || v["envVariable"] == "yes" {
		arg.EnvVariable = true
	}

	return arg
}

func cleanUpLabel(value yamlMapStringInterface) Label {
	v := convertMapSIToMapSS(value)
	var l Label

	if v["name"] != "" {
		l.Name = v["name"]
	}

	if v["value"] != "" {
		l.Value = v["value"]
	}

	return l
}

func cleanUpVolume(value yamlMapStringInterface) Volume {
	v := convertMapSIToMapSS(value)
	var vlm Volume

	if v["source"] != "" {
		vlm.Source = v["source"]
	}

	if v["destination"] != "" {
		vlm.Destination = v["destination"]
	}

	return vlm
}

func cleanUpRunCommand(value yamlMapInterfaceInterface) RunCommand {
	var r RunCommand
	v := convertMapIIToMapSS(value)

	params, err := convertSliceInterfaceToString(value["params"])
	if err != nil {
		panic("Failed to parse run instruction params!")
	}
	r.Params = params

	r.RunForm = RunCommandDefaultRunForm
	if v["runForm"] == "exec" {
		r.RunForm = ExecForm
	} else if v["runForm"] == "shell" {
		r.RunForm = ShellForm
	}

	return r
}

func cleanUpEnvVariable(value yamlMapStringInterface) EnvVariable {
	v := convertMapSIToMapSS(value)
	var e EnvVariable

	if v["name"] != "" {
		e.Name = v["name"]
	}

	if v["value"] != "" {
		e.Value = v["value"]
	}

	return e
}

func cleanUpCopyCommand(value yamlMapInterfaceInterface) CopyCommand {
	var c CopyCommand
	v := convertMapIIToMapSS(value)

	params, err := convertSliceInterfaceToString(value["sources"])
	if err != nil {
		panic("Failed to parse copy instruction sources!")
	}
	c.Sources = params

	if v["destination"] != "" {
		c.Destination = v["destination"]
	}

	if v["chown"] != "" {
		c.Chown = v["chown"]
	}

	if v["from"] != "" {
		c.From = v["from"]
	}

	return c
}

func cleanUpCmd(value yamlMapInterfaceInterface) Cmd {
	var c Cmd
	v := convertMapIIToMapSS(value)

	params, err := convertSliceInterfaceToString(value["params"])
	if err != nil {
		panic("Failed to parse cmd instruction params!")
	}
	c.Params = params

	c.RunForm = CmdDefaultRunForm
	if v["runForm"] == "exec" {
		c.RunForm = ExecForm
	} else if v["runForm"] == "shell" {
		c.RunForm = ShellForm
	}

	return c
}

func cleanUpEntrypoint(value yamlMapInterfaceInterface) Entrypoint {
	var e Entrypoint
	v := convertMapIIToMapSS(value)

	params, err := convertSliceInterfaceToString(value["params"])
	if err != nil {
		panic("Failed to parse entrypoint instruction params!")
	}
	e.Params = params

	e.RunForm = EntrypointDefaultRunForm
	if v["runForm"] == "exec" {
		e.RunForm = ExecForm
	} else if v["runForm"] == "shell" {
		e.RunForm = ShellForm
	}

	return e
}

func cleanUpOnbuild(value yamlMapInterfaceInterface) Onbuild {
	var o Onbuild

	params, err := convertSliceInterfaceToString(value["params"])
	if err != nil {
		panic("Failed to parse onBuild instruction params!")
	}
	o.Params = params

	return o
}

func cleanUpHealthCheck(value yamlMapInterfaceInterface) HealthCheck {
	var h HealthCheck

	params, err := convertSliceInterfaceToString(value["params"])
	if err != nil {
		panic("Failed to parse healthCheck instruction params!")
	}
	h.Params = params

	return h
}

func cleanUpShell(value yamlMapInterfaceInterface) Shell {
	var s Shell

	params, err := convertSliceInterfaceToString(value["params"])
	if err != nil {
		panic("Failed to parse shell instruction params!")
	}
	s.Params = params

	return s
}

func cleanUpWorkdir(value yamlMapStringInterface) Workdir {
	v := convertMapSIToMapSS(value)
	var w Workdir

	if v["dir"] != "" {
		w.Dir = v["dir"]
	}

	return w
}

func cleanUpUserString(value string) User {
	return User{User: value}
}

func cleanUpUserMap(value yamlMapStringInterface) User {
	v := convertMapSIToMapSS(value)
	var u User

	if v["user"] != "" {
		u.User = v["user"]
	}

	if v["group"] != "" {
		u.Group = v["group"]
	}

	return u
}

func cleanUpInterfaceArray(in []interface{}) []Instruction {
	result := make([]Instruction, len(in))
	for i, v := range in {
		result[i] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpMapSI(in map[string]interface{}) Instruction {
	for key, value := range in {
		switch key {
		case "user":
			v := ensureMapString(value)
			return cleanUpUserString(v)
		}
	}

	panic("Unknown instruction in yaml!")
}

func cleanUpMapII(in map[interface{}]interface{}) Instruction {
	for key, value := range in {
		switch key {
		case "from":
			v := ensureMapStringInterface(value)
			return cleanUpFrom(v)
		case "arg":
			v := ensureMapInterfaceInterface(value)
			return cleanUpArg(v)
		case "label":
			v := ensureMapStringInterface(value)
			return cleanUpLabel(v)
		case "volume":
			v := ensureMapStringInterface(value)
			return cleanUpVolume(v)
		case "run":
			v := ensureMapInterfaceInterface(value)
			return cleanUpRunCommand(v)
		case "envVariable":
			v := ensureMapStringInterface(value)
			return cleanUpEnvVariable(v)
		case "copy":
			v := ensureMapInterfaceInterface(value)
			return cleanUpCopyCommand(v)
		case "cmd":
			v := ensureMapInterfaceInterface(value)
			return cleanUpCmd(v)
		case "entrypoint":
			v := ensureMapInterfaceInterface(value)
			return cleanUpEntrypoint(v)
		case "onbuild":
			v := ensureMapInterfaceInterface(value)
			return cleanUpOnbuild(v)
		case "healthCheck":
			v := ensureMapInterfaceInterface(value)
			return cleanUpHealthCheck(v)
		case "shell":
			v := ensureMapInterfaceInterface(value)
			return cleanUpShell(v)
		case "workdir":
			v := ensureMapStringInterface(value)
			return cleanUpWorkdir(v)
		case "user":
			v := ensureMapStringInterface(value)
			return cleanUpUserMap(v)
		}

	}

	panic("Unknown instruction in yaml!")
}

func cleanUpMapValue(v interface{}) Instruction {
	switch v := v.(type) {
	case map[string]interface{}:
		return cleanUpMapSI(v)
	case map[interface{}]interface{}:
		return cleanUpMapII(v)
	default:
		panic("Invalid instruction type in yaml!")
	}
}

func unmarshallYamlFile(filename string, node *yaml.Node, data *DockerfileDataYaml) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, node)
	if err != nil {
		return fmt.Errorf("Unmarshal: %v", err)
	}

	err = node.Decode(data)
	if err != nil {
		return fmt.Errorf("Unmarshal: %v", err)
	}

	return nil
}

func getStagesOrderFromYamlNode(node *yaml.Node) ([]string, error) {
	var stages []string
	parentMapNode := node.Content[0]

	if parentMapNode.Kind != yaml.MappingNode {
		return nil, errors.New("Yaml should contain a map that contains 'stages' key!")
	}

	stagesKeyNode := parentMapNode.Content[0]
	if stagesKeyNode.Kind != yaml.ScalarNode {
		return nil, errors.New("Yaml should contain a 'stages' key!")
	}

	stagesMapNode := parentMapNode.Content[1]
	if stagesMapNode.Kind != yaml.MappingNode {
		return nil, errors.New("Yaml should contain a 'stages' map that has stage names as keys!")
	}

	for i, stage := range stagesMapNode.Content {
		if i%2 == 0 {
			if stage.Kind != yaml.ScalarNode {
				return nil, errors.New("Yaml should contain stage keys in 'staging' map")
			}
			stages = append(stages, stage.Value)
		} else {
			if stage.Kind != yaml.SequenceNode {
				return nil, errors.New("Yaml should contain stage sequences in 'staging' map")
			}
		}
	}

	return stages, nil
}
