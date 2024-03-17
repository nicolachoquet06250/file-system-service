package main

import (
	"encoding/json"
	"filesystem_service/types"
	"net/http"
)

func CheckValidity(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Content-Type", "application/json")

	response := &types.CheckValidity{true}

	_json, err := json.Marshal(response)
	if err != nil {
		println(err)
		return
	}
	println(string(_json))

	_, _ = writer.Write(_json)
}
