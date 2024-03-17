package auth

import (
	"encoding/json"
	"filesystem_service/customHttp"
	"filesystem_service/data"
	"filesystem_service/types"
	"fmt"
	"net/http"
	"strings"
)

func CheckToken(request *http.Request) (bool, error) {
	signature := request.Header.Get("Signature-Token")
	token := strings.Replace(request.Header.Get("Authorization"), "Bearer ", "", 1)

	if token == "" || signature == "" {
		return false, fmt.Errorf("you are not authorized")
	}

	ip, err := customHttp.GetUserIp(request)
	if err != nil {
		return false, err
	}

	db, err := data.InitDatabase()
	defer db.Close()
	if err != nil {
		return false, err
	}

	results, err := db.Query(fmt.Sprintf(
		"SELECT * FROM tokens WHERE IP=\"%v\" AND active=TRUE AND type=\"classic\" AND token=\"%s\" AND signature=\"%s\"",
		ip, token, signature,
	))
	if err != nil {
		return false, err
	}

	tokens, err := data.ReadRows[Token](results, func(t *Token) error {
		return results.Scan(&t.Id, &t.Ip, &t.Token, &t.Signature, &t.Type, &t.Active, &t.CreatedAt)
	})
	if err != nil {
		return false, err
	}

	if len(tokens) == 0 {
		return false, fmt.Errorf("invalid access token")
	}

	return true, nil
}

func IsAuthorized(writer http.ResponseWriter, request *http.Request) bool {
	if _, err := CheckToken(request); err != nil {
		writer.WriteHeader(403)
		response, _ := json.Marshal(&types.HttpError{
			Code:    403,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return false
	}
	return true
}
