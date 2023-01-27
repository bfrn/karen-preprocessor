package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bfrn/karen-preprocessor/pkg/preprocessor"
)

func main() {
	callParameter := os.Args[1]
	if callParameter == "--plan-file" {
		inputPath := os.Args[2]
		outputPath := os.Args[3]
		url := os.Args[4]
		filePath := os.Args[5]

		data, err := os.ReadFile(inputPath)
		if err != nil {
			fmt.Printf("Could not read the input file: %s", err.Error())
			return
		}
		parsedModel, err := preprocessor.ParsePlanFile(data, url, filePath)
		if err != nil {
			fmt.Printf("Could not parse the given file: %s", err.Error())
			return
		}

		output, err := json.Marshal(parsedModel)
		if err != nil {
			fmt.Printf("Could not unmarshal the given file: %s", err.Error())
			return
		}
		err = os.WriteFile(outputPath, output, 0644)
		if err != nil {
			fmt.Printf("Could not unmarshal the given file: %s", err.Error())
			return
		}
	} else if callParameter == "--state-file" {
		inputPath := os.Args[2]
		outputPath := os.Args[3]
		data, err := os.ReadFile(inputPath)
		if err != nil {
			fmt.Printf("Could not read the input file: %s", err.Error())
			return
		}
		parsedModel, err := preprocessor.ParseStateFile(data)
		if err != nil {
			fmt.Printf("Could not parse the given file: %s", err.Error())
			return
		}
		output, err := json.Marshal(parsedModel)
		if err != nil {
			fmt.Printf("Could not unmarshal the given file: %s", err.Error())
			return
		}
		err = os.WriteFile(outputPath, output, 0644)
		if err != nil {
			fmt.Printf("Could not unmarshal the given file: %s", err.Error())
			return
		}
	} else {
		fmt.Printf("Invalid call parameter: %s", callParameter)
	}

}
