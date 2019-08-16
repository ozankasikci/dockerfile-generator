package dockerfile_generator

import (
	"errors"
	"fmt"
)

type DockerfileDataYaml struct {
	Stages map[string]Stage `yaml:stages`
}

func ensureMapInterfaceInterface(value interface{}) map[interface{}]interface{} {
	v, ok := value.(map[interface{}]interface{})
	if !ok {
		panic("Invalid value for 'from'")
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

func cleanUpFrom(value map[interface{}]interface{}) From {
	v := convertMapInterfaceToString(value)
	var from From

	if v["image"] != "" {
		from.Image = v["image"]
	}

	if v["as"] != "" {
		from.As = v["as"]
	}

	return from
}

func cleanUpArg(value map[interface{}]interface{}) Arg {
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

func cleanUpLabel(value map[interface{}]interface{}) Label {
	v := convertMapInterfaceToString(value)
	var l Label

	if v["name"] != "" {
		l.Name = v["name"]
	}

	if v["value"] != "" {
		l.Value = v["value"]
	}

	return l
}

func cleanUpVolume(value map[interface{}]interface{}) Volume {
	v := convertMapInterfaceToString(value)
	var vlm Volume

	if v["source"] != "" {
		vlm.Source = v["source"]
	}

	if v["destination"] != "" {
		vlm.Destination = v["destination"]
	}

	return vlm
}

func cleanUpRunCommand(value map[interface{}]interface{}) RunCommand {
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

func cleanUpEnvVariable(value map[interface{}]interface{}) EnvVariable {
	v := convertMapInterfaceToString(value)
	var e EnvVariable

	if v["name"] != "" {
		e.Name = v["name"]
	}

	if v["value"] != "" {
		e.Value = v["value"]
	}

	return e
}

func cleanUpCopyCommand(value map[interface{}]interface{}) CopyCommand {
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

func cleanUpCmd(value map[interface{}]interface{}) Cmd {
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

func cleanUpEntrypoint(value map[interface{}]interface{}) Entrypoint {
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

func cleanUpOnbuild(value map[interface{}]interface{}) Onbuild {
	var o Onbuild

	params, err := convertSliceInterfaceToString(value["params"])
	if err != nil {
		panic("Failed to parse onBuild instruction params!")
	}
	o.Params = params

	return o
}

func cleanUpHealthCheck(value map[interface{}]interface{}) HealthCheck {
	var h HealthCheck

	params, err := convertSliceInterfaceToString(value["params"])
	if err != nil {
		panic("Failed to parse healthCheck instruction params!")
	}
	h.Params = params

	return h
}

func cleanUpShell(value map[interface{}]interface{}) Shell {
	var s Shell

	params, err := convertSliceInterfaceToString(value["params"])
	if err != nil {
		panic("Failed to parse shell instruction params!")
	}
	s.Params = params

	return s
}

func cleanUpWorkdir(value map[interface{}]interface{}) Workdir {
	v := convertMapInterfaceToString(value)
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
		v := ensureMapInterfaceInterface(value)
		switch key {
		case "from":
			return cleanUpFrom(v)
		case "arg":
			return cleanUpArg(v)
		case "label":
			return cleanUpLabel(v)
		case "volume":
			return cleanUpVolume(v)
		case "run":
			return cleanUpRunCommand(v)
		case "envVariable":
			return cleanUpEnvVariable(v)
		case "copy":
			return cleanUpCopyCommand(v)
		case "cmd":
			return cleanUpCmd(v)
		case "entrypoint":
			return cleanUpEntrypoint(v)
		case "onbuild":
			return cleanUpOnbuild(v)
		case "healthCheck":
			return cleanUpHealthCheck(v)
		case "shell":
			return cleanUpShell(v)
		case "workdir":
			return cleanUpWorkdir(v)
		}

	}
    panic("Unknown instruction in yaml!")
}

func cleanUpMapValue(v interface{}) Instruction {
	switch v := v.(type) {
	case map[interface {}]interface {}:
		return cleanUpInterfaceMap(v)
	default:
        panic("Invalid instruction type in yaml!")
	}
}

