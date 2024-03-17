package main

import (
	"filesystem_service/actions"
	"filesystem_service/auth"
	"filesystem_service/customHttp"
	"filesystem_service/directories"
	"filesystem_service/files"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "modernc.org/sqlite"
)

func main() {
	if !actions.Exec() {
		addr := customHttp.GetAddr()

		server := http.NewServeMux()

		server.HandleFunc("/check-validity", CheckValidity) // ok

		server.HandleFunc("POST /auth/get-token", auth.GetToken)
		server.HandleFunc("PUT /auth/get-token", auth.RefreshToken)

		server.HandleFunc("/file-system/{path...}", directories.GetFileSystem) // ok

		server.HandleFunc("POST /directory", directories.CreateDirectory)             // ok
		server.HandleFunc("PATCH /directory/{path...}", directories.RenameDirectory)  // ok
		server.HandleFunc("DELETE /directory/{path...}", directories.DeleteDirectory) // ok

		server.HandleFunc("/file/{path...}", files.GetFileContent)        // ok
		server.HandleFunc("POST /file", files.CreateFile)                 // ok
		server.HandleFunc("PATCH /file/{path...}", files.RenameFile)      // ok
		server.HandleFunc("PUT /file/{path...}", files.UpdateFileContent) // ok
		server.HandleFunc("DELETE /file/{path...}", files.DeleteFile)     // ok

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
