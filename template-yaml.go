package dockerfile_generator

import (
	"errors"
	"fmt"
)

type DockerfileDataYaml struct {
	Stages map[string]Stage `yaml:stages`
}

type YAMLMapStringInterface map[string]interface{}
type YAMLMapInterfaceInterface map[interface{}]interface{}

func ensureMapInterfaceInterface(value interface{}) (map[interface{}]interface{}) {
	v, ok := value.(map[interface{}]interface{})
	if !ok {
		panic("no")
		//return nil, errors.New(fmt.Sprintf("Expected map[interface]interface found, %T", value))
	}

	return v
}

func ensureMapStringInterface(value interface{}) (map[string]interface{}) {
	v, ok := value.(map[string]interface{})
	if !ok {
		panic("no")
		//return nil, errors.New(fmt.Sprintf("Expected map[string]interface found, %T", value))
	}

	return v
}

func convertMapInterfaceToString(mapInterface map[interface{}]interface{}) map[string]string {
	mapString := make(map[string]string)

	for key, value := range mapInterface {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		mapString[strKey] = strValue
	}

	return mapString
}

func convertMapStringInterfaceToString(mapInterface map[string]interface{}) map[string]string {
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

func cleanUpFrom(value YAMLMapStringInterface) From {
	v := convertMapStringInterfaceToString(value)
	var from From

	if v["image"] != "" {
		from.Image = v["image"]
	}

	if v["as"] != "" {
		from.As = v["as"]
	}

	return from
}

func cleanUpArg(value YAMLMapInterfaceInterface) Arg {
	v := convertMapInterfaceToString(value)
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

func cleanUpLabel(value YAMLMapStringInterface) Label {
	v := convertMapStringInterfaceToString(value)
	var l Label

	if v["name"] != "" {
		l.Name = v["name"]
	}

	if v["value"] != "" {
		l.Value = v["value"]
	}

	return l
}

func cleanUpVolume(value YAMLMapStringInterface) Volume {
	v := convertMapStringInterfaceToString(value)
	var vlm Volume

	if v["source"] != "" {
		vlm.Source = v["source"]
	}

	if v["destination"] != "" {
		vlm.Destination = v["destination"]
	}

	return vlm
}

func cleanUpRunCommand(value YAMLMapInterfaceInterface) RunCommand {
	var r RunCommand
	v := convertMapInterfaceToString(value)

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

func cleanUpEnvVariable(value YAMLMapStringInterface) EnvVariable {
	v := convertMapStringInterfaceToString(value)
	var e EnvVariable

	if v["name"] != "" {
		e.Name = v["name"]
	}

	if v["value"] != "" {
		e.Value = v["value"]
	}

	return e
}

func cleanUpCopyCommand(value YAMLMapInterfaceInterface) CopyCommand {
	var c CopyCommand
	v := convertMapInterfaceToString(value)

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

func cleanUpCmd(value YAMLMapInterfaceInterface) Cmd {
	var c Cmd
	v := convertMapInterfaceToString(value)

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

func cleanUpEntrypoint(value YAMLMapInterfaceInterface) Entrypoint {
	var e Entrypoint
	v := convertMapInterfaceToString(value)

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

func cleanUpOnbuild(value YAMLMapInterfaceInterface) Onbuild {
	var o Onbuild

	params, err := convertSliceInterfaceToString(value["params"])
	if err != nil {
		panic("Failed to parse onBuild instruction params!")
	}
	o.Params = params

	return o
}

func cleanUpHealthCheck(value YAMLMapInterfaceInterface) HealthCheck {
	var h HealthCheck

	params, err := convertSliceInterfaceToString(value["params"])
	if err != nil {
		panic("Failed to parse healthCheck instruction params!")
	}
	h.Params = params

	return h
}

func cleanUpShell(value YAMLMapInterfaceInterface) Shell {
	var s Shell

	params, err := convertSliceInterfaceToString(value["params"])
	if err != nil {
		panic("Failed to parse shell instruction params!")
	}
	s.Params = params

	return s
}

func cleanUpWorkdir(value YAMLMapStringInterface) Workdir {
	v := convertMapStringInterfaceToString(value)
	var w Workdir

	if v["dir"] != "" {
		w.Dir = v["dir"]
	}

	return w
}

func cleanUpInterfaceArray(in []interface{}) []Instruction {
	result := make([]Instruction, len(in))
	for i, v := range in {
		result[i] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpInterfaceMap(in map[interface{}]interface{}) Instruction {
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
		}

	}
    //panic("Unknown instruction in yaml!")
    return Workdir{Dir: "test"}
}

func cleanUpMapValue(v interface{}) Instruction {
	switch v := v.(type) {
	case map[interface {}]interface {}:
		return cleanUpInterfaceMap(v)
	default:
        panic("Invalid instruction type in yaml!")
	}
}

