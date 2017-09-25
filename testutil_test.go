// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil

import "testing"

func TestEqualJSON(t *testing.T) {
	tests := []struct {
		A, B []byte
		Want bool
	}{
		{
			A:    nil,
			B:    nil,
			Want: true,
		},
		{
			A:    nil,
			B:    []byte{},
			Want: true,
		},
		{
			A:    []byte(`{"a":1,"b":2}`),
			B:    []byte(`{"a":1,"b":2}`),
			Want: true,
		},
		{
			A:    []byte(`{"a":1,"b":2}     `),
			B:    []byte(`{"a":1,"b":2}`),
			Want: true,
		},
		{
			A: []byte(`{"a":1,"b":2}`),
			B: []byte(`{
	"a":1,
	"b":2
}
`),
			Want: true,
		},
		{
			A:    []byte(`{"a":1,"b":2}`),
			B:    []byte(`{"b":2,"a":1}`),
			Want: false,
		},
	}

	for _, tt := range tests {
		if want, have := tt.Want, EqualJSON(tt.A, tt.B); want != have {
			t.Errorf("EqualJSON(%v,%v): want %t, have %t", tt.A, tt.B, want, have)
		}
	}
}
