package files

import (
	"encoding/json"
	"filesystem_service/auth"
	fs "filesystem_service/customFs"
	"filesystem_service/types"
	"net/http"
	"strings"
)

func DeleteFile(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Content-Type", "application/json")

	if !auth.IsAuthorized(writer, request) {
		return
	}

	path := fs.GetRoot() + strings.Replace(request.PathValue("path"), "%2F", "/", -1)

	f := fs.NewFile(path)

	if _, err := f.Delete(); err != nil {
		writer.WriteHeader(400)
		response, _ := json.Marshal(&types.HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	response, _ := json.Marshal(&types.ResponseStatus{
		Status: "success",
	})
	_, _ = writer.Write(response)
}
