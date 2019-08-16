package dockerfile_generator

import (
	"fmt"
)

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
	v := convertMapInterfaceToString(value)
	var r RunCommand

	if v["params"] != "" {
		r.Params = []string{"sdf", "adf"}
	}

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
	v := convertMapInterfaceToString(value)
	var c CopyCommand

	if v["sources"] != "" {
		c.Sources = []string{"source1"}
	}

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
	v := convertMapInterfaceToString(value)
	var c Cmd

	if v["params"] != "" {
		c.Params = []string{"sdf", "adf"}
	}

	c.RunForm = CmdDefaultRunForm
	if v["runForm"] == "exec" {
		c.RunForm = ExecForm
	} else if v["runForm"] == "shell" {
		c.RunForm = ShellForm
	}

	return c
}

func cleanUpEntrypoint(value map[interface{}]interface{}) Entrypoint {
	v := convertMapInterfaceToString(value)
	var e Entrypoint

	if v["params"] != "" {
		e.Params = []string{"sdf", "adf"}
	}

	e.RunForm = EntrypointDefaultRunForm
	if v["runForm"] == "exec" {
		e.RunForm = ExecForm
	} else if v["runForm"] == "shell" {
		e.RunForm = ShellForm
	}

	return e
}

func cleanUpOnbuild(value map[interface{}]interface{}) Onbuild {
	v := convertMapInterfaceToString(value)
	var o Onbuild

	if v["params"] != "" {
		o.Params = []string{"sdf", "adf"}
	}

	return o
}

func cleanUpHealthCheck(value map[interface{}]interface{}) HealthCheck {
	v := convertMapInterfaceToString(value)
	var h HealthCheck

	if v["params"] != "" {
		h.Params = []string{"sdf", "adf"}
	}

	return h
}

func cleanUpShell(value map[interface{}]interface{}) Shell {
	v := convertMapInterfaceToString(value)
	var s Shell

	if v["params"] != "" {
		s.Params = []string{"sdf", "adf"}
	}

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

