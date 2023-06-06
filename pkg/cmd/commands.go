package commands

import (
	"github.com/bfrn/karen-preprocessor/pkg/cmd/parse"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// Options is a struct to support root command
type KarenProcessorOptions struct {
	debug bool
	args  []string
}

// NewOptions returns initialized Options
func NewKarenProcessorOptions() KarenProcessorOptions {
	return KarenProcessorOptions{}
}

// NewKarenProcessorCommand returns a root cobra command for KarenProcessor
func NewKarenProcessorCommand(o KarenProcessorOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "karenprocessor",
		Short: "karenprocessor - parses terraform state files",
		Run:   runHelp,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if o.debug {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			} else {
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}
			return nil
		},
	}

	cmd.PersistentFlags().BoolVarP(&o.debug, "debug", "d", o.debug, "enables debug log level")

	// add sub-commands to root
	cmd.AddCommand(parse.NewCmdParse())

	return cmd
}

// Execute runs the root command
func Execute() error {
	o := NewKarenProcessorOptions()
	// default command
	cmd := NewKarenProcessorCommand(o)
	return cmd.Execute()
}

// runHelp prints the help text of the root cmd
func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
