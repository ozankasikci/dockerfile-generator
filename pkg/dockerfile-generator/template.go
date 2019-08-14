package dockerfile_generator

import (
	"io"
	"os"
	"text/template"
)

type DockerFileTemplate struct {
	Data *DockerfileData
}

func NewDockerFileTemplate(data *DockerfileData) *DockerFileTemplate {
    return &DockerFileTemplate{ Data: data }
}

func (d *DockerFileTemplate) Render(writer io.Writer) error {
	funcMap := template.FuncMap{}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}


	templateFilePath := wd + "/pkg/template/dockerfile.template"

	tmpl, err := template.New("dockerfile.template").Funcs(funcMap).ParseFiles(templateFilePath)
	if err != nil {
		return err
	}

	err = tmpl.Execute(writer, d.Data)
	if err != nil {
        return err
	}

	return nil
}
