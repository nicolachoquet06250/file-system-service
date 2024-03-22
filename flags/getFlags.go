package flags

import (
	"flag"
	"os"
	"regexp"
)

var port = flag.Int("port", 3000, "Port d'exposition de l'application.")
var host = flag.String("host", "127.0.0.1", "Domaine ou IP de la machine qui expose le sercice.")

var portEnv = flag.String("port-env", "", "Variable d'environement où trouver le port d'exposition de l'application.")
var hostEnv = flag.String("host-env", "", "Variable d'environement où trouver le domaine ou l'IP de la machine qui expose le sercice.")

var generateCredentials = flag.Bool("generate-credentials", false, "Active l'option de génération des credentials pour l'utilisateur.")
var updateCredentials = flag.Bool("update-credentials", false, "Active l'option de modification des credentials pour l'utilisateur.")
var showUserRole = flag.Bool("show-user-role", false, "Active l'option de d'affichage du rôle de l'utilisateur.")

var showRoles = flag.Bool("show-roles", false, "Afficher la liste des rôles disponibles.")
var role = flag.String("role", "readonly", "Rensègne le rôle selectionné.")
var clientId = flag.String("client_id", "", "Rensègne le client_id.")

var isBuiltVersion = func() bool {
	isMatch := true
	for _, _ = range regexp.MustCompile(`(?m)/go-build[0-9]*`).FindAllString(os.Args[0], -1) {
		isMatch = false
	}
	return isMatch
}()

func GetFlags() Flags {
	if !flag.Parsed() {
		flag.Parse()
	}

	return Flags{
		port,
		host,
		portEnv,
		hostEnv,
		generateCredentials,
		updateCredentials,
		showRoles,
		showUserRole,
		role,
		clientId,
	}
}

func IsProd() bool {
	return isBuiltVersion
}
