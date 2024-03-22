package files

import (
	"encoding/json"
	"filesystem_service/auth/tokens"
	fs "filesystem_service/customFs"
	"filesystem_service/customHttp"
	"filesystem_service/types"
	"io"
	"net/http"
	"strings"
)

func RenameFile(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Content-Type", "application/json")

	if !tokens.IsAuthorized(writer, request) {
		return
	}

	var renamedFile fs.File

	path := fs.GetRoot() + strings.Replace(request.PathValue("path"), "%2F", "/", -1)

	if err := json.NewDecoder(request.Body).Decode(&renamedFile); err != nil {
		customHttp.WriteError(writer, 400, err)
		return
	}

	f := fs.NewFile(path)

	if _, err := f.Rename(fs.NewFile(fs.BuildFileCompletePath(renamedFile))); err != nil {
		customHttp.WriteError(writer, 400, err)
		return
	}

	response, _ := json.Marshal(&types.ResponseStatus{
		Status: "success",
	})
	_, _ = writer.Write(response)
}

func UpdateFileContent(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Content-Type", "application/json")

	if !tokens.IsAuthorized(writer, request) {
		return
	}

	path := fs.GetRoot() + strings.Replace(request.PathValue("path"), "%2F", "/", -1)

	var body []byte

	f := fs.NewFile(path)

	text, err := io.ReadAll(request.Body)
	if err != nil {
		customHttp.WriteError(writer, 400, err)
		return
	}
	body = text

	if _, err = f.SetContent(body); err != nil {
		customHttp.WriteError(writer, 400, err)
		return
	}

	response, _ := json.Marshal(&types.ResponseStatus{
		Status: "success",
	})
	_, _ = writer.Write(response)
}
