package cmd

import (
	dfg "github.com/ozankasikci/dockerfile-generator"
	"github.com/spf13/cobra"
	"io"
	"os"
)

const (
	// YAMLFileInput specifies that the input channel will be a yaml file, this is the default
	YAMLFileInput = "yaml-file"
)

type cmdGenerateConfig struct {
	input     string
	output    string
	inputType string
	stdout    bool
}

// NewCmdGenerate generates a command that is responsible for generating a Dockerfile output
func NewCmdGenerate() *cobra.Command {
	cfg := &cmdGenerateConfig{}

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generates a Dockerfile based on input",
		RunE: func(cmd *cobra.Command, args []string) error {
			switch cfg.inputType {
			case YAMLFileInput:
			default:
				if err := generateFromYAMLFile(cfg); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&cfg.input, "input", "i", "", "Input path")
	cmd.PersistentFlags().StringVarP(&cfg.output, "out", "o", "", "Output file path")
	cmd.PersistentFlags().BoolVar(&cfg.stdout, "stdout", false, "When true, output will be redirected to stdout")
	cmd.PersistentFlags().StringVarP(&cfg.inputType, "type", "t", "", "Input type (yaml-file)")

	return cmd
}

func generateFromYAMLFile(cfg *cmdGenerateConfig) error {
	var outputTarget io.Writer
	data, err := dfg.NewDockerFileDataFromYamlFile(cfg.input)
	if err != nil {
		return err
	}

	tmpl := dfg.NewDockerfileTemplate(data)

	if cfg.stdout {
		outputTarget = os.Stdout
	} else {
		file, err := os.Create(cfg.output)
		outputTarget = file
		if err != nil {
			return err
		}
		defer file.Close()
	}

	err = tmpl.Render(outputTarget)
	if err != nil {
		return err
	}

	return nil
}
