package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseStatus struct {
	Status string `json:"status"`
}

func getAddr(host *string, port *int, hostEnv *string, portEnv *string) string {
	if hostEnv != nil {
		if envHost := os.Getenv(*hostEnv); envHost != "" {
			*host = envHost
		}
	}
	if portEnv != nil {
		if envPort := os.Getenv(*portEnv); envPort != "" {
			*port, _ = strconv.Atoi(envPort)
		}
	}

	if host != nil && strings.Contains(*host, ":") {
		*host = fmt.Sprintf("[%s]", *host)
	}

	return *host + ":" + strconv.Itoa(*port)
}

func main() {
	port := flag.Int("port", 3000, "Port d'exposition de l'application.")
	host := flag.String("host", "127.0.0.1", "Domaine ou IP de la machine qui expose le sercice.")

	portEnvVar := flag.String("portEnv", "", "Variable d'environement où trouver le port d'exposition de l'application.")
	hostEnvVar := flag.String("hostEnv", "", "Variable d'environement où trouver le domaine ou l'IP de la machine qui expose le sercice.")
	flag.Parse()

	addr := getAddr(host, port, hostEnvVar, portEnvVar)

	server := http.NewServeMux()

	server.HandleFunc("/check-validity", checkValidity) // ok

	server.HandleFunc("/file-system/{path...}", getFileSystem) // ok

	server.HandleFunc("POST /directory", createDirectory)             // ok
	server.HandleFunc("PATCH /directory", renameDirectory)            // ok
	server.HandleFunc("PATCH /directory/{path...}", renameDirectory)  // ok
	server.HandleFunc("DELETE /directory", deleteDirectory)           // ok
	server.HandleFunc("DELETE /directory/{path...}", deleteDirectory) // ok

	server.HandleFunc("/file/{path...}", getFileContent)        // ok
	server.HandleFunc("POST /file", createFile)                 // ok
	server.HandleFunc("PATCH /file/{path...}", renameFile)      // ok
	server.HandleFunc("PUT /file/{path...}", updateFileContent) // ok
	server.HandleFunc("DELETE /file/{path...}", deleteFile)     // ok

	fmt.Printf("Listening on http://%s", addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		println(err)
	}
}
