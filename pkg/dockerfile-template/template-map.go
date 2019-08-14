package dockerfile_template

import (
	"io"
	"text/template"
)

type DockerFileMapTemplate struct {
	Data *DockerfileDataSlice
}

func NewDockerFileMapTemplate(data *DockerfileDataSlice) *DockerFileMapTemplate {
    return &DockerFileMapTemplate{ Data: data }
}

func (d *DockerFileMapTemplate) Render(writer io.Writer, templateFilePath string) error {
	funcMap := template.FuncMap{}

	tmpl, err := template.New("dockerfile.map.template").Funcs(funcMap).ParseFiles(templateFilePath)
	if err != nil {
		return err
	}

	err = tmpl.Execute(writer, d.Data)
	if err != nil {
        return err
	}

	return nil
}
