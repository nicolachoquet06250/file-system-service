package database

import "os"

var prodPath = func() string {
	path := "C://Program Files/filesystem-service"
	_, err := os.Stat(path)
	if err != nil {
		_ = os.MkdirAll(path, 0777)
	}
	return path
}()
