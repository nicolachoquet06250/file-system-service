package main

import (
	"encoding/json"
	"net/http"
)

type CheckValidity struct {
	IsValid bool `json:"isValid"`
}

func checkValidity(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Content-Type", "application/json")

	response := &CheckValidity{true}

	_json, err := json.Marshal(response)
	if err != nil {
		println(err)
		return
	}
	println(string(_json))

	_, _ = writer.Write(_json)
}
