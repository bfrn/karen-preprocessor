package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/bfrn/karen-preprocessor/pkg/preprocessor"
)

// Options is a struct to support parse command in wasm
type Options struct {
	InputType string

	data     []byte
	Url      string
	FilePath string

	args []string
}

// NewOptions returns initialized Options
func NewOptions() *Options {
	return &Options{}
}

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("wasmKarenPreprocessor", js.FuncOf(exec))
	<-done
}

func exec(this js.Value, args []js.Value) interface{} {
	o := NewOptions()
	o.InputType = args[0].String()
	array := args[1]
	buf := make([]byte, array.Get("length").Int())
	js.CopyBytesToGo(buf, array)
	o.data = buf
	if o.InputType == "plan" {
		o.Url = args[3].String()
		o.FilePath = args[4].String()
	}

	var err error
	var parsedModel map[string]preprocessor.Node
	switch o.InputType {
	case "plan":
		parsedModel, err = preprocessor.ParsePlanFile(o.data, o.Url, o.FilePath)
		if err != nil {
			return fmt.Sprintln(err.Error())
		}

	case "state":
		parsedModel, err = preprocessor.ParseStateFile(o.data)
		if err != nil {
			return fmt.Sprintln(err.Error())
		}
	default:
		return "Invalid type"
	}
	output, err := json.Marshal(parsedModel)
	if err != nil {
		return fmt.Sprintln(err.Error())
	}

	arrayConstructor := js.Global().Get("Uint8Array")
	dataJS := arrayConstructor.New(len(output))
	js.CopyBytesToJS(dataJS, output)
	// doc := js.Global().Get("document")
	// body := doc.Call("getElementById", "outputHash")
	// body.Set("innerHTML", dataJS)

	return dataJS
}
