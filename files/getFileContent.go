package files

import (
	"encoding/json"
	"filesystem_service/auth"
	fs "filesystem_service/customFs"
	"filesystem_service/types"
	"fmt"
	"net/http"
	"strings"
)

func GetFileContent(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	if _, err := auth.CheckToken(request); err != nil {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(403)
		response, _ := json.Marshal(&types.HttpError{
			Code:    403,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	path := fs.GetRoot() + strings.Replace(request.PathValue("path"), "%2F", "/", -1)
	file := fs.NewFile(path)

	_, err := file.IsFile()

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

	extension := file.GetExtension()
	writer.Header().Add("Content-Type", fileFormats[extension])
	fileContent, err := file.GetContent()
	if err != nil {
		response, _ := json.Marshal(&types.HttpError{
			Code:    404,
			Message: fmt.Sprintf("open %s not found", path),
		})
		_, _ = writer.Write(response)
		return
	}

	_, _ = writer.Write(fileContent)
}
