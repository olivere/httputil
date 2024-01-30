// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestMustFormString(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		defer Recover(w, r)
		fmt.Fprint(w, MustFormString(r, "name"))
	}

	values := url.Values{"name": {"Oliver"}}
	req, err := http.NewRequest("POST", "http://localhost/", strings.NewReader(values.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	h(w, req)
	if w.Code != 200 {
		t.Fatalf("expected status = %d; got: %d", 200, w.Code)
	}
	body := w.Body.String()
	if body != "Oliver" {
		t.Errorf("expected body = %q; got: %q", "Oliver", body)
	}
}

func TestMustFormStringFailure(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		defer Recover(w, r)
		fmt.Fprint(w, MustFormString(r, "name"))
	}

	values := url.Values{}
	req, err := http.NewRequest("POST", "http://localhost/", strings.NewReader(values.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	h(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status = %d; got: %d", http.StatusBadRequest, w.Code)
	}
}
