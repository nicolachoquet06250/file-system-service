package main

import (
	"encoding/json"
	fs "filesystem_service/customFs"
	"io"
	"net/http"
)

func renameFile(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	var renamedFile fs.File

	path := "/" + request.PathValue("path")

	if err := json.NewDecoder(request.Body).Decode(&renamedFile); err != nil {
		writer.WriteHeader(400)
		response, _ := json.Marshal(&HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	f := fs.NewFile(path)

	if _, err := f.Rename(fs.NewFile(fs.BuildFileCompletePath(renamedFile))); err != nil {
		writer.WriteHeader(400)
		response, _ := json.Marshal(&HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	response, _ := json.Marshal(&ResponseStatus{
		Status: "success",
	})
	_, _ = writer.Write(response)
}

func updateFileContent(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	path := "/" + request.PathValue("path")

	var body []byte

	f := fs.NewFile(path)

	text, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(400)
		response, _ := json.Marshal(&HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}
	body = text

	if _, err = f.SetContent(body); err != nil {
		writer.WriteHeader(400)
		response, _ := json.Marshal(&HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	response, _ := json.Marshal(&ResponseStatus{
		Status: "success",
	})
	_, _ = writer.Write(response)
}
