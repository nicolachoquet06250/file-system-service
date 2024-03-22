package main

import (
	"filesystem_service/actions"
	"filesystem_service/auth/tokens"
	"filesystem_service/customHttp"
	"filesystem_service/directories"
	"filesystem_service/files"
	"filesystem_service/swagger"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "modernc.org/sqlite"
)

func main() {
	// TODO maintenant que la génération de credentials est faite, gérer la génération et la validation de tokens.
	if !actions.Exec() {
		addr := customHttp.GetAddr()

		server := http.NewServeMux()

		server.HandleFunc("/check-validity", CheckValidity) // ok

		server.HandleFunc("OPTIONS /auth/get-token", customHttp.Cors(customHttp.AcceptOptions)) // ok
		server.HandleFunc("POST /auth/get-token", customHttp.Cors(tokens.GetToken))             // ok
		server.HandleFunc("PUT /auth/get-token", customHttp.Cors(tokens.RefreshToken))          // ok

		server.HandleFunc("OPTIONS /file-system/{path...}", customHttp.Cors(customHttp.AcceptOptions)) // ok
		server.HandleFunc("/file-system/{path...}", customHttp.Cors(directories.GetFileSystem))        // ok

		server.HandleFunc("OPTIONS /directory", customHttp.Cors(customHttp.AcceptOptions))             // ok
		server.HandleFunc("POST /directory", customHttp.Cors(directories.CreateDirectory))             // ok
		server.HandleFunc("OPTIONS /directory/{path...}", customHttp.Cors(customHttp.AcceptOptions))   // ok
		server.HandleFunc("PATCH /directory/{path...}", customHttp.Cors(directories.RenameDirectory))  // ok
		server.HandleFunc("DELETE /directory/{path...}", customHttp.Cors(directories.DeleteDirectory)) // ok

		server.HandleFunc("OPTIONS /file/{path...}", customHttp.Cors(customHttp.AcceptOptions)) // ok
		server.HandleFunc("/file/{path...}", customHttp.Cors(files.GetFileContent))             // ok
		server.HandleFunc("PATCH /file/{path...}", customHttp.Cors(files.RenameFile))           // ok
		server.HandleFunc("PUT /file/{path...}", customHttp.Cors(files.UpdateFileContent))      // ok
		server.HandleFunc("DELETE /file/{path...}", customHttp.Cors(files.DeleteFile))          // ok
		server.HandleFunc("OPTIONS /file", customHttp.Cors(customHttp.AcceptOptions))           // ok
		server.HandleFunc("POST /file", customHttp.Cors(files.CreateFile))                      // ok

		server.HandleFunc("/swagger.json", swagger.DefinitionRoute) // ok
		server.HandleFunc("/swagger.yaml", swagger.DefinitionRoute) // ok
		server.HandleFunc("/swagger.yml", swagger.DefinitionRoute)  // ok

		server.HandleFunc("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger.json"))) // ok

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			println("\b\bNous fermons le serveur.")
			// Run Cleanup
			os.Exit(0)
		}()

		fmt.Printf("Listening on http://%s\n", addr)
		if err := http.ListenAndServe(addr, server); err != nil {
			println(err)
		}
	}
}
