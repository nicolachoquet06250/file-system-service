package customHttp

import (
	"encoding/json"
	"filesystem_service/types"
	"net/http"
)

func WriteError(writer http.ResponseWriter, code int, err error) {
	writer.WriteHeader(code)
	response, _ := json.Marshal(types.HttpError{
		Code:    code,
		Message: err.Error(),
	})
	_, _ = writer.Write(response)
}
