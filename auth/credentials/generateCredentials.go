package credentials

import (
	"filesystem_service/auth/roles"
	"filesystem_service/database"
	"fmt"
)

func GenerateCredentials(role string) {
	db, err := database.Init()
	defer db.Close()
	if err != nil {
		fmt.Printf("Une erreur est survenue lors de l'initialisation de la base de donnée.\n")
		fmt.Printf(err.Error() + "\n")
		return
	}

	roleId := roles.GetRoleIdFromName(db, role)

	if roleId != -1 {
		fmt.Printf("Generation du client_id ...\n")
		fmt.Printf("Generation du client_secret ...\n")
		clientId := GenerateClientId()
		clientSecret := GenerateClientSecret()

		if _, err = db.Exec(`UPDATE credentials 
			SET active = FALSE 
		WHERE active = TRUE;`); err != nil {
			fmt.Printf("\nUne erreur est survenue lors de la création du token de signature.\n")
			fmt.Printf(err.Error() + "\n")
			return
		}

		if _, err = db.Exec(`INSERT INTO credentials (
			client_id, 
			client_secret, 
			role, 
			active
		) VALUES (?, ?, ?, TRUE);`, clientId, clientSecret, roleId); err != nil {
			fmt.Printf("\nUne erreur est survenue lors de la création du token de signature.\n")
			fmt.Printf(err.Error() + "\n")
			return
		}

		fmt.Printf("\nVous devrez les saisir dans le système d'exploitation web.\n")
		fmt.Printf("\n=> client_id: %v\n=> client_secret: %v\n", clientId, clientSecret)
	}
}
