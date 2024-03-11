package main

import (
	"encoding/json"
	fs "filesystem_service/customFs"
	"net/http"
)

func deleteDirectory(writer http.ResponseWriter, request *http.Request) {
	isForce := request.URL.Query().Has("force")
	writer.Header().Add("Content-Type", "application/json")

	path := "/" + request.PathValue("path")

	d := fs.NewDirectory(path)

	/*fi, err := os.Stat(path)
	if err != nil {
		writer.WriteHeader(400)
		response, _ := json.Marshal(&HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	if fi.Mode().IsDir() {*/
	var err error
	if isForce {
		_, err = d.DeepDelete()
	} else {
		_, err = d.Delete()
	}

	if err != nil {
		writer.WriteHeader(400)
		response, _ := json.Marshal(&HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}
	//}

	response, _ := json.Marshal(&ResponseStatus{
		Status: "success",
	})

	_, _ = writer.Write(response)
}
