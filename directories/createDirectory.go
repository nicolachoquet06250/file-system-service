package directories

import (
	"encoding/json"
	"filesystem_service/auth/tokens"
	fs "filesystem_service/customFs"
	"filesystem_service/customHttp"
	"net/http"
)

func CreateDirectory(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Content-Type", "application/json")

	if !tokens.IsAuthorized(writer, request) {
		return
	}

	var body fs.Directory

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		customHttp.WriteError(writer, 400, err)
		return
	}

	d := fs.NewDirectory(fs.BuildDirectoryCompletePath(body))

	if _, err := d.Create(); err != nil {
		customHttp.WriteError(writer, 400, err)
		return
	}

	response, _ := json.Marshal(&body)
	_, _ = writer.Write(response)
}
