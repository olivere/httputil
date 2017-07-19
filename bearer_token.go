// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil

import (
	"net/http"
	"strings"
)

// BearerToken extracts the Bearer token from the request.
func BearerToken(r *http.Request) (string, bool) {
	const prefix = "bearer "
	auth := r.Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(strings.ToLower(auth), prefix) {
		return "", false
	}
	ts := auth[len(prefix):]
	if ts == "" {
		return "", false
	}
	return ts, true
}
