package actions

import (
	"filesystem_service/auth"
	"filesystem_service/flags"
)

func Exec() bool {
	fl := flags.GetFlags()

	if fl.IsGenerateSignature() {
		auth.GenerateSignatureTokenAction()
		return true
	}

	return false
}
