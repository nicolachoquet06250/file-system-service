package actions

import (
	"filesystem_service/auth/credentials"
	"filesystem_service/auth/roles"
	"filesystem_service/flags"
)

func Exec() bool {
	fl := flags.GetFlags()

	if fl.IsShowRoles() {
		roles.ShowRoles()
		return true
	}

	if fl.IsUpdateCredentials() {
		credentials.UpdateCredentials(fl.GetClientId(), fl.GetRole())
		return true
	}

	if fl.IsShowUserRole() {
		roles.GetUserRole(fl.GetClientId())
		return true
	}

	if fl.IsGenerateCredentials() {
		credentials.GenerateCredentials(fl.GetRole())
		return true
	}

	return false
}
