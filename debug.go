// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil

import (
	"fmt"
	"io"
	"net/http"
	stdhttputil "net/http/httputil"
)

// DumpRequestOut prints the request to the given io.Writer.
func DumpRequestOut(w io.Writer, r *http.Request) {
	data, _ := stdhttputil.DumpRequestOut(r, true)
	fmt.Fprint(w, string(data))
}
