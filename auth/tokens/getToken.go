package tokens

import (
	"encoding/base64"
	"encoding/json"
	"filesystem_service/auth"
	"filesystem_service/customHttp"
	"filesystem_service/database"
	"filesystem_service/types"
	"fmt"
	"net/http"
	"strings"
)

func GetToken(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	db, err := database.Init()
	defer db.Close()
	if err != nil {
		customHttp.WriteError(writer, 500, err)
		return
	}

	queryString := request.URL.Query()
	clientId := queryString.Get("client_id")
	clientSecret := queryString.Get("client_secret")

	if clientId == "" || clientSecret == "" {
		basicToken := request.Header.Get("Authorization")
		if strings.Contains(basicToken, "Basic ") {
			basicToken = strings.Replace(basicToken, "Basic ", "", 1)
			_basicToken, err := base64.StdEncoding.DecodeString(basicToken)
			if err != nil {
				customHttp.WriteError(writer, 403, fmt.Errorf("You must input your client_id and client_secret credentials."))
				return
			}
			basicToken = string(_basicToken)

			splitToken := strings.Split(basicToken, ":")
			clientId = splitToken[0]
			clientSecret = splitToken[1]
		} else {
			customHttp.WriteError(writer, 403, fmt.Errorf("You must input your client_id and client_secret credentials."))
			return
		}
	}

	foundCredentials, err := db.Query(`SELECT 
		client_id, client_secret,
		role_name as role,
		creation_date, updated_date,
		expires_in, c.active as active
		FROM credentials c INNER JOIN roles r ON r.id = c.role
		WHERE c.active = TRUE AND 
		      client_id = ? AND 
		      client_secret = ?;`,
		clientId, clientSecret)
	if err != nil {
		customHttp.WriteError(writer, 500, err)
		return
	}

	credentials, err := database.ReadRows[auth.Credential](foundCredentials, func(t *auth.Credential) error {
		return foundCredentials.Scan(
			&t.ClientId, &t.ClientSecret,
			&t.Role, &t.CreationDate,
			&t.UpdatedDate, &t.ExpiresIn,
			&t.IsActive,
		)
	})
	if err != nil {
		customHttp.WriteError(writer, 500, err)
		return
	}

	if len(credentials) == 0 {
		customHttp.WriteError(writer, 403, fmt.Errorf("You are not identified"))
		return
	}

	token, err := GenerateToken()
	if err != nil {
		customHttp.WriteError(writer, 500, err)
		return
	}
	refreshToken, err := GenerateToken()
	if err != nil {
		customHttp.WriteError(writer, 500, err)
		return
	}

	ip, err := customHttp.GetUserIp(request)
	if err != nil {
		customHttp.WriteError(writer, 500, err)
		return
	}

	_, err = db.Exec(`UPDATE tokens SET active = FALSE WHERE ip = ? AND active = TRUE`, ip)
	if err != nil {
		writer.WriteHeader(500)
		response, _ := json.Marshal(types.HttpError{
			Code:    500,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	result, err := db.Exec(
		`INSERT INTO tokens (ip, client_id, access_token, refresh_token, active) VALUES (?, ?, ?, ?, TRUE)`,
		ip, clientId, token, refreshToken,
	)
	if err != nil {
		writer.WriteHeader(500)
		response, _ := json.Marshal(types.HttpError{
			Code:    500,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	lastInsertedId, err := result.LastInsertId()

	rows, err := db.Query(`SELECT id, ip, 
       client_id, access_token, 
       refresh_token, creation_date, 
       expires_in, active 
	FROM tokens 
	WHERE id = ?;`, lastInsertedId)
	if err != nil {
		writer.WriteHeader(500)
		response, _ := json.Marshal(types.HttpError{
			Code:    500,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}
	tokens, err := database.ReadRows[auth.Token](rows, func(t *auth.Token) error {
		return rows.Scan(
			&t.Id, &t.Ip, &t.ClientId,
			&t.Token, &t.RefreshToken,
			&t.CreatedAt, &t.ExpiresIn,
			&t.Active,
		)
	})
	if err != nil {
		writer.WriteHeader(500)
		response, _ := json.Marshal(types.HttpError{
			Code:    500,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}
	if len(tokens) == 0 {
		writer.WriteHeader(403)
		response, _ := json.Marshal(types.HttpError{
			Code:    403,
			Message: "Invalid token.",
		})
		_, _ = writer.Write(response)
		return
	}
	lastCreatedToken := tokens[0]

	response, _ := json.Marshal(lastCreatedToken.ToJsonToken())
	_, _ = writer.Write(response)
}
