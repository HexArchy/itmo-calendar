package api

import (
	_ "embed"
	"net/http"
)

//go:embed openapi/openapi.yaml
var openapiSpec []byte

// SpecHandler returns an [http.Handler] that serves the OpenAPI spec.
func SpecHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		if _, err := w.Write(openapiSpec); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// ScalarDocsHandler returns an [http.Handler] that serves the Scalar API reference UI.
func ScalarDocsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := w.Write([]byte(scalarHTML)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

const scalarHTML = `<!DOCTYPE html>
<html>
<head>
  <title>ITMO Calendar API</title>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
  <script id="api-reference" data-url="/openapi/openapi.yaml"></script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`
