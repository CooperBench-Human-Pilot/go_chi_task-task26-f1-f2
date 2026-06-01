//go:build !go1.22
// +build !go1.22

package chi

import "net/http"

// supportPathValue indicates whether the running Go version provides the
// http.Request.SetPathValue / PathValue methods (Go 1.22+). On older versions
// these methods don't exist, so the value is false.
const supportPathValue = false

// setPathValues is a no-op on Go versions before 1.22, which lack the
// http.Request.SetPathValue method.
func setPathValues(rctx *Context, r *http.Request) {}
