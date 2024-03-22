package directories

import (
	"encoding/json"
	"filesystem_service/auth/tokens"
	fs "filesystem_service/customFs"
	"filesystem_service/customHttp"
	"net/http"
	"strings"
)

func GetFileSystem(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Content-Type", "application/json")

	if !tokens.IsAuthorized(writer, request) {
		return
	}

	path := fs.GetRoot() + strings.Replace(request.PathValue("path"), "%2F", "/", -1)

	d := fs.NewDirectory(path)
	list, err := d.GetFlatContent()

	if err != nil {
		code := 400
		if strings.Contains(err.Error(), "no such file or directory") {
			code = 404
		}

		customHttp.WriteError(writer, code, err)
		return
	}

	response, _ := json.Marshal(list)
	_, _ = writer.Write(response)
}
