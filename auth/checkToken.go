package auth

import (
	"filesystem_service/customHttp"
	"filesystem_service/data"
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
