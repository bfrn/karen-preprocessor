package parse

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	cmdutil "github.com/bfrn/karen-preprocessor/pkg/cmd/util"
	"github.com/bfrn/karen-preprocessor/pkg/preprocessor"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Options is a struct to support parse command
type Options struct {
	InputType string

	InputPath  string
	OutputPath string
	Url        string
	FilePath   string

	args []string
}

// NewOptions returns initialized Options
func NewOptions() *Options {
	return &Options{}
}

// NewCmdParse returns a cobra command for parsing terraform files
func NewCmdParse() *cobra.Command {
	o := NewOptions()
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse a terraform state file to karen intermediate format",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}
	cmd.Flags().StringVarP(&o.InputType, "type", "t", o.InputType, "One of 'plan' or 'state'.")
	cmd.Flags().StringVarP(&o.inputPath, "input", "i", o.inputPath, "relative path to input file.")
	cmd.Flags().StringVarP(&o.outputPath, "output", "o", o.outputPath, "relative path to output location.")
	cmd.Flags().StringVar(&o.url, "url", o.url, "url of the remote repository where the terraform files are located")
	cmd.Flags().StringVar(&o.filePath, "filePath", o.filePath, "relative path under which the terraform files are located in the remote repository")

	cmd.MarkFlagRequired("type")

	return cmd
}

// Complete completes all the required options
func (o *Options) Complete(args []string) error {
	o.args = args
	return nil
}

// Validate validates the provided options
func (o *Options) Validate() error {
	if len(o.args) != 0 {
		return errors.New(fmt.Sprintf("extra arguments: %v", o.args))
	}

	if o.InputType == "" || o.InputType != "plan" && o.InputType != "state" {
		return errors.New(`--type must be 'state' or 'plan'`)
	}

	switch o.InputType {
	case "plan":
		if o.InputPath == "" || o.OutputPath == "" || o.Url == "" || o.FilePath == "" {
			return errors.New("type 'plan' requires inputPath, outputPath, url and filePath")
		}
	case "state":
		if o.InputPath == "" || o.OutputPath == "" {
			return errors.New("type 'state' requires inputPath and outputPath")
		}
	}

	return nil
}

// Run executes parse command
func (o *Options) Run() error {

	log.Debug().Msgf("read file %s", o.InputPath)
	data, err := os.ReadFile(o.InputPath)
	if err != nil {
		return err
	}

	var parsedModel map[string]preprocessor.Node

	switch o.InputType {
	case "plan":
		log.Debug().Msgf("parse plan file url=%s filepath=%s", o.Url, o.FilePath)
		parsedModel, err = preprocessor.ParsePlanFile(data, o.Url, o.FilePath)
		if err != nil {
			return err
		}

	case "state":
		log.Debug().Msgf("parse state file")
		parsedModel, err = preprocessor.ParseStateFile(data)
		if err != nil {
			return err
		}
	}
	output, err := json.Marshal(parsedModel)
	if err != nil {
		return err
	}

	log.Debug().Msgf("write file %s", o.OutputPath)
	err = os.WriteFile(o.OutputPath, output, 0644)
	if err != nil {
		return err
	}
	return nil
}
