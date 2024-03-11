package main

import (
	"encoding/json"
	fs "filesystem_service/customFs"
	"net/http"
	"strings"
)

func getFileSystem(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	path := "/" + request.PathValue("path")

	d := fs.NewDirectory(path)
	list, err := d.GetFlatContent()

	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			response, _ := json.Marshal(&HttpError{
				Code:    404,
				Message: err.Error(),
			})
			_, _ = writer.Write(response)
			return
		}

		response, _ := json.Marshal(&HttpError{
			Code:    400,
			Message: err.Error(),
		})
		writer.WriteHeader(400)
		_, _ = writer.Write(response)
		return
	}

	response, _ := json.Marshal(list)
	_, _ = writer.Write(response)
}
