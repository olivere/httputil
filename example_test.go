// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"net/http/httptest"

	"github.com/olivere/httputil"
)

func ExampleMustQueryString() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		defer httputil.RecoverJSON(w, r)

		name := httputil.MustQueryString(r, "name")
		fmt.Fprintf(w, "Hello %s", name)
	}

	req := httptest.NewRequest("GET", "http://example.com/hello?name=Oliver", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
	// Output: Hello Oliver
}

func ExampleMissingParameterError() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		defer httputil.RecoverJSON(w, r)

		name := httputil.QueryString(r, "name", "")
		if name == "" {
			panic(httputil.MissingParameterError("name"))
		}

		httputil.WriteJSONCode(w, http.StatusOK, struct {
			Message string `json:"name"`
		}{
			Message: fmt.Sprintf("Hello %s", name),
		})
	}

	req := httptest.NewRequest("GET", "http://example.com/hello?name=", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("Status: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Println(string(body))
	// Output:
	// Status: 400 Bad Request
	// {
	//   "error": {
	//     "code": 400,
	//     "message": "Missing parameter \"name\""
	//   }
	// }
}

func ExampleReadJSON() {
	req := &http.Request{
		Body: ioutil.NopCloser(strings.NewReader(`{"name":"Oliver"}`)),
	}

	var person struct {
		Name string `json:"name"`
	}
	httputil.MustReadJSON(req, &person)

	fmt.Println(person.Name)
	// Output: Oliver
}

func ExampleWriteJSON() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		person := struct {
			Name string `json:"name"`
		}{
			Name: "Oliver",
		}
		httputil.WriteJSON(w, person)
	}

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
	// Output:
	// {
	//   "name": "Oliver"
	// }
}

func ExampleWriteJSONCode() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		person := struct {
			Name string `json:"name"`
		}{
			Name: "Oliver",
		}
		httputil.WriteJSONCode(w, http.StatusCreated, person)
	}

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("Status: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Println(string(body))
	// Output:
	// Status: 201 Created
	// {
	//   "name": "Oliver"
	// }
}
