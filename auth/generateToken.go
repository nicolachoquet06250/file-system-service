package auth

import "fmt"

func GenerateToken(signature string) (string, error) {
	if signature == "" {
		return "", fmt.Errorf("signature is empty")
	}
	return getNewToken(signature, 60), nil
}
