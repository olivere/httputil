// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSONError(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		WriteJSONError(w, "something happened")
	}

	req, err := http.NewRequest("GET", "http://example.com/", nil)
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
	if fail.Error.Code != 500 {
		t.Errorf("expected error code = %d; got: %d", 500, fail.Error.Code)
	}
	if fail.Error.Message != "something happened" {
		t.Errorf("expected error message = %q; got: %q", "something happened", fail.Error.Message)
	}
}

func TestWriteJSONErrorWithCoder(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		WriteJSONError(w, InvalidParameterError("pin"))
	}

	req, err := http.NewRequest("GET", "http://example.com/", nil)
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
	if fail.Error.Message != `Invalid parameter "pin"` {
		t.Errorf("expected error message = %q; got: %q", `Invalid parameter "pin"`, fail.Error.Message)
	}
}

func TestWriteJSONErrorWithDetails(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		err := UnprocessableEntityError{}
		err.Errors = []string{"A has failed", "B is invalid"}
		WriteJSONError(w, err)
	}

	req, err := http.NewRequest("GET", "http://example.com/", nil)
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
	if fail.Error.Code != 422 {
		t.Errorf("expected error code = %d; got: %d", 422, fail.Error.Code)
	}
	if fail.Error.Message != `Record has semantic errors` {
		t.Errorf("expected error message = %q; got: %q", `Record has semantic errors`, fail.Error.Message)
	}
	if len(fail.Error.Details) != 2 {
		t.Fatalf("expected 2 error details; got: %d", len(fail.Error.Details))
	}
	if fail.Error.Details[0] != `A has failed` {
		t.Errorf("expected error details[0] = %q; got: %q", `A has failed`, fail.Error.Details[0])
	}
	if fail.Error.Details[1] != `B is invalid` {
		t.Errorf("expected error details[1] = %q; got: %q", `B is invalid`, fail.Error.Details[1])
	}
}
