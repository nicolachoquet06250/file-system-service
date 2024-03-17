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

	generateSignature *bool
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

func (f *Flags) IsGenerateSignature() bool {
	return !(f.generateSignature == nil || *f.generateSignature == false)
}
