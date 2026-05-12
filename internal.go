package aigchq

import (
	"net/http"
	"net/url"
)

const (
	httpMethodGet    = http.MethodGet
	httpMethodPost   = http.MethodPost
	httpMethodPatch  = http.MethodPatch
	httpMethodDelete = http.MethodDelete
)

func pathEscape(value string) string {
	return url.PathEscape(value)
}
