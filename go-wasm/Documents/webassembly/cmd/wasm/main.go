package main

import (
	"fmt"
	"encoding/json"
	"syscall/js"
	"errors"
)

func prettyJson(input string) (string, error) {
    var raw any
    if err := json.Unmarshal([]byte(input), &raw); err != nil {
        return "", err
    }
    pretty, err := json.MarshalIndent(raw, "", "  ")
    if err != nil {
        return "", err
    }
    return string(pretty), nil
}

func jsonWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return errors.New("Invalid number of arguments")
		}
		inputJSON := args[0].String()
		fmt.Printf("input %s\n", inputJSON)
		pretty, err := prettyJson(inputJSON)
		if err != nil {
			errStr := fmt.Sprintf("unable to convert %s\n", err)
			return errors.New(errStr)
		}

		jsDoc := js.Global().Get("document")
		if !jsDoc.Truthy() { return errors.New("get object doc") }
		jsonOutputTextArea := jsDoc.Call("getElementById", "jsonoutput")
		if !jsonOutputTextArea.Truthy() { return errors.New("getElementByid") }
		jsonOutputTextArea.Set("value", pretty)
		return nil

	})
	return jsonFunc
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("formatJSON", jsonWrapper())
	nullSource := make(chan bool)
	<-nullSource
}
