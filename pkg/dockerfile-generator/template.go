package dockerfile_generator

import (
	"io"
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
	templateFilePath := "pkg/template/dockerfile.template"

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
