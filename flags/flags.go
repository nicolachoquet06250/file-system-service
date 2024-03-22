package flags

import (
	"os"
	"strconv"
)

type Flags struct {
	port    *int
	host    *string
	portEnv *string
	hostEnv *string

	generateCredentials *bool
	updateCredentials   *bool
	showRoles           *bool
	showUserRole        *bool

	role     *string
	clientId *string
}

func (f *Flags) GetPort() int {
	if envPort := os.Getenv(*f.portEnv); f.portEnv != nil && envPort != "" {
		port, _ := strconv.Atoi(envPort)
		return port
	}
	return *f.port
}

func (f *Flags) GetHost() string {
	if envHost := os.Getenv(*f.hostEnv); f.hostEnv != nil && envHost != "" {
		return envHost
	}
	return *f.host
}

func (f *Flags) IsGenerateCredentials() bool {
	return !(f.generateCredentials == nil || *f.generateCredentials == false)
}

func (f *Flags) IsUpdateCredentials() bool {
	return !(f.updateCredentials == nil || *f.updateCredentials == false)
}

func (f *Flags) IsShowRoles() bool {
	return !(f.showRoles == nil || *f.showRoles == false)
}

func (f *Flags) IsShowUserRole() bool {
	return !(f.showUserRole == nil || *f.showUserRole == false)
}

func (f *Flags) GetRole() string {
	if !(f.role == nil || *f.role == "") {
		return *f.role
	}
	return ""
}

func (f *Flags) GetClientId() string {
	return *f.clientId
}
