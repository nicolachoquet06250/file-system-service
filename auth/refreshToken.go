package auth

import (
	"encoding/json"
	"filesystem_service/arrays"
	"filesystem_service/customHttp"
	"filesystem_service/data"
	"filesystem_service/types"
	"fmt"
	"net/http"
	"time"
)

func RefreshToken(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Content-Type", "application/json")
	db, err := data.InitDatabase()
	defer db.Close()
	if err != nil {
		writer.WriteHeader(500)
		response, _ := json.Marshal(types.HttpError{
			Code:    500,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}
	refreshTokenHeader := request.Header.Get("Refresh-Token")
	signatureTokenHeader := request.Header.Get("Signature-Token")

	if refreshTokenHeader == "" {
		writer.WriteHeader(400)
		response, _ := json.Marshal(types.HttpError{
			Code:    400,
			Message: "Invalid refresh token",
		})
		_, _ = writer.Write(response)
		return
	}

	ip, err := customHttp.GetUserIp(request)
	if err != nil {
		writer.WriteHeader(500)
		response, _ := json.Marshal(types.HttpError{
			Code:    500,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	tokens, err := db.Query(fmt.Sprintf(
		"SELECT * FROM tokens WHERE IP=\"%v\" AND active=TRUE AND signature=\"%v\"",
		ip, signatureTokenHeader,
	))
	if err != nil {
		writer.WriteHeader(500)
		response, _ := json.Marshal(types.HttpError{
			Code:    500,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	formattedTokens, err := data.ReadRows[Token](tokens, func(t *Token) error {
		return tokens.Scan(&t.Id, &t.Ip, &t.Token, &t.Signature, &t.Type, &t.Active, &t.CreatedAt)
	})

	_, err = json.Marshal(&formattedTokens)

	if err != nil {
		println(err.Error())
	}

	if len(formattedTokens) == 0 {
		writer.WriteHeader(400)
		response, _ := json.Marshal(types.HttpError{
			Code:    400,
			Message: "Invalid couple signature and refresh token",
		})
		_, _ = writer.Write(response)
		return
	}

	refreshTokenObj := arrays.Filter[Token](formattedTokens, func(token Token) bool {
		return token.Type == "refresh"
	})[0]

	if refreshTokenHeader != refreshTokenObj.Token {
		writer.WriteHeader(400)
		response, _ := json.Marshal(types.HttpError{
			Code:    400,
			Message: "Invalid refresh token",
		})
		_, _ = writer.Write(response)
		return
	}

	token, err := GenerateToken(signatureTokenHeader + "." + refreshTokenHeader)
	if err != nil {
		writer.WriteHeader(500)
		response, _ := json.Marshal(types.HttpError{
			Code:    500,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}
	refreshToken, err := GenerateToken(token)
	if err != nil {
		writer.WriteHeader(500)
		response, _ := json.Marshal(types.HttpError{
			Code:    500,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	_, err = db.Exec(`UPDATE tokens SET active=FALSE WHERE IP=? AND active=TRUE`, ip)
	if err != nil {
		writer.WriteHeader(500)
		response, _ := json.Marshal(types.HttpError{
			Code:    500,
			Message: err.Error(),
		})
		_, _ = writer.Write(response)
		return
	}

	_, err = db.Exec(
		`INSERT INTO tokens (IP, token, signature, type, created_at) VALUES (?, ?, ?, ?, ?)`,
		ip, token, signatureTokenHeader, "classic", time.Now().Unix(),
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

	_, err = db.Exec(
		`INSERT INTO tokens (IP, token, signature, type, created_at) VALUES (?, ?, ?, ?, ?)`,
		ip, refreshToken, signatureTokenHeader, "refresh", time.Now().Unix(),
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

	response, _ := json.Marshal(&GetTokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now().Unix(),
		// 1h (en ms)
		ExpiresIn: 3600000,
	})
	_, _ = writer.Write(response)
}
