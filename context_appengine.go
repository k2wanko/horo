//+build appengine

package horo

import "net/http"

func requestID(r *http.Request) string {
	return r.Header.Get("X-AppEngine-Request-Log-Id")
}
