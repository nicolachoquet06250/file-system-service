package flags

import "flag"

var port = flag.Int("port", 3000, "Port d'exposition de l'application.")
var host = flag.String("host", "127.0.0.1", "Domaine ou IP de la machine qui expose le sercice.")

var portEnv = flag.String("port-env", "", "Variable d'environement où trouver le port d'exposition de l'application.")
var hostEnv = flag.String("host-env", "", "Variable d'environement où trouver le domaine ou l'IP de la machine qui expose le sercice.")

var generateSignature = flag.Bool("generate-signature", false, "Active l'option de génération du token de signature pour l'utilisateur.")

func GetFlags() Flags {
	if !flag.Parsed() {
		flag.Parse()
	}

	return Flags{
		port,
		host,
		portEnv,
		hostEnv,
		generateSignature,
	}
}
