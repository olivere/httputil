// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil

import (
	"bytes"
	"encoding/json"
)

// EqualJSON compares the two serialized byte slices for equality.
//
// Both a and b are expected to be JSON serialized byte slices.
// EqualJSON simply ensures that insiginificant white space is
// removed both from a and b before comparing for equality.
//
// EqualJSON returns true in the following cases:
// 1. a or b are both nil
// 2. a or b have both a length of 0
// 3. a or b are equal if all siginificant white space is removed,
//    i.e. newlines, tabs, and space.
//
// In all other cases, EqualJSON returns false. Notice that while the
// two JSON objects `{"a":1,"b":2}` and `{"b":2,"a":1}` may be semantically
// equal, EqualJSON will return false.
func EqualJSON(a, b []byte) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	var dsta, dstb bytes.Buffer
	if err := json.Compact(&dsta, a); err != nil {
		return false
	}
	if err := json.Compact(&dstb, b); err != nil {
		return false
	}
	return bytes.Equal(dsta.Bytes(), dstb.Bytes())
}
