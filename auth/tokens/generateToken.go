package tokens

import (
	"filesystem_service/auth/credentials"
)

func GenerateToken() (string, error) {
	return credentials.CreateNewToken(credentials.AvailableCharacters, 60), nil
}
