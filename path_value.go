//go:build go1.22
// +build go1.22

package chi

import "net/http"

// supportPathValue indicates whether the running Go version provides the
// http.Request.SetPathValue / PathValue methods (Go 1.22+).
const supportPathValue = true

// setPathValues copies the URL parameters captured by chi during routing onto
// the request using http.Request.SetPathValue, so that r.PathValue(key) acts
// as an alias for chi.URLParam(r, key) on Go 1.22+.
func setPathValues(rctx *Context, r *http.Request) {
	for i, key := range rctx.URLParams.Keys {
		value := rctx.URLParams.Values[i]
		r.SetPathValue(key, value)
	}
}
