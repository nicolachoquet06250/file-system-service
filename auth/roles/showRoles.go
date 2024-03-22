package roles

import (
	"database/sql"
	"filesystem_service/database"
	"fmt"

	"github.com/TwiN/go-color"
)

func GetRoleIdFromName(db *sql.DB, role string) int {
	rows, _ := db.Query(`SELECT id, role_name as name FROM roles WHERE role_name = ?`, role)
	findRoles, _ := database.ReadRows[Role](rows, func(r *Role) error {
		return rows.Scan(&r.Id, &r.Name)
	})
	if len(findRoles) == 0 {
		return -1
	}
	return findRoles[0].GetId()
}

func ShowRoles() {
	db, err := database.Init()
	defer db.Close()
	if err != nil {
		fmt.Printf("Une erreur est survenue lors de l'initialisation de la base de donnée.\n")
		fmt.Printf(err.Error() + "\n")
		return
	}

	rows, err := db.Query(`SELECT role as id,
       role_name as name,
       group_concat(role_action, '|') as role_actions_ids,
       group_concat(role_action_name, '|') as role_actions_names,
       active as is_active
    FROM roles_link_role_actions
        INNER JOIN roles
                 ON roles_link_role_actions.role = roles.id
        INNER JOIN role_actions
                 ON roles_link_role_actions.role_action = role_actions.id
    GROUP BY role_name`)
	if err != nil {
		fmt.Printf("Une erreur est survenue lors de la recherche des rôles.\n")
		fmt.Printf(err.Error() + "\n")
		return
	}

	roles, _ := database.ReadRows[Role](rows, func(r *Role) error {
		return rows.Scan(&r.Id, &r.Name, &r.RoleActionsIds, &r.RoleActionsNames, &r.IsActive)
	})

	fmt.Printf("Liste des roles disponibles:\n")
	if len(roles) == 0 {
		fmt.Printf(" == La liste de roles est vide, veuillez en créer ==\n")
		return
	}
	for _, role := range roles {
		fmt.Printf(" => %v \n", color.InUnderline(role.GetName()))
		for _, roleAction := range role.GetRoleActions() {
			fmt.Printf("\t-> %v \n", color.OverGray(color.InBold(roleAction.GetName())))
		}
	}
}
