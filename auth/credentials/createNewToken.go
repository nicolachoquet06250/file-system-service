package credentials

import (
	"filesystem_service/arrays"
	"filesystem_service/integer"
	"strings"
)

func CreateNewToken(availableCharacters string, charNumber int) string {
	return strings.Join(
		arrays.Map[string, string](
			arrays.Generate[string](charNumber),
			func(t string) string {
				return strings.Split(availableCharacters, "")[integer.RandomBetween(0, len(availableCharacters)-1)]
			},
		),
		"",
	)
}
