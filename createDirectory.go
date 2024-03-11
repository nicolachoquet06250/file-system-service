package main

import (
	"encoding/json"
	fs "filesystem_service/customFs"
	"net/http"
)

func createDirectory(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	var body fs.Directory

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		writer.WriteHeader(400)
		response, _ := json.Marshal(&HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	d := fs.NewDirectory(fs.BuildDirectoryCompletePath(body))

	if _, err := d.Create(); err != nil {
		writer.WriteHeader(400)
		response, _ := json.Marshal(&HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	response, _ := json.Marshal(&body)
	_, _ = writer.Write(response)
}
