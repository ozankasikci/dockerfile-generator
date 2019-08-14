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

func (d *DockerFileTemplate) Render(writer io.Writer, templateFilePath string) error {
	funcMap := template.FuncMap{}

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
