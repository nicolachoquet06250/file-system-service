package files

import (
	"encoding/json"
	"filesystem_service/auth"
	fs "filesystem_service/customFs"
	"filesystem_service/types"
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

	if _, err := auth.CheckToken(request); err != nil {
		writer.WriteHeader(403)
		response, _ := json.Marshal(&types.HttpError{
			Code:    403,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	file, content, err := getMultipartKeys(request)
	if err != nil {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(400)
		response, _ := json.Marshal(&types.HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	f := fs.NewFile(fs.BuildFileCompletePath(file))

	if _, err = f.Create(); err != nil {
		writer.WriteHeader(400)
		response, _ := json.Marshal(&types.HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	if _, err = f.SetContent([]byte(content)); err != nil {
		writer.WriteHeader(400)
		response, _ := json.Marshal(&types.HttpError{
			Code:    400,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	response, _ := json.Marshal(&file)
	_, _ = writer.Write(response)
}
