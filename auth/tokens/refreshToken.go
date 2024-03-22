package tokens

import (
	"encoding/json"
	"filesystem_service/auth"
	"filesystem_service/customHttp"
	"filesystem_service/database"
	"fmt"
	"net/http"
)

func RefreshToken(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	clientId := request.URL.Query().Get("client_id")
	clientSecret := request.URL.Query().Get("client_secret")

	if clientId == "" || clientSecret == "" {
		customHttp.WriteError(writer, 403, fmt.Errorf("you must input client_id and client_secret for refresh your token."))
		return
	}

	db, err := database.Init()
	defer db.Close()
	if err != nil {
		customHttp.WriteError(writer, 500, err)
		return
	}

	refreshToken := request.Header.Get("Refresh-Token")
	if refreshToken == "" {
		var body *struct {
			RefreshToken string `json:"refresh_token"`
		}
		_ = json.NewDecoder(request.Body).Decode(body)

		if body.RefreshToken == "" {
			customHttp.WriteError(writer, 400, fmt.Errorf("Invalid refresh token"))
			return
		}
		refreshToken = body.RefreshToken
	}

	ip, err := customHttp.GetUserIp(request)
	if err != nil {
		customHttp.WriteError(writer, 500, err)
		return
	}

	activeTokens, _ := db.Query(`SELECT client_id FROM tokens WHERE active = TRUE AND ip = ?`, ip)
	_tokens, _ := database.ReadRows[auth.Token](activeTokens, func(t *auth.Token) error {
		return activeTokens.Scan(&t.ClientId)
	})

	tokenCredentials, _ := _tokens[0].GetCredentials()
	if clientSecret != tokenCredentials.ClientSecret {
		customHttp.WriteError(writer, 403, fmt.Errorf("Invalid credentials."))
		return
	}

	rows, err := db.Query(`SELECT id, ip, 
       client_id, access_token, 
       refresh_token, creation_date, 
       expires_in, active 
	FROM tokens 
	WHERE active = TRUE AND 
	      ip = ? AND 
	      refresh_token = ?;`, ip, refreshToken)
	if err != nil {
		customHttp.WriteError(writer, 500, err)
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
		customHttp.WriteError(writer, 500, err)
		return
	}

	if len(tokens) == 0 {
		customHttp.WriteError(writer, 400, fmt.Errorf("Invalid couple signature and refresh token"))
		return
	}

	if refreshToken != tokens[0].RefreshToken {
		customHttp.WriteError(writer, 400, fmt.Errorf("Invalid refresh token"))
		return
	}

	newToken, err := GenerateToken()
	if err != nil {
		customHttp.WriteError(writer, 500, err)
		return
	}
	newRefreshToken, err := GenerateToken()
	if err != nil {
		customHttp.WriteError(writer, 500, err)
		return
	}

	_, err = db.Exec(`UPDATE tokens SET active = FALSE 
              WHERE ip = ? AND 
                    active = TRUE`, ip)
	if err != nil {
		customHttp.WriteError(writer, 500, err)
		return
	}

	result, err := db.Exec(
		`INSERT INTO tokens (ip, client_id, access_token, refresh_token, active) VALUES (?, ?, ?, ?, TRUE)`,
		ip, tokens[0].ClientId, newToken, newRefreshToken,
	)
	if err != nil {
		customHttp.WriteError(writer, 500, err)
		return
	}

	lastInsertedId, err := result.LastInsertId()

	rows, err = db.Query(`SELECT id, ip, 
       client_id, access_token, 
       refresh_token, creation_date, 
       expires_in, active 
	FROM tokens 
	WHERE id = ?`, lastInsertedId)
	if err != nil {
		customHttp.WriteError(writer, 500, err)
		return
	}
	tokens, err = database.ReadRows[auth.Token](rows, func(t *auth.Token) error {
		return rows.Scan(
			&t.Id, &t.Ip, &t.ClientId,
			&t.Token, &t.RefreshToken,
			&t.CreatedAt, &t.ExpiresIn,
			&t.Active,
		)
	})
	if err != nil {
		customHttp.WriteError(writer, 500, err)
		return
	}
	lastCreatedToken := tokens[0]

	response, _ := json.Marshal(lastCreatedToken.ToJsonToken())
	_, _ = writer.Write(response)
}
