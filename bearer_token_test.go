// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.
package httputil

import (
	"net/http"
	"testing"
)

func TestBearerToken(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "Bearer secret")

	token, ok := BearerToken(r)
	if !ok {
		t.Fatal("expected to find bearer token")
	}
	if token != "secret" {
		t.Fatalf("expected %q, got %q", "sceret", token)
	}
}
