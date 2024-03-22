package files

import (
	"filesystem_service/auth/tokens"
	fs "filesystem_service/customFs"
	"filesystem_service/customHttp"
	"fmt"
	"net/http"
	"strings"
)

func GetFileContent(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	if !tokens.IsAuthorized(writer, request) {
		return
	}

	path := fs.GetRoot() + strings.Replace(request.PathValue("path"), "%2F", "/", -1)
	file := fs.NewFile(path)

	_, err := file.IsFile()

	if err != nil {
		writer.Header().Add("Content-Type", "application/json")
		customHttp.WriteError(writer, 400, err)

		return
	}

	extension := file.GetExtension()
	writer.Header().Add("Content-Type", fileFormats[extension])
	fileContent, err := file.GetContent()
	if err != nil {
		customHttp.WriteError(writer, 404, fmt.Errorf(fmt.Sprintf("open %s not found", path)))
		return
	}

	_, _ = writer.Write(fileContent)
}
