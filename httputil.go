// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// IsGetOrHead returns true if r is a GET or HEAD request.
func IsGetOrHead(r *http.Request) bool {
	return r.Method == "GET" || r.Method == "HEAD"
}

// IsWebsocketUpgrade returns true if this is a WebSocket upgrade.
func IsWebsocketUpgrade(req *http.Request) bool {
	return req.Method == "GET" && req.Header.Get("Upgrade") == "websocket"
}

// IsXHR returns true if r is an XHR request. It inspects the
// Content-Type header for that.
func IsXHR(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "application/json")
}

var byteBufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// ReadJSON deserializes the body of the request into dst as JSON.
// A maximum size of 8 MB of JSON are permitted.
func ReadJSON(r *http.Request, dst interface{}) error {
	buf := byteBufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		byteBufPool.Put(buf)
	}()
	// Limit to 8 MB of JSON
	if err := json.NewDecoder(io.TeeReader(io.LimitReader(r.Body, 8<<20), buf)).Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON data: %v, on input: %s", err, buf.Bytes())
	}
	return nil
}

// MustReadJSON is like ReadJSON, but panics on errors.
func MustReadJSON(r *http.Request, dst interface{}) {
	if err := ReadJSON(r, dst); err != nil {
		panic(BadRequestError{Message: "Invalid JSON data", Err: err})
	}
}

// CloseBody closes rc.
func CloseBody(rc io.ReadCloser) {
	if rc != nil {
		rc.Close()
	}
}

// WriteJSON writes data as JSON into w with HTTP status code 200.
func WriteJSON(w http.ResponseWriter, data interface{}) {
	WriteJSONCode(w, http.StatusOK, data)
}

// WriteJSONCode writes data as JSON into w and sets the HTTP status code.
func WriteJSONCode(w http.ResponseWriter, code int, data interface{}) {
	js, _ := json.MarshalIndent(data, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
	w.Write([]byte("\n"))
}

// Recover can be used as a deferred func to catch panics in an HTTP handler.
func Recover(w http.ResponseWriter, r *http.Request) {
	err := recover()
	if err != nil {
		WriteError(w, err)
	}
}

// RecoverJSON can be used as a deferred func to catch panics in an HTTP handler
// and print a JSON error.
//
// Example:
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//	  defer httputil.RecoverJSON(w, r)
//	  ...
//	  panic(errors.New("kaboom"))
//	}
func RecoverJSON(w http.ResponseWriter, r *http.Request) {
	err := recover()
	if err != nil {
		WriteJSONError(w, err)
	}
}
