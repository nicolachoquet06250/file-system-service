package main

import (
	"encoding/json"
	fs "filesystem_service/customFs"
	"net/http"
)

func renameDirectory(writer http.ResponseWriter, request *http.Request) {
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

	path := "/" + request.PathValue("path")

	d := fs.NewDirectory(path)
	nd := fs.NewDirectory(fs.BuildDirectoryCompletePath(body))

	if _, err := d.Rename(nd); err != nil {
		writer.WriteHeader(400)
		response, _ := json.Marshal(&HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

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

	if fi.Mode().IsDir() {
		newPath := body.Path
		if newPath != "/" {
			newPath += "/"
		}
		newPath += body.Name

		err = os.Rename(path, newPath)
		if err != nil {
			writer.WriteHeader(400)
			response, _ := json.Marshal(&HttpError{
				Code:    400,
				Message: err.Error(),
			})
			_, _ = writer.Write(response)
			return
		}
	}*/

	response, _ := json.Marshal(&body)

	_, _ = writer.Write(response)
}
