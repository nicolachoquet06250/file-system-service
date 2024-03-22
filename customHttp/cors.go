package customHttp

import (
	"net/http"
)

func Cors(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(
			"Access-Control-Allow-Origin",
			request.Header.Get("Origin"),
		)
		writer.Header().Set(
			"Access-Control-Allow-Methods",
			"POST,GET,OPTIONS,PUT,PATCH,DELETE",
		)
		writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Access-Control-Allow-Origin, Access-Control-Allow-Methods, Accept, Content-Length, Accept-Encoding, X-CSRF-Token, Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With",
		)

		if request.Method == "OPTIONS" {
			writer.WriteHeader(http.StatusOK)
			return
		}

		handler(writer, request)
	}
}
