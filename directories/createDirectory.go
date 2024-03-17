package directories

import (
	"encoding/json"
	"filesystem_service/auth"
	fs "filesystem_service/customFs"
	"filesystem_service/types"
	"net/http"
)

func CreateDirectory(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Content-Type", "application/json")

	if _, err := auth.CheckToken(request); err != nil {
		writer.WriteHeader(403)
		response, _ := json.Marshal(&types.HttpError{
			Code:    403,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	var body fs.Directory

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		writer.WriteHeader(400)
		response, _ := json.Marshal(&types.HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	d := fs.NewDirectory(fs.BuildDirectoryCompletePath(body))

	if _, err := d.Create(); err != nil {
		writer.WriteHeader(400)
		response, _ := json.Marshal(&types.HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	response, _ := json.Marshal(&body)
	_, _ = writer.Write(response)
}