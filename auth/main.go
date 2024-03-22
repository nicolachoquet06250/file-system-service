package auth

import (
	"filesystem_service/database"
	"strings"
	"time"
)

type Credential struct {
	ClientId     string
	ClientSecret string
	Role         string
	CreationDate string
	UpdatedDate  *string
	ExpiresIn    int
	IsActive     bool
}

func (c *Credential) GetCreationDate() time.Time {
	creationDate, _ := time.Parse(time.DateTime, c.CreationDate)
	return creationDate
}

func (c *Credential) GetUpdatedDate() time.Time {
	if creationDate, err := time.Parse(time.DateTime, c.CreationDate); err != nil {
		return c.GetCreationDate()
	} else {
		return creationDate
	}
}

type JsonToken struct {
	ref          *Token
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    int64  `json:"created_at"`
	ExpiresIn    int    `json:"expires_in"`
}

func (t *JsonToken) ToToken() *Token {
	return t.ref
}

type Token struct {
	Id           int
	Ip           string
	ClientId     string
	credentials  *Credential
	Token        string
	RefreshToken string
	CreatedAt    string
	UpdatedAt    int64
	ExpiresIn    int
	Active       bool
}

func (t *Token) GetCreationDate() int64 {
	format := time.RFC3339
	if strings.Contains(t.CreatedAt, " UTC") {
		format = time.DateTime
	}
	creationDate, _ := time.Parse(format, t.CreatedAt)
	return creationDate.Unix()
}

func (t *Token) GetCredentials() (*Credential, error) {
	if t.ClientId != "" && t.credentials == nil {
		db, err := database.Init()
		if err != nil {
			return nil, err
		}

		rows, err := db.Query(`SELECT client_id, client_secret,
		   role_name as role,
		   creation_date, updated_date,
		   expires_in, c.active as active
		FROM credentials c INNER JOIN roles r ON r.id = c.role
		WHERE c.active = TRUE AND 
		      client_id = ?;`, t.ClientId)
		if err != nil {
			return nil, err
		}

		credentials, err := database.ReadRows[Credential](rows, func(c *Credential) error {
			return rows.Scan(&c.ClientId, &c.ClientSecret, &c.Role, &c.CreationDate, &c.UpdatedDate, &c.ExpiresIn, &c.IsActive)
		})
		if err != nil {
			return nil, err
		}

		t.credentials = &credentials[0]
	}

	return t.credentials, nil
}

func (t *Token) ToJsonToken() *JsonToken {
	return &JsonToken{
		ref:          t,
		Token:        t.Token,
		RefreshToken: t.RefreshToken,
		CreatedAt:    t.GetCreationDate(),
		ExpiresIn:    t.ExpiresIn,
	}
}
