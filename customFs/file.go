package customFs

import (
	"fmt"
	"os"
	"strings"
)

var fileFormats = map[string]string{
	".json": "application/json",
	".xml":  "text/xml",
	".txt":  "text/plain",
	".md":   "text/markdown",
	".":     "text/plain",
	".pdf":  "application/pdf",
}

type File struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
}

/****************************************************************/
/**** Constructor                                            ****/
/****************************************************************/

func NewFile(path string) *File {
	path, name, extension := splitCompletePathForFile(path)

	return &File{path, name, extension}
}

/****************************************************************/
/**** Body Methods                                           ****/
/****************************************************************/

/********************************************************/
/**** Private Methods                                ****/
/********************************************************/

func (f *File) buildCompletePath() string {
	return BuildFileCompletePath(*f)
}

/********************************************************/
/**** Public Methods                                 ****/
/********************************************************/

func (f *File) IsFile() (bool, error) {
	fi, err := os.Stat(f.GetPath())

	if err != nil {
		return false, err
	}

	if fi.Mode().IsRegular() || !fi.Mode().IsDir() {
		return false, fmt.Errorf("open %s is a not a directory", f.GetPath())
	}

	fi, err = os.Stat(f.buildCompletePath())

	if err != nil {
		return false, err
	}

	if fi.Mode().IsDir() || !fi.Mode().IsRegular() {
		return false, fmt.Errorf("open %s is a directory", f.buildCompletePath())
	}

	return true, nil
}

func (f *File) Create() (bool, error) {
	if exists, err := f.IsFile(); err != nil {
		if !strings.Contains(err.Error(), f.buildCompletePath()) {
			return false, err
		}
	} else if exists {
		return false, fmt.Errorf("create %s already exists", f.buildCompletePath())
	}

	if file, err := os.Create(f.buildCompletePath()); err != nil {
		defer file.Close()
		return false, err
	}
	return true, nil
}

func (f *File) Rename(newFile *File) (bool, error) {
	if isExists, err := f.IsFile(); !isExists {
		return false, err
	}

	if err := os.Rename(f.buildCompletePath(), newFile.buildCompletePath()); err != nil {
		return false, err
	}

	f.Path = newFile.GetPath()
	f.Name = newFile.GetFileName()
	f.Extension = newFile.GetExtension()

	return true, nil
}

func (f *File) Delete() (bool, error) {
	if isExists, err := f.IsFile(); !isExists {
		return false, err
	}

	err := os.Remove(f.buildCompletePath())
	if err != nil {
		return false, err
	}
	return true, nil
}

/************************************************/
/**** Built Getters                          ****/
/************************************************/

func (f *File) GetContent() ([]byte, error) {
	path := f.buildCompletePath()

	if isFile, _ := f.IsFile(); isFile {
		return os.ReadFile(path)
	}
	return []byte{}, fmt.Errorf("open %s no such file or directory", path)
}

func (f *File) SetContent(content []byte) (bool, error) {
	if isFile, _ := f.IsFile(); !isFile {
		created, err := f.Create()
		if err != nil {
			return created, err
		}
	}

	err := os.WriteFile(f.buildCompletePath(), content, 0777)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (f *File) GetMimeType() string {
	if isFile, _ := f.IsFile(); isFile {
		return fileFormats[f.GetExtension()]
	}
	return fileFormats["."]
}

/************************************************/
/**** Getters                                ****/
/************************************************/

func (f *File) GetExtension() string {
	return f.Extension
}

func (f *File) GetPath() string {
	return f.Path
}

func (f *File) GetFileName() string {
	return f.Name
}
