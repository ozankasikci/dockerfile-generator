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

//func cleanUpRunCommand(value map[interface{}]interface{}) RunCommand {
//	v := convertMapInterfaceToString(value)
//	var r RunCommand
//
//	if v["params"] != "" {
//		r.Source = v["source"]
//	}
//
//	if v["destination"] != "" {
//		r.Destination = v["destination"]
//	}
//
//	return r
//}

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
			v := ensureMapInterfaceInterface(value)
			return cleanUpFrom(v)
		case "arg":
			v := ensureMapInterfaceInterface(value)
			return cleanUpArg(v)
		case "label":
			v := ensureMapInterfaceInterface(value)
			return cleanUpLabel(v)
		case "volume":
			v := ensureMapInterfaceInterface(value)
			return cleanUpVolume(v)
		//case "run":
		//	v := ensureMapInterfaceInterface(value)
		//	return cleanUpRunCommand(v)
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

