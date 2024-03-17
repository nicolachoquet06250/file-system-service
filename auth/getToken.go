package auth

import (
	"encoding/json"
	"filesystem_service/customHttp"
	"filesystem_service/data"
	"filesystem_service/types"
	"net/http"
	"time"
)

func GetToken(writer http.ResponseWriter, request *http.Request) {
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
	signatureToken := request.Header.Get("Signature-Token")

	if signatureToken == "" {
		writer.WriteHeader(400)
		response, _ := json.Marshal(types.HttpError{
			Code:    400,
			Message: "You must input your signature token",
		})
		_, _ = writer.Write(response)
		return
	}

	foundSignatures, err := db.Query(
		`SELECT * FROM signatures WHERE active=TRUE AND signature=?`,
		signatureToken,
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

	signatures, err := data.ReadRows[Signature](foundSignatures, func(t *Signature) error {
		return foundSignatures.Scan(&t.Id, &t.Signature, &t.Active)
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

	if len(signatures) == 0 {
		writer.WriteHeader(403)
		response, _ := json.Marshal(types.HttpError{
			Code:    403,
			Message: "You are not identified",
		})
		_, _ = writer.Write(response)
		return
	}

	signature := signatures[0].Signature
	token, err := GenerateToken(signature)
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
		ip, token, signature, "classic", time.Now().Unix(),
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
		ip, refreshToken, signature, "refresh", time.Now().Unix(),
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

	tokenResponse := &GetTokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now().Unix(),
		// 1h (en ms)
		ExpiresIn: 3600000,
	}

	response, _ := json.Marshal(tokenResponse)
	_, _ = writer.Write(response)
}
