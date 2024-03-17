package auth

import (
	"filesystem_service/arrays"
	"filesystem_service/data"
	"fmt"
	"math/rand"
	"strings"

	"github.com/atotto/clipboard"
)

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

var availableCharacters = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ0123456789+-*/=!&@#'()|\\{}[]_~"

func getNewToken(availableCharacters string, charNumber int) string {
	return strings.Join(
		arrays.Map(
			arrays.Generate[string](charNumber),
			func(t string) string {
				return strings.Split(availableCharacters, "")[randInt(0, len(availableCharacters)-1)]
			},
		),
		"",
	)
}

func generateSignatureToken() string {
	return getNewToken(availableCharacters, 60)
}

func GenerateSignatureTokenAction() {
	fmt.Printf("Generation du token de signature ...\n")
	token := generateSignatureToken()
	err := clipboard.WriteAll(token)
	db, err := data.InitDatabase()
	defer db.Close()
	if err != nil {
		fmt.Printf("Une erreur est survenue lors de l'initialisation de la base de donnée.\n")
		fmt.Printf(err.Error() + "\n")
		return
	}

	_, err = db.Exec(`UPDATE signatures SET active=FALSE where active=TRUE`, token)
	if err != nil {
		fmt.Printf("Une erreur est survenue lors de la création du token de signature.\n")
		fmt.Printf(err.Error() + "\n")
		return
	}

	_, err = db.Exec(`INSERT INTO signatures (signature) VALUES (?)`, token)
	if err != nil {
		fmt.Printf("Une erreur est survenue lors de la création du token de signature.\n")
		fmt.Printf(err.Error() + "\n")
		return
	}

	if err != nil {
		fmt.Printf("Vous devrez le saisir dans le système d'exploitation web.\n")
		fmt.Printf("Nous n'avons pas pu le mettre dans votre press-papier mais le voici ci-dessous:\n")
	} else {
		fmt.Printf("Nous l'avons mis dans votre presse-papier.\n")
		fmt.Printf("Vous devrez le saisir dans le système d'exploitation web.\n")
	}
	fmt.Printf("\n=> %v\n", token)
}
