// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadJSON(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString(`{"message":"hello"}`)
	req, err := http.NewRequest("GET", "http://localhost/", &buf)
	if err != nil {
		t.Fatal(err)
	}
	type coding struct {
		Message string `json:"message"`
	}
	var dst coding
	err = ReadJSON(req, &dst)
	if err != nil {
		t.Fatal(err)
	}
	if dst.Message != "hello" {
		t.Errorf("expected %q; got: %q", "hello", dst.Message)
	}
}

func TestReadJSONFailure(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString(`{"message"}`)
	req, err := http.NewRequest("GET", "http://localhost/", &buf)
	if err != nil {
		t.Fatal(err)
	}
	type coding struct {
		Message string `json:"message"`
	}
	var dst coding
	err = ReadJSON(req, &dst)
	if err == nil {
		t.Fatal("expected ReadJSON to fail")
	}
}

func TestMustReadJSON(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		defer RecoverJSON(w, r)

		type coding struct {
			Message string `json:"message"`
		}
		var dst coding
		MustReadJSON(r, &dst)

		fmt.Fprint(w, "ok\n")
	}

	var buf bytes.Buffer
	buf.WriteString(`{"message"}`)
	req, err := http.NewRequest("GET", "http://localhost/", &buf)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	h(w, req)

	got := w.Header().Get("Content-Type")
	if got != "application/json" {
		t.Errorf("expected Content-Type = %q; got: %q", "application/json", got)
	}
	type failure struct {
		Error struct {
			Code    int      `json:"code"`
			Message string   `json:"message"`
			Details []string `json:"details"`
		} `json:"error"`
	}
	var fail failure
	err = json.NewDecoder(w.Body).Decode(&fail)
	if err != nil {
		t.Fatal(err)
	}
	if fail.Error.Code != http.StatusBadRequest {
		t.Errorf("expected error code = %d; got: %d", http.StatusBadRequest, fail.Error.Code)
	}
	if !strings.HasPrefix(fail.Error.Message, "Invalid JSON data") {
		t.Errorf("unexpected error message: %q", fail.Error.Message)
	}
}
