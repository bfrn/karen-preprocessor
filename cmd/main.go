package main

import (
	commands "github.com/bfrn/karen-preprocessor/pkg/cmd"
)

func main() {
	o := commands.NewKarenProcessorOptions()
	command := commands.NewKarenProcessorCommand(o)
	command.Execute()
}
