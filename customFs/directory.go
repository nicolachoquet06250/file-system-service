package customFs

import (
	"fmt"
	"os"
	"strings"
)

type Directory struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type ElementListItem struct {
	Type    string `json:"type"`
	Path    string `json:"path"`
	Name    string `json:"name"`
	SymLink string `json:"symlink,omitempty"`
}

/****************************************************************/
/**** Constructor                                            ****/
/****************************************************************/

func NewDirectory(path string) *Directory {
	path, name := splitCompletePathForDirectory(path)

	return &Directory{path, name}
}

/****************************************************************/
/**** Body Methods                                           ****/
/****************************************************************/

/********************************************************/
/**** Private Methods                                ****/
/********************************************************/

func (d *Directory) buildCompletePath() string {
	return BuildDirectoryCompletePath(*d)
}

/********************************************************/
/**** Public Methods                                 ****/
/********************************************************/

func (d *Directory) IsDir() (bool, error) {
	fi, err := os.Stat(d.GetPath())

	if err != nil {
		return false, err
	}

	if fi.Mode().IsRegular() || !fi.Mode().IsDir() {
		return false, fmt.Errorf("open %s is a not a directory", d.GetPath())
	}

	fi, err = os.Stat(d.buildCompletePath())

	if err != nil {
		return false, err
	}

	if fi.Mode().IsRegular() || !fi.Mode().IsDir() {
		return false, fmt.Errorf("open %s is not a directory", d.buildCompletePath())
	}

	return true, nil
}

func (d *Directory) Create() (bool, error) {
	if exists, err := d.IsDir(); err != nil {
		if !strings.Contains(err.Error(), d.buildCompletePath()) {
			return false, err
		}
	} else if exists {
		return false, fmt.Errorf("create %s already exists", d.buildCompletePath())
	}

	if err := os.Mkdir(d.buildCompletePath(), 0777); err != nil {
		return false, err
	}
	return true, nil
}

func (d *Directory) Delete() (bool, error) {
	if _, err := d.IsDir(); err != nil {
		return false, err
	}

	if err := os.Remove(d.buildCompletePath()); err != nil {
		return false, err
	}
	return true, nil
}

func (d *Directory) DeepDelete() (bool, error) {
	if isExists, err := d.IsDir(); !isExists {
		return false, err
	}

	err := os.RemoveAll(d.buildCompletePath())
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *Directory) Rename(newDir *Directory) (bool, error) {
	if _, err := d.IsDir(); err != nil {
		return false, err
	}

	if err := os.Rename(d.buildCompletePath(), newDir.buildCompletePath()); err != nil {
		return false, err
	}
	return true, nil
}

/************************************************/
/**** Built Getters                          ****/
/************************************************/

func (d *Directory) GetFlatContent() (list []ElementListItem, err error) {
	path := d.buildCompletePath()

	entries, err := os.ReadDir(path)
	if err != nil {
		return list, err
	}

	for _, e := range entries {
		Type := "file"
		if e.IsDir() {
			Type = "directory"
		}

		item := ElementListItem{
			Type: Type,
			Name: e.Name(),
			Path: path,
		}

		_path := path
		if _path == "/" {
			_path += "/"
		}
		_path += e.Name()

		if link, err := os.Readlink(_path); err == nil {
			item.SymLink = link

			if !strings.Contains(_path, ".") {
				item.Type = "directory"
			}
		}

		list = append(list, item)
	}

	return list, nil
}

/************************************************/
/**** Getters                                ****/
/************************************************/

func (d *Directory) GetPath() string {
	return d.Path
}

func (d *Directory) GetDirName() string {
	return d.Name
}
