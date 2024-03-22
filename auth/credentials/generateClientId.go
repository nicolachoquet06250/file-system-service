package credentials

import "fmt"

func GenerateClientId() string {
	return fmt.Sprintf("fs_service_i@%v", CreateNewToken(AvailableCharacters, 15))
}
