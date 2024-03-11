package customFs

import (
	"strings"
)

func getExtension(path []string) string {
	var extParts []string

	for _, part := range path {
		if !strings.Contains(part, "/") {
			extParts = append(extParts, part)
		}
	}

	return "." + strings.Join(extParts, ".")
}

func getName(path string) string {
	parts := strings.Split(path, "/")

	return parts[len(parts)-1]
}

func getParentPath(path string) string {
	parts := strings.Split(path, "/")

	finalPath := strings.Join(parts[0:len(parts)-1], "/")
	if finalPath == "" {
		finalPath = "/" + finalPath
	}
	return finalPath
}

func splitCompletePathForFile(
	completePath string,
) (path string, name string, ext string) {
	ext = getExtension(strings.Split(completePath, "."))
	completePath = strings.Replace(completePath, ext, "", 1)
	name = getName(completePath)
	path = getParentPath(completePath)

	return
}

func splitCompletePathForDirectory(
	completePath string,
) (path string, name string) {
	name = getName(completePath)
	path = getParentPath(completePath)

	return
}

func BuildFileCompletePath(f File) string {
	path := f.Path
	if path != "/" {
		path += "/"
	}
	path += f.Name + f.Extension

	return path
}

func BuildDirectoryCompletePath(d Directory) string {
	path := d.Path
	if path != "/" {
		path += "/"
	}
	path += d.Name

	return path
}
