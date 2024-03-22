package swagger

import (
	_ "embed"
	"net/http"
	"strings"
)

//go:embed swagger.json
var swaggerJson []byte

//go:embed swagger.yaml
var swaggerYaml []byte

type associatedToExtension struct {
	mimeType string
	content  []byte
}

var extensions = map[string]associatedToExtension{
	".json": {
		mimeType: "application/json",
		content:  swaggerJson,
	},
	".yaml": {
		mimeType: "application/x-yaml",
		content:  swaggerYaml,
	},
	".yml": {
		mimeType: "application/x-yaml",
		content:  swaggerYaml,
	},
}

func DefinitionRoute(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path

	if strings.HasPrefix(path, "/swagger") {
		for ext, detail := range extensions {
			if strings.HasSuffix(path, ext) {
				writer.Header().Add("Content-Type", detail.mimeType)
				_, _ = writer.Write(detail.content)
				return
			}
		}
	}

	http.NotFound(writer, request)
}
