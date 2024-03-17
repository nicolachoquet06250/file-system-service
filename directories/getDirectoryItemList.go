package directories

import (
	"encoding/json"
	"filesystem_service/auth"
	fs "filesystem_service/customFs"
	"filesystem_service/types"
	"net/http"
	"strings"
)

func GetFileSystem(writer http.ResponseWriter, request *http.Request) {
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

	path := fs.ROOT + strings.Replace(request.PathValue("path"), "%2F", "/", -1)

	d := fs.NewDirectory(path)
	list, err := d.GetFlatContent()

	if err != nil {
		code := 400
		if strings.Contains(err.Error(), "no such file or directory") {
			code = 404
		}

		writer.WriteHeader(code)
		response, _ := json.Marshal(&types.HttpError{
			Code:    code,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	response, _ := json.Marshal(list)
	_, _ = writer.Write(response)
}
