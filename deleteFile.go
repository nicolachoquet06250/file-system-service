package main

import (
	"encoding/json"
	fs "filesystem_service/customFs"
	"net/http"
)

func deleteFile(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	path := "/" + request.PathValue("path")

	f := fs.NewFile(path)

	if _, err := f.Delete(); err != nil {
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
