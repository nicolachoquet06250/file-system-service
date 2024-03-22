package credentials

import "fmt"

func GenerateClientSecret() string {
	return fmt.Sprintf("fs_service_s@%v", CreateNewToken(AvailableCharacters, 15))
}
