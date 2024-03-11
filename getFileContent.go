package main

import (
	"encoding/json"
	fs "filesystem_service/customFs"
	"fmt"
	"net/http"
)

func getFileContent(writer http.ResponseWriter, request *http.Request) {
	path := "/" + request.PathValue("path")
	file := fs.NewFile(path)

	_, err := file.IsFile()

	if err != nil {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(400)
		response, _ := json.Marshal(&HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)

		return
	}

	extension := file.GetExtension()
	writer.Header().Add("Content-Type", fileFormats[extension])
	fileContent, err := file.GetContent()
	if err != nil {
		response, _ := json.Marshal(&HttpError{
			Code:    404,
			Message: fmt.Sprintf("open %s not found", path),
		})
		_, _ = writer.Write(response)
		return
	}

	_, _ = writer.Write(fileContent)
}
