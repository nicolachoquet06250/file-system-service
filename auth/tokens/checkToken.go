package tokens

import (
	"encoding/json"
	"filesystem_service/auth"
	"filesystem_service/database"
	//"filesystem_service/auth"
	//"filesystem_service/customHttp"
	//"filesystem_service/database"
	"filesystem_service/types"
	"fmt"
	//"fmt"
	"net/http"
	"strings"
	"time"
)

func CheckToken(request *http.Request) (bool, *[]byte, error) {
	accessToken := request.Header.Get("Authorization")
	if accessToken == "" {
		return false, nil, fmt.Errorf("invalid token")
	}
	accessToken = strings.Replace(accessToken, "Bearer ", "", 1)

	db, err := database.Init()
	if err != nil {
		response, _ := json.Marshal(types.HttpError{
			Code:    500,
			Message: err.Error(),
		})
		return false, &response, nil
	}

	rows, err := db.Query(`SELECT id, ip, 
       client_id, access_token, 
       refresh_token, creation_date, 
       expires_in, active 
   FROM tokens
   WHERE active = TRUE AND 
         access_token = ?;`, accessToken)
	if err != nil {
		response, _ := json.Marshal(types.HttpError{
			Code:    500,
			Message: err.Error(),
		})
		return false, &response, nil
	}

	tokens, err := database.ReadRows[auth.Token](rows, func(t *auth.Token) error {
		return rows.Scan(
			&t.Id, &t.Ip,
			&t.ClientId, &t.Token,
			&t.RefreshToken, &t.CreatedAt,
			&t.ExpiresIn, &t.Active,
		)
	})
	if len(tokens) == 0 {
		return false, nil, fmt.Errorf("invalid token")
	}
	token := tokens[0]

	if token.GetCreationDate()+int64(token.ExpiresIn) < time.Now().Unix() {
		if _, err = db.Exec(`UPDATE tokens SET active = FALSE WHERE active = TRUE AND access_token = ?`, accessToken); err != nil {
			response, _ := json.Marshal(types.HttpError{
				Code:    500,
				Message: err.Error(),
			})
			return false, &response, nil
		}

		return false, nil, fmt.Errorf("invalid token")
	}

	return true, nil, nil
}

func IsAuthorized(writer http.ResponseWriter, request *http.Request) bool {
	if _, resp, err := CheckToken(request); err != nil || resp != nil {
		var message types.HttpError
		code := 403
		if err != nil {
			message = types.HttpError{
				Code:    403,
				Message: err.Error(),
			}
		} else if resp != nil {
			var data types.HttpError
			_ = json.Unmarshal(*resp, &data)
			code = data.Code
			message = data
		}
		writer.WriteHeader(code)
		response, _ := json.Marshal(&message)
		_, _ = writer.Write(response)
		return false
	}
	return true
}
