package roles

import (
	"filesystem_service/database"
	"fmt"

	"github.com/TwiN/go-color"
)

func GetUserRole(clientId string) {
	if clientId == "" {
		fmt.Printf("Vous devez renseigner le client_id.\n")
		return
	}

	db, err := database.Init()
	defer db.Close()
	if err != nil {
		fmt.Printf("Une erreur est survenue lors de l'initialisation de la base de donnée.\n")
		fmt.Printf(err.Error() + "\n")
		return
	}

	rows, err := db.Query(`SELECT r.id as id, 
       role_name as name,
       group_concat(role_action, '|') as role_actions_ids,
       group_concat(role_action_name, '|') as role_actions_names,
       r.active as is_active FROM credentials
                  INNER JOIN roles r ON r.id = credentials.role
                  INNER JOIN main.roles_link_role_actions rlra ON r.id = rlra.role
                  INNER JOIN main.role_actions ra on ra.id = rlra.role_action
		WHERE client_id = ? AND r.active = TRUE`, clientId)
	if err != nil {
		fmt.Printf("Une erreur est survenue.\n")
		fmt.Printf(err.Error() + "\n")
		return
	}

	roles, _ := database.ReadRows[Role](rows, func(r *Role) error {
		return rows.Scan(&r.Id, &r.Name, &r.RoleActionsIds, &r.RoleActionsNames, &r.IsActive)
	})

	if len(roles) == 0 {
		fmt.Printf("Aucun utilisateur actif à été trouvé avec ce client_id.\n")
		return
	}

	for _, role := range roles {
		fmt.Printf(" => %v \n", color.InUnderline(role.GetName()))
		for _, roleAction := range role.GetRoleActions() {
			fmt.Printf("\t-> %v \n", color.OverGray(color.InBold(roleAction.GetName())))
		}
	}
}
