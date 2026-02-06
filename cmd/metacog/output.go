package main

import (
	"encoding/json"
)

type OutputError struct {
	Message    string `json:"error"`
	Code       int    `json:"code"`
	Suggestion string `json:"suggestion"`
}

func FormatOutput(asJSON bool, output string, err *OutputError) string {
	if !asJSON {
		if err != nil {
			return err.Message
		}
		return output
	}

	if err != nil {
		data, _ := json.Marshal(err)
		return string(data)
	}

	data, _ := json.Marshal(map[string]string{"output": output})
	return string(data)
}
