package credentials

import (
	"filesystem_service/auth/roles"
	"filesystem_service/database"
	"fmt"
)

func UpdateCredentials(clientId string, newRole string) {
	if clientId == "" {
		fmt.Printf("Vous devez renseigner le client_id.\n")
		return
	}

	if newRole == "" {
		fmt.Printf("Vous devez renseigner votre nouveau rôle.\n")
		return
	}

	db, err := database.Init()
	defer db.Close()
	if err != nil {
		fmt.Printf("Une erreur est survenue lors de l'initialisation de la base de donnée.\n")
		fmt.Printf(err.Error() + "\n")
		return
	}

	newRoleId := roles.GetRoleIdFromName(db, newRole)

	if _, err = db.Exec(
		`UPDATE credentials SET role = ? WHERE client_id = ?;`,
		newRoleId, clientId,
	); err != nil {
		fmt.Printf("Une erreur est survenue lors de laa modification de votre rôle.\n")
		fmt.Printf(err.Error() + "\n")
		return
	}

	fmt.Printf("Votre rôle à bien été modifié.\n")
}
