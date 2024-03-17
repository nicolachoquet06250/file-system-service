package directories

import (
	"encoding/json"
	fs "filesystem_service/customFs"
	"filesystem_service/types"
	"net/http"
	"strings"
)

func DeleteDirectory(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

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
