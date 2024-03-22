package files

import (
	"encoding/json"
	"filesystem_service/auth/tokens"
	fs "filesystem_service/customFs"
	"filesystem_service/customHttp"
	"net/http"
)

func getMultipartKeys(request *http.Request) (file fs.File, content string, err error) {
	if _err := request.ParseMultipartForm(0); _err != nil {
		if __err := json.NewDecoder(request.Body).Decode(&file); __err != nil {
			err = _err
			return
		}
	} else {
		content = request.FormValue("content")

		if _err = json.Unmarshal([]byte(request.FormValue("file")), &file); _err != nil {
			err = _err
			return
		}
	}

	err = nil
	return
}

func CreateFile(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Content-Type", "application/json")

	if !tokens.IsAuthorized(writer, request) {
		return
	}

	file, content, err := getMultipartKeys(request)
	if err != nil {
		writer.Header().Add("Content-Type", "application/json")
		customHttp.WriteError(writer, 400, err)
		return
	}

	f := fs.NewFile(fs.BuildFileCompletePath(file))

	if _, err = f.Create(); err != nil {
		customHttp.WriteError(writer, 400, err)
		return
	}

	if _, err = f.SetContent([]byte(content)); err != nil {
		customHttp.WriteError(writer, 400, err)
		return
	}

	response, _ := json.Marshal(&file)
	_, _ = writer.Write(response)
}
