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

// Render iterates through the given dockerfile instruction instances and executes the template.
// The output would be a generated Dockerfile.
func (d *DockerFileTemplate) Render(writer io.Writer) error {
	templateString := `
{{ range .Stages -}}
{{ range . }}
{{ .Render }}
{{- end }}
{{ end }}`

	tmpl, err := template.New("dockerfile.template").Parse(templateString)
	if err != nil {
		return err
	}

	err = tmpl.Execute(writer, d.Data)
	if err != nil {
        return err
	}

	return nil
}
