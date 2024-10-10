package parkserver

import (
	"net/http"
	"strings"
	"time"
)

var scalarDocsPage = strings.NewReader(`<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`)

func handleDocs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodHead:
	default:
		http.Error(w, "", http.StatusNotImplemented)
	}

	w.Header().Set("Content-Type", "text/html")
	http.ServeContent(w, r, "", time.Time{}, scalarDocsPage)
}
