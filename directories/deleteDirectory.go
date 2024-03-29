package directories

import (
	"encoding/json"
	"filesystem_service/auth/tokens"
	fs "filesystem_service/customFs"
	"filesystem_service/customHttp"
	"filesystem_service/types"
	"net/http"
	"strings"
)

func DeleteDirectory(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	if !tokens.IsAuthorized(writer, request) {
		return
	}

	isForce := request.URL.Query().Has("force")

	path := fs.GetRoot() + strings.Replace(request.PathValue("path"), "%2F", "/", -1)

	d := fs.NewDirectory(path)

	var err error
	if isForce {
		_, err = d.DeepDelete()
	} else {
		_, err = d.Delete()
	}

	if err != nil {
		customHttp.WriteError(writer, 400, err)
		return
	}

	response, _ := json.Marshal(&types.ResponseStatus{
		Status: "success",
	})

	_, _ = writer.Write(response)
}
