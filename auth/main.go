package auth

type GetTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    int64  `json:"created_at"`
	ExpiresIn    int    `json:"expires_in"`
}

type Signature struct {
	Id        int    `json:"id"`
	Signature string `json:"signature"`
	Active    bool   `json:"active"`
}

type Token struct {
	Id        int `json:"id"`
	Ip        string
	Token     string `json:"token"`
	Signature string `json:"signature"`
	Type      string `json:"type"`
	Active    bool   `json:"active"`
	CreatedAt int    `json:"created_at"`
}
