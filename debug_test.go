// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func chunk(s string) string {
	return fmt.Sprintf("%x\r\n%s\r\n", len(s), s)
}

func TestDumpRequestOut(t *testing.T) {
	body := []byte("<html>Body</html>")
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   "www.alt-f4.de",
			Path:   "/",
		},
		ProtoMajor:       1,
		ProtoMinor:       1,
		TransferEncoding: []string{"chunked"},
		Header: http.Header{
			"Authorization": []string{"Bearer secret"},
		},
		Body: ioutil.NopCloser(bytes.NewReader(body)),
	}

	var buf bytes.Buffer
	DumpRequestOut(&buf, req)

	expected := "GET / HTTP/1.1\r\n" +
		"Host: www.alt-f4.de\r\n" +
		"User-Agent: Go-http-client/1.1\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"Authorization: Bearer secret\r\n" +
		"Accept-Encoding: gzip\r\n" +
		"\r\n" +
		chunk(string(body)) +
		chunk("")
	if have, want := buf.String(), expected; have != want {
		t.Fatalf("have:\n%q\nwant:\n%q", have, want)
	}
}
