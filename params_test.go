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
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
)

var (
	zeroTime        time.Time
	zeroDuration    time.Duration
	defaultTime     = time.Now()
	defaultDuration = time.Duration(42198) * time.Millisecond
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

var (
	paramsTests = []struct {
		Values         url.Values
		ExpectedErr    error
		ExpectedOutput interface{}
		Invoke         func(r *http.Request, key string) interface{}
	}{
		// -- MustFormString --
		// #0
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: "",
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormString(r, key)
			},
		},
		// #1
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: "",
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormString(r, key)
			},
		},
		// #2
		{
			Values:         url.Values{"key": []string{"Rob"}},
			ExpectedErr:    nil,
			ExpectedOutput: "Rob",
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormString(r, key)
			},
		},
		// -- MustFormBool --
		// #3
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormBool(r, key)
			},
		},
		// #4
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormBool(r, key)
			},
		},
		// #5
		{
			Values:         url.Values{"key": []string{"true"}},
			ExpectedErr:    nil,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormBool(r, key)
			},
		},
		// #6
		{
			Values:         url.Values{"key": []string{"false"}},
			ExpectedErr:    nil,
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormBool(r, key)
			},
		},
		// #7
		{
			Values:         url.Values{"key": []string{"1"}},
			ExpectedErr:    nil,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormBool(r, key)
			},
		},
		// #8
		{
			Values:         url.Values{"key": []string{"0"}},
			ExpectedErr:    nil,
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormBool(r, key)
			},
		},
		// -- MustFormInt --
		// #9
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: 0,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt(r, key)
			},
		},
		// #10
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: 0,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt(r, key)
			},
		},
		// #11
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt(r, key)
			},
		},
		// -- MustFormInt32 --
		// #12
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt32(r, key)
			},
		},
		// #13
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt32(r, key)
			},
		},
		// #14
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt32(r, key)
			},
		},
		// -- MustFormInt64 --
		// #15
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt64(r, key)
			},
		},
		// #16
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt64(r, key)
			},
		},
		// #17
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt64(r, key)
			},
		},
		// -- MustFormFloat32 --
		// #18
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormFloat32(r, key)
			},
		},
		// #19
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormFloat32(r, key)
			},
		},
		// #20
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: float32(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormFloat32(r, key)
			},
		},
		// -- MustFormFloat64 --
		// #21
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormFloat64(r, key)
			},
		},
		// #22
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormFloat64(r, key)
			},
		},
		// #23
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: float64(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormFloat64(r, key)
			},
		},
		// -- MustFormTime --
		// #24
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTime(r, key, time.RFC3339)
			},
		},
		// #25
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTime(r, key, time.RFC3339)
			},
		},
		// #26
		{
			Values:         url.Values{"key": []string{"2018-12-31T18:47:59.999999999Z"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 18, 47, 59, 999999999, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTime(r, key, time.RFC3339)
			},
		},
		// #27
		{
			Values:         url.Values{"key": []string{"2018-12-31"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTime(r, key, "2006-01-02")
			},
		},
		// -- MustFormTimeWithDefault --
		// #28
		{
			Values:         url.Values{},
			ExpectedErr:    nil,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTimeWithDefault(r, key, time.RFC3339, defaultTime)
			},
		},
		// #29
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    nil,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTimeWithDefault(r, key, time.RFC3339, defaultTime)
			},
		},
		// #30
		{
			Values:         url.Values{"key": []string{"2018-12-31T18:47:59.999999999Z"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 18, 47, 59, 999999999, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTimeWithDefault(r, key, time.RFC3339Nano, defaultTime)
			},
		},
		// #31
		{
			Values:         url.Values{"key": []string{"2018-12-31"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTimeWithDefault(r, key, "2006-01-02", defaultTime)
			},
		},
		// -- MustFormDuration --
		// #32
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormDuration(r, key)
			},
		},
		// #33
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormDuration(r, key)
			},
		},
		// #34
		{
			Values:         url.Values{"key": []string{"1m12s"}},
			ExpectedErr:    nil,
			ExpectedOutput: 72 * time.Second,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormDuration(r, key)
			},
		},
		// -- MustFormDurationWithDefault --
		// #35
		{
			Values:         url.Values{},
			ExpectedErr:    nil,
			ExpectedOutput: defaultDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormDurationWithDefault(r, key, defaultDuration)
			},
		},
		// #36
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    nil,
			ExpectedOutput: defaultDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormDurationWithDefault(r, key, defaultDuration)
			},
		},
		// #37
		{
			Values:         url.Values{"key": []string{"1m12s"}},
			ExpectedErr:    nil,
			ExpectedOutput: 72 * time.Second,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormDurationWithDefault(r, key, defaultDuration)
			},
		},
	}
)

func TestParams(t *testing.T) {
	t.Run("MustForm", testParamsMustFormXXX)
	t.Run("Form", testParamsFormXXX)
	t.Run("MustQuery", testParamsMustQueryXXX)
	t.Run("Query", testParamsQueryXXX)
	t.Run("MustParams", testParamsMustParamsXXX)
	t.Run("Params", testParamsParamsXXX)
}

func testParamsMustFormXXX(t *testing.T) {
	tests := []struct {
		Values         url.Values
		ExpectedErr    error
		ExpectedOutput interface{}
		Invoke         func(r *http.Request, key string) interface{}
	}{
		// -- MustFormString --
		// #0
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: "",
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormString(r, key)
			},
		},
		// #1
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: "",
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormString(r, key)
			},
		},
		// #2
		{
			Values:         url.Values{"key": []string{"Rob"}},
			ExpectedErr:    nil,
			ExpectedOutput: "Rob",
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormString(r, key)
			},
		},
		// -- MustFormBool --
		// #3
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormBool(r, key)
			},
		},
		// #4
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormBool(r, key)
			},
		},
		// #5
		{
			Values:         url.Values{"key": []string{"true"}},
			ExpectedErr:    nil,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormBool(r, key)
			},
		},
		// #6
		{
			Values:         url.Values{"key": []string{"false"}},
			ExpectedErr:    nil,
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormBool(r, key)
			},
		},
		// #7
		{
			Values:         url.Values{"key": []string{"1"}},
			ExpectedErr:    nil,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormBool(r, key)
			},
		},
		// #8
		{
			Values:         url.Values{"key": []string{"0"}},
			ExpectedErr:    nil,
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormBool(r, key)
			},
		},
		// -- MustFormInt --
		// #9
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: 0,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt(r, key)
			},
		},
		// #10
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: 0,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt(r, key)
			},
		},
		// #11
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt(r, key)
			},
		},
		// -- MustFormInt32 --
		// #12
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt32(r, key)
			},
		},
		// #13
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt32(r, key)
			},
		},
		// #14
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt32(r, key)
			},
		},
		// -- MustFormInt64 --
		// #15
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt64(r, key)
			},
		},
		// #16
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt64(r, key)
			},
		},
		// #17
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormInt64(r, key)
			},
		},
		// -- MustFormFloat32 --
		// #18
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormFloat32(r, key)
			},
		},
		// #19
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormFloat32(r, key)
			},
		},
		// #20
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: float32(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormFloat32(r, key)
			},
		},
		// -- MustFormFloat64 --
		// #21
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormFloat64(r, key)
			},
		},
		// #22
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormFloat64(r, key)
			},
		},
		// #23
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: float64(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormFloat64(r, key)
			},
		},
		// -- MustFormTime --
		// #24
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTime(r, key, time.RFC3339)
			},
		},
		// #25
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTime(r, key, time.RFC3339)
			},
		},
		// #26
		{
			Values:         url.Values{"key": []string{"2018-12-31T18:47:59.999999999Z"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 18, 47, 59, 999999999, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTime(r, key, time.RFC3339)
			},
		},
		// #27
		{
			Values:         url.Values{"key": []string{"2018-12-31"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTime(r, key, "2006-01-02")
			},
		},
		// -- MustFormTimeWithDefault --
		// #28
		{
			Values:         url.Values{},
			ExpectedErr:    nil,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTimeWithDefault(r, key, time.RFC3339, defaultTime)
			},
		},
		// #29
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    nil,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTimeWithDefault(r, key, time.RFC3339, defaultTime)
			},
		},
		// #30
		{
			Values:         url.Values{"key": []string{"2018-12-31T18:47:59.999999999Z"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 18, 47, 59, 999999999, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTimeWithDefault(r, key, time.RFC3339Nano, defaultTime)
			},
		},
		// #31
		{
			Values:         url.Values{"key": []string{"2018-12-31"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormTimeWithDefault(r, key, "2006-01-02", defaultTime)
			},
		},
		// -- MustFormDuration --
		// #32
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormDuration(r, key)
			},
		},
		// #33
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormDuration(r, key)
			},
		},
		// #34
		{
			Values:         url.Values{"key": []string{"1m12s"}},
			ExpectedErr:    nil,
			ExpectedOutput: 72 * time.Second,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormDuration(r, key)
			},
		},
		// -- MustFormDurationWithDefault --
		// #35
		{
			Values:         url.Values{},
			ExpectedErr:    nil,
			ExpectedOutput: defaultDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormDurationWithDefault(r, key, defaultDuration)
			},
		},
		// #36
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    nil,
			ExpectedOutput: defaultDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormDurationWithDefault(r, key, defaultDuration)
			},
		},
		// #37
		{
			Values:         url.Values{"key": []string{"1m12s"}},
			ExpectedErr:    nil,
			ExpectedOutput: 72 * time.Second,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustFormDurationWithDefault(r, key, defaultDuration)
			},
		},
	}
	for i, tt := range tests {
		var err error
		var got interface{}
		f := func(req *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					e, ok := r.(error)
					if !ok {
						t.Fatalf("expected an error type, got %T", r)
					}
					err = e
				}
			}()
			got = tt.Invoke(req, "key")
			//got = MustFormString(req, "key")
		}
		req, _ := http.NewRequest("GET", "/", http.NoBody)
		req.Form = tt.Values
		f(req)
		if err != nil {
			if tt.ExpectedErr != nil {
				if err.Error() != tt.ExpectedErr.Error() {
					t.Fatalf("#%d. expected Error = %v; got %v", i, tt.ExpectedErr, err)
				}
			} else {
				t.Fatalf("#%d. expected no error; got %v", i, err)
			}
		} else {
			if tt.ExpectedErr != nil {
				t.Fatalf("#%d. expected error %v; got %v", i, tt.ExpectedErr, err)
			} else {
				if got != tt.ExpectedOutput {
					t.Fatalf("#%d. expected Output = %v; got %v", i, tt.ExpectedOutput, got)
				}
			}
		}
	}
}

func testParamsFormXXX(t *testing.T) {
	tests := []struct {
		Values         url.Values
		DefaultValue   interface{}
		ExpectedOutput interface{}
		Invoke         func(r *http.Request, key string, defaultValue interface{}) interface{}
	}{
		// -- FormString --
		// #0
		{
			Values:         url.Values{},
			DefaultValue:   "",
			ExpectedOutput: "",
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormString(r, key, defaultValue.(string))
			},
		},
		// #1
		{
			Values:         url.Values{"key": []string{""}},
			DefaultValue:   "",
			ExpectedOutput: "",
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormString(r, key, defaultValue.(string))
			},
		},
		// #2
		{
			Values:         url.Values{},
			DefaultValue:   "Mary",
			ExpectedOutput: "Mary",
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormString(r, key, defaultValue.(string))
			},
		},
		// #3
		{
			Values:         url.Values{"key": []string{""}},
			DefaultValue:   "Mary",
			ExpectedOutput: "Mary",
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormString(r, key, defaultValue.(string))
			},
		},
		// -- FormBool --
		// #4
		{
			Values:         url.Values{},
			DefaultValue:   false,
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormBool(r, key, defaultValue.(bool))
			},
		},
		// #5
		{
			Values:         url.Values{"key": []string{"true"}},
			DefaultValue:   false,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormBool(r, key, defaultValue.(bool))
			},
		},
		// #6
		{
			Values:         url.Values{},
			DefaultValue:   true,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormBool(r, key, defaultValue.(bool))
			},
		},
		// #7
		{
			Values:         url.Values{"key": []string{"true"}},
			DefaultValue:   false,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormBool(r, key, defaultValue.(bool))
			},
		},
		// #8
		{
			Values:         url.Values{"key": []string{"invalid-bool"}},
			DefaultValue:   true,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormBool(r, key, defaultValue.(bool))
			},
		},
		// -- FormInt --
		// #9
		{
			Values:         url.Values{},
			DefaultValue:   0,
			ExpectedOutput: 0,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt(r, key, defaultValue.(int))
			},
		},
		// #10
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   0,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt(r, key, defaultValue.(int))
			},
		},
		// #11
		{
			Values:         url.Values{},
			DefaultValue:   42,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt(r, key, defaultValue.(int))
			},
		},
		// #12
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   87,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt(r, key, defaultValue.(int))
			},
		},
		// #13
		{
			Values:         url.Values{"key": []string{"invalid-int"}},
			DefaultValue:   42,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt(r, key, defaultValue.(int))
			},
		},
		// -- FormInt32 --
		// #14
		{
			Values:         url.Values{},
			DefaultValue:   int32(0),
			ExpectedOutput: int32(0),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt32(r, key, defaultValue.(int32))
			},
		},
		// #15
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   int32(0),
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt32(r, key, defaultValue.(int32))
			},
		},
		// #16
		{
			Values:         url.Values{},
			DefaultValue:   int32(42),
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt32(r, key, defaultValue.(int32))
			},
		},
		// #17
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   int32(87),
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt32(r, key, defaultValue.(int32))
			},
		},
		// #18
		{
			Values:         url.Values{"key": []string{"invalid-int"}},
			DefaultValue:   int32(42),
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt32(r, key, defaultValue.(int32))
			},
		},
		// -- FormInt64 --
		// #19
		{
			Values:         url.Values{},
			DefaultValue:   int64(0),
			ExpectedOutput: int64(0),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt64(r, key, defaultValue.(int64))
			},
		},
		// #20
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   int64(0),
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt64(r, key, defaultValue.(int64))
			},
		},
		// #21
		{
			Values:         url.Values{},
			DefaultValue:   int64(42),
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt64(r, key, defaultValue.(int64))
			},
		},
		// #22
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   int64(87),
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt64(r, key, defaultValue.(int64))
			},
		},
		// #23
		{
			Values:         url.Values{"key": []string{"invalid-int"}},
			DefaultValue:   int64(42),
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormInt64(r, key, defaultValue.(int64))
			},
		},
		// -- FormFloat32 --
		// #24
		{
			Values:         url.Values{},
			DefaultValue:   float32(0),
			ExpectedOutput: float32(0),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormFloat32(r, key, defaultValue.(float32))
			},
		},
		// #25
		{
			Values:         url.Values{"key": []string{"42.7"}},
			DefaultValue:   float32(0),
			ExpectedOutput: float32(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormFloat32(r, key, defaultValue.(float32))
			},
		},
		// #26
		{
			Values:         url.Values{},
			DefaultValue:   float32(42),
			ExpectedOutput: float32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormFloat32(r, key, defaultValue.(float32))
			},
		},
		// #27
		{
			Values:         url.Values{"key": []string{"42.7"}},
			DefaultValue:   float32(87.3),
			ExpectedOutput: float32(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormFloat32(r, key, defaultValue.(float32))
			},
		},
		// #28
		{
			Values:         url.Values{"key": []string{"invalid-value"}},
			DefaultValue:   float32(42.7),
			ExpectedOutput: float32(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormFloat32(r, key, defaultValue.(float32))
			},
		},
		// -- FormFloat64 --
		// #29
		{
			Values:         url.Values{},
			DefaultValue:   float64(0),
			ExpectedOutput: float64(0),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormFloat64(r, key, defaultValue.(float64))
			},
		},
		// #30
		{
			Values:         url.Values{"key": []string{"42.7"}},
			DefaultValue:   float64(0),
			ExpectedOutput: float64(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormFloat64(r, key, defaultValue.(float64))
			},
		},
		// #31
		{
			Values:         url.Values{},
			DefaultValue:   float64(42),
			ExpectedOutput: float64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormFloat64(r, key, defaultValue.(float64))
			},
		},
		// #32
		{
			Values:         url.Values{"key": []string{"42.7"}},
			DefaultValue:   float64(87.3),
			ExpectedOutput: float64(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormFloat64(r, key, defaultValue.(float64))
			},
		},
		// #33
		{
			Values:         url.Values{"key": []string{"invalid-value"}},
			DefaultValue:   float64(42.7),
			ExpectedOutput: float64(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormFloat64(r, key, defaultValue.(float64))
			},
		},
		// -- FormTime --
		// #34
		{
			Values:         url.Values{},
			DefaultValue:   zeroTime,
			ExpectedOutput: zeroTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormTime(r, key, time.RFC3339, defaultValue.(time.Time))
			},
		},
		// #35
		{
			Values:         url.Values{"key": []string{defaultTime.Truncate(time.Second).Format(time.RFC3339)}},
			DefaultValue:   zeroTime,
			ExpectedOutput: defaultTime.Truncate(time.Second),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormTime(r, key, time.RFC3339, defaultValue.(time.Time))
			},
		},
		// #36
		{
			Values:         url.Values{},
			DefaultValue:   defaultTime,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormTime(r, key, time.RFC3339, defaultValue.(time.Time))
			},
		},
		// #37
		{
			Values:         url.Values{"key": []string{defaultTime.Format(time.RFC3339Nano)}},
			DefaultValue:   zeroTime,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormTime(r, key, time.RFC3339Nano, defaultValue.(time.Time))
			},
		},
		// #38
		{
			Values:         url.Values{"key": []string{"invalid-value"}},
			DefaultValue:   defaultTime,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormTime(r, key, time.RFC3339, defaultValue.(time.Time))
			},
		},
		// #39
		{
			Values:         url.Values{"key": []string{"2018-12-31"}},
			DefaultValue:   defaultTime,
			ExpectedOutput: time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormTime(r, key, "2006-01-02", defaultValue.(time.Time))
			},
		},
		// -- FormDuration --
		// #34
		{
			Values:         url.Values{},
			DefaultValue:   zeroDuration,
			ExpectedOutput: zeroDuration,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormDuration(r, key, defaultValue.(time.Duration))
			},
		},
		// #35
		{
			Values:         url.Values{"key": []string{"1m12s"}},
			DefaultValue:   zeroDuration,
			ExpectedOutput: 72 * time.Second,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormDuration(r, key, defaultValue.(time.Duration))
			},
		},
		// #36
		{
			Values:         url.Values{},
			DefaultValue:   12 * time.Second,
			ExpectedOutput: 12 * time.Second,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormDuration(r, key, defaultValue.(time.Duration))
			},
		},
		// #37
		{
			Values:         url.Values{"key": []string{"invalid-value"}},
			DefaultValue:   12 * time.Second,
			ExpectedOutput: 12 * time.Second,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return FormDuration(r, key, defaultValue.(time.Duration))
			},
		},
	}
	for i, tt := range tests {
		req, _ := http.NewRequest("GET", "/", http.NoBody)
		req.Form = tt.Values
		got := tt.Invoke(req, "key", tt.DefaultValue)
		if !cmp.Equal(got, tt.ExpectedOutput) {
			t.Fatalf("#%d. expected Output = %v; got %v", i, tt.ExpectedOutput, got)
		}
	}
}

func testParamsMustQueryXXX(t *testing.T) {
	tests := []struct {
		Values         url.Values
		ExpectedErr    error
		ExpectedOutput interface{}
		Invoke         func(r *http.Request, key string) interface{}
	}{
		// -- MustQueryString --
		// #0
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: "",
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryString(r, key)
			},
		},
		// #1
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: "",
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryString(r, key)
			},
		},
		// #2
		{
			Values:         url.Values{"key": []string{"Rob"}},
			ExpectedErr:    nil,
			ExpectedOutput: "Rob",
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryString(r, key)
			},
		},
		// -- MustQueryBool --
		// #3
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryBool(r, key)
			},
		},
		// #4
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryBool(r, key)
			},
		},
		// #5
		{
			Values:         url.Values{"key": []string{"true"}},
			ExpectedErr:    nil,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryBool(r, key)
			},
		},
		// #6
		{
			Values:         url.Values{"key": []string{"false"}},
			ExpectedErr:    nil,
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryBool(r, key)
			},
		},
		// #7
		{
			Values:         url.Values{"key": []string{"1"}},
			ExpectedErr:    nil,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryBool(r, key)
			},
		},
		// #8
		{
			Values:         url.Values{"key": []string{"0"}},
			ExpectedErr:    nil,
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryBool(r, key)
			},
		},
		// -- MustQueryInt --
		// #9
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: 0,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryInt(r, key)
			},
		},
		// #10
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: 0,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryInt(r, key)
			},
		},
		// #11
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryInt(r, key)
			},
		},
		// -- MustQueryInt32 --
		// #12
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryInt32(r, key)
			},
		},
		// #13
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryInt32(r, key)
			},
		},
		// #14
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryInt32(r, key)
			},
		},
		// -- MustQueryInt64 --
		// #15
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryInt64(r, key)
			},
		},
		// #16
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryInt64(r, key)
			},
		},
		// #17
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryInt64(r, key)
			},
		},
		// -- MustQueryFloat32 --
		// #18
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryFloat32(r, key)
			},
		},
		// #19
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryFloat32(r, key)
			},
		},
		// #20
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: float32(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryFloat32(r, key)
			},
		},
		// -- MustQueryFloat64 --
		// #21
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryFloat64(r, key)
			},
		},
		// #22
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryFloat64(r, key)
			},
		},
		// #23
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: float64(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryFloat64(r, key)
			},
		},
		// -- MustQueryTime --
		// #24
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryTime(r, key, time.RFC3339)
			},
		},
		// #25
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryTime(r, key, time.RFC3339)
			},
		},
		// #26
		{
			Values:         url.Values{"key": []string{"2018-12-31T18:47:59.999999999Z"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 18, 47, 59, 999999999, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryTime(r, key, time.RFC3339)
			},
		},
		// #27
		{
			Values:         url.Values{"key": []string{"2018-12-31"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryTime(r, key, "2006-01-02")
			},
		},
		// -- MustQueryTimeWithDefault --
		// #28
		{
			Values:         url.Values{},
			ExpectedErr:    nil,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryTimeWithDefault(r, key, time.RFC3339, defaultTime)
			},
		},
		// #29
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    nil,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryTimeWithDefault(r, key, time.RFC3339, defaultTime)
			},
		},
		// #30
		{
			Values:         url.Values{"key": []string{"2018-12-31T18:47:59.999999999Z"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 18, 47, 59, 999999999, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryTimeWithDefault(r, key, time.RFC3339Nano, defaultTime)
			},
		},
		// #31
		{
			Values:         url.Values{"key": []string{"2018-12-31"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryTimeWithDefault(r, key, "2006-01-02", defaultTime)
			},
		},
		// -- MustQueryDuration --
		// #32
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryDuration(r, key)
			},
		},
		// #33
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryDuration(r, key)
			},
		},
		// #34
		{
			Values:         url.Values{"key": []string{"1m12s"}},
			ExpectedErr:    nil,
			ExpectedOutput: 72 * time.Second,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryDuration(r, key)
			},
		},
		// -- MustQueryDurationWithDefault --
		// #35
		{
			Values:         url.Values{},
			ExpectedErr:    nil,
			ExpectedOutput: defaultDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryDurationWithDefault(r, key, defaultDuration)
			},
		},
		// #36
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    nil,
			ExpectedOutput: defaultDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryDurationWithDefault(r, key, defaultDuration)
			},
		},
		// #37
		{
			Values:         url.Values{"key": []string{"1m12s"}},
			ExpectedErr:    nil,
			ExpectedOutput: 72 * time.Second,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustQueryDurationWithDefault(r, key, defaultDuration)
			},
		},
	}
	for i, tt := range tests {
		var err error
		var got interface{}
		f := func(req *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					e, ok := r.(error)
					if !ok {
						t.Fatalf("expected an error type, got %T", r)
					}
					err = e
				}
			}()
			got = tt.Invoke(req, "key")
		}
		req, _ := http.NewRequest("GET", "/?"+tt.Values.Encode(), http.NoBody)
		f(req)
		if err != nil {
			if tt.ExpectedErr != nil {
				if err.Error() != tt.ExpectedErr.Error() {
					t.Fatalf("#%d. expected Error = %v; got %v", i, tt.ExpectedErr, err)
				}
			} else {
				t.Fatalf("#%d. expected no error; got %v", i, err)
			}
		} else {
			if tt.ExpectedErr != nil {
				t.Fatalf("#%d. expected error %v; got %v", i, tt.ExpectedErr, err)
			} else {
				if got != tt.ExpectedOutput {
					t.Fatalf("#%d. expected Output = %v; got %v", i, tt.ExpectedOutput, got)
				}
			}
		}
	}
}

func testParamsQueryXXX(t *testing.T) {
	tests := []struct {
		Values         url.Values
		DefaultValue   interface{}
		ExpectedOutput interface{}
		Invoke         func(r *http.Request, key string, defaultValue interface{}) interface{}
	}{
		// -- QueryString --
		// #0
		{
			Values:         url.Values{},
			DefaultValue:   "",
			ExpectedOutput: "",
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryString(r, key, defaultValue.(string))
			},
		},
		// #1
		{
			Values:         url.Values{"key": []string{""}},
			DefaultValue:   "",
			ExpectedOutput: "",
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryString(r, key, defaultValue.(string))
			},
		},
		// #2
		{
			Values:         url.Values{},
			DefaultValue:   "Mary",
			ExpectedOutput: "Mary",
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryString(r, key, defaultValue.(string))
			},
		},
		// #3
		{
			Values:         url.Values{"key": []string{""}},
			DefaultValue:   "Mary",
			ExpectedOutput: "Mary",
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryString(r, key, defaultValue.(string))
			},
		},
		// -- QueryBool --
		// #4
		{
			Values:         url.Values{},
			DefaultValue:   false,
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryBool(r, key, defaultValue.(bool))
			},
		},
		// #5
		{
			Values:         url.Values{"key": []string{"true"}},
			DefaultValue:   false,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryBool(r, key, defaultValue.(bool))
			},
		},
		// #6
		{
			Values:         url.Values{},
			DefaultValue:   true,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryBool(r, key, defaultValue.(bool))
			},
		},
		// #7
		{
			Values:         url.Values{"key": []string{"true"}},
			DefaultValue:   false,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryBool(r, key, defaultValue.(bool))
			},
		},
		// #8
		{
			Values:         url.Values{"key": []string{"invalid-bool"}},
			DefaultValue:   true,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryBool(r, key, defaultValue.(bool))
			},
		},
		// -- QueryInt --
		// #9
		{
			Values:         url.Values{},
			DefaultValue:   0,
			ExpectedOutput: 0,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt(r, key, defaultValue.(int))
			},
		},
		// #10
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   0,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt(r, key, defaultValue.(int))
			},
		},
		// #11
		{
			Values:         url.Values{},
			DefaultValue:   42,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt(r, key, defaultValue.(int))
			},
		},
		// #12
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   87,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt(r, key, defaultValue.(int))
			},
		},
		// #13
		{
			Values:         url.Values{"key": []string{"invalid-int"}},
			DefaultValue:   42,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt(r, key, defaultValue.(int))
			},
		},
		// -- QueryInt32 --
		// #14
		{
			Values:         url.Values{},
			DefaultValue:   int32(0),
			ExpectedOutput: int32(0),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt32(r, key, defaultValue.(int32))
			},
		},
		// #15
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   int32(0),
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt32(r, key, defaultValue.(int32))
			},
		},
		// #16
		{
			Values:         url.Values{},
			DefaultValue:   int32(42),
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt32(r, key, defaultValue.(int32))
			},
		},
		// #17
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   int32(87),
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt32(r, key, defaultValue.(int32))
			},
		},
		// #18
		{
			Values:         url.Values{"key": []string{"invalid-int"}},
			DefaultValue:   int32(42),
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt32(r, key, defaultValue.(int32))
			},
		},
		// -- QueryInt64 --
		// #19
		{
			Values:         url.Values{},
			DefaultValue:   int64(0),
			ExpectedOutput: int64(0),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt64(r, key, defaultValue.(int64))
			},
		},
		// #20
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   int64(0),
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt64(r, key, defaultValue.(int64))
			},
		},
		// #21
		{
			Values:         url.Values{},
			DefaultValue:   int64(42),
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt64(r, key, defaultValue.(int64))
			},
		},
		// #22
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   int64(87),
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt64(r, key, defaultValue.(int64))
			},
		},
		// #23
		{
			Values:         url.Values{"key": []string{"invalid-int"}},
			DefaultValue:   int64(42),
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryInt64(r, key, defaultValue.(int64))
			},
		},
		// -- QueryFloat32 --
		// #24
		{
			Values:         url.Values{},
			DefaultValue:   float32(0),
			ExpectedOutput: float32(0),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryFloat32(r, key, defaultValue.(float32))
			},
		},
		// #25
		{
			Values:         url.Values{"key": []string{"42.7"}},
			DefaultValue:   float32(0),
			ExpectedOutput: float32(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryFloat32(r, key, defaultValue.(float32))
			},
		},
		// #26
		{
			Values:         url.Values{},
			DefaultValue:   float32(42),
			ExpectedOutput: float32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryFloat32(r, key, defaultValue.(float32))
			},
		},
		// #27
		{
			Values:         url.Values{"key": []string{"42.7"}},
			DefaultValue:   float32(87.3),
			ExpectedOutput: float32(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryFloat32(r, key, defaultValue.(float32))
			},
		},
		// #28
		{
			Values:         url.Values{"key": []string{"invalid-value"}},
			DefaultValue:   float32(42.7),
			ExpectedOutput: float32(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryFloat32(r, key, defaultValue.(float32))
			},
		},
		// -- QueryFloat64 --
		// #29
		{
			Values:         url.Values{},
			DefaultValue:   float64(0),
			ExpectedOutput: float64(0),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryFloat64(r, key, defaultValue.(float64))
			},
		},
		// #30
		{
			Values:         url.Values{"key": []string{"42.7"}},
			DefaultValue:   float64(0),
			ExpectedOutput: float64(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryFloat64(r, key, defaultValue.(float64))
			},
		},
		// #31
		{
			Values:         url.Values{},
			DefaultValue:   float64(42),
			ExpectedOutput: float64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryFloat64(r, key, defaultValue.(float64))
			},
		},
		// #32
		{
			Values:         url.Values{"key": []string{"42.7"}},
			DefaultValue:   float64(87.3),
			ExpectedOutput: float64(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryFloat64(r, key, defaultValue.(float64))
			},
		},
		// #33
		{
			Values:         url.Values{"key": []string{"invalid-value"}},
			DefaultValue:   float64(42.7),
			ExpectedOutput: float64(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryFloat64(r, key, defaultValue.(float64))
			},
		},
		// -- QueryTime --
		// #34
		{
			Values:         url.Values{},
			ExpectedOutput: zeroTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryTime(r, key, time.RFC3339)
			},
		},
		// #35
		{
			Values:         url.Values{"key": []string{defaultTime.Truncate(time.Second).Format(time.RFC3339)}},
			ExpectedOutput: defaultTime.Truncate(time.Second),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryTime(r, key, time.RFC3339)
			},
		},
		// #36
		{
			Values:         url.Values{},
			ExpectedOutput: zeroTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryTime(r, key, time.RFC3339)
			},
		},
		// #37
		{
			Values:         url.Values{"key": []string{defaultTime.Format(time.RFC3339Nano)}},
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryTime(r, key, time.RFC3339Nano)
			},
		},
		// #38
		{
			Values:         url.Values{"key": []string{"invalid-value"}},
			ExpectedOutput: zeroTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryTime(r, key, time.RFC3339)
			},
		},
		// #39
		{
			Values:         url.Values{"key": []string{"2018-12-31"}},
			ExpectedOutput: time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryTime(r, key, "2006-01-02")
			},
		},
		// -- QueryTimeWithDefault --
		// #40
		{
			Values:         url.Values{},
			DefaultValue:   zeroTime,
			ExpectedOutput: zeroTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryTimeWithDefault(r, key, time.RFC3339, defaultValue.(time.Time))
			},
		},
		// #41
		{
			Values:         url.Values{"key": []string{defaultTime.Truncate(time.Second).Format(time.RFC3339)}},
			DefaultValue:   zeroTime,
			ExpectedOutput: defaultTime.Truncate(time.Second),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryTimeWithDefault(r, key, time.RFC3339, defaultValue.(time.Time))
			},
		},
		// #42
		{
			Values:         url.Values{},
			DefaultValue:   defaultTime,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryTimeWithDefault(r, key, time.RFC3339, defaultValue.(time.Time))
			},
		},
		// #43
		{
			Values:         url.Values{"key": []string{defaultTime.Format(time.RFC3339Nano)}},
			DefaultValue:   zeroTime,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryTimeWithDefault(r, key, time.RFC3339Nano, defaultValue.(time.Time))
			},
		},
		// #44
		{
			Values:         url.Values{"key": []string{"invalid-value"}},
			DefaultValue:   defaultTime,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryTimeWithDefault(r, key, time.RFC3339, defaultValue.(time.Time))
			},
		},
		// #45
		{
			Values:         url.Values{"key": []string{"2018-12-31"}},
			DefaultValue:   defaultTime,
			ExpectedOutput: time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryTimeWithDefault(r, key, "2006-01-02", defaultValue.(time.Time))
			},
		},
		// -- QueryDuration --
		// #46
		{
			Values:         url.Values{},
			ExpectedOutput: zeroDuration,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryDuration(r, key)
			},
		},
		// #47
		{
			Values:         url.Values{"key": []string{"1m12s"}},
			ExpectedOutput: 72 * time.Second,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryDuration(r, key)
			},
		},
		// #48
		{
			Values:         url.Values{"key": []string{"invalid-value"}},
			ExpectedOutput: zeroDuration,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryDuration(r, key)
			},
		},
		// -- QueryDurationWithDefault --
		// #49
		{
			Values:         url.Values{},
			DefaultValue:   zeroDuration,
			ExpectedOutput: zeroDuration,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryDurationWithDefault(r, key, defaultValue.(time.Duration))
			},
		},
		// #50
		{
			Values:         url.Values{"key": []string{"1m12s"}},
			DefaultValue:   zeroDuration,
			ExpectedOutput: 72 * time.Second,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryDurationWithDefault(r, key, defaultValue.(time.Duration))
			},
		},
		// #51
		{
			Values:         url.Values{},
			DefaultValue:   12 * time.Second,
			ExpectedOutput: 12 * time.Second,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryDurationWithDefault(r, key, defaultValue.(time.Duration))
			},
		},
		// #52
		{
			Values:         url.Values{"key": []string{"invalid-value"}},
			DefaultValue:   12 * time.Second,
			ExpectedOutput: 12 * time.Second,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return QueryDurationWithDefault(r, key, defaultValue.(time.Duration))
			},
		},
	}
	for i, tt := range tests {
		req, _ := http.NewRequest("GET", "/?"+tt.Values.Encode(), http.NoBody)
		got := tt.Invoke(req, "key", tt.DefaultValue)
		if !cmp.Equal(got, tt.ExpectedOutput) {
			t.Fatalf("#%d. expected Output = %v; got %v", i, tt.ExpectedOutput, got)
		}
	}
}

func testParamsMustParamsXXX(t *testing.T) {
	tests := []struct {
		Values         url.Values
		ExpectedErr    error
		ExpectedOutput interface{}
		Invoke         func(r *http.Request, key string) interface{}
	}{
		// -- MustParamsString --
		// #0
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: "",
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsString(r, key)
			},
		},
		// #1
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: "",
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsString(r, key)
			},
		},
		// #2
		{
			Values:         url.Values{"key": []string{"Rob"}},
			ExpectedErr:    nil,
			ExpectedOutput: "Rob",
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsString(r, key)
			},
		},
		// -- MustParamsBool --
		// #3
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsBool(r, key)
			},
		},
		// #4
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsBool(r, key)
			},
		},
		// #5
		{
			Values:         url.Values{"key": []string{"true"}},
			ExpectedErr:    nil,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsBool(r, key)
			},
		},
		// #6
		{
			Values:         url.Values{"key": []string{"false"}},
			ExpectedErr:    nil,
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsBool(r, key)
			},
		},
		// #7
		{
			Values:         url.Values{"key": []string{"1"}},
			ExpectedErr:    nil,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsBool(r, key)
			},
		},
		// #8
		{
			Values:         url.Values{"key": []string{"0"}},
			ExpectedErr:    nil,
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsBool(r, key)
			},
		},
		// -- MustParamsInt --
		// #9
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: 0,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsInt(r, key)
			},
		},
		// #10
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: 0,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsInt(r, key)
			},
		},
		// #11
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsInt(r, key)
			},
		},
		// -- MustParamsInt32 --
		// #12
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsInt32(r, key)
			},
		},
		// #13
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsInt32(r, key)
			},
		},
		// #14
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsInt32(r, key)
			},
		},
		// -- MustParamsInt64 --
		// #15
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsInt64(r, key)
			},
		},
		// #16
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: int64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsInt64(r, key)
			},
		},
		// #17
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsInt64(r, key)
			},
		},
		// -- MustParamsFloat32 --
		// #18
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsFloat32(r, key)
			},
		},
		// #19
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float32(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsFloat32(r, key)
			},
		},
		// #20
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: float32(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsFloat32(r, key)
			},
		},
		// -- MustParamsFloat64 --
		// #21
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsFloat64(r, key)
			},
		},
		// #22
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: float64(0),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsFloat64(r, key)
			},
		},
		// #23
		{
			Values:         url.Values{"key": []string{"42"}},
			ExpectedErr:    nil,
			ExpectedOutput: float64(42),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsFloat64(r, key)
			},
		},
		// -- MustParamsTime --
		// #24
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsTime(r, key, time.RFC3339)
			},
		},
		// #25
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsTime(r, key, time.RFC3339)
			},
		},
		// #26
		{
			Values:         url.Values{"key": []string{"2018-12-31T18:47:59.999999999Z"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 18, 47, 59, 999999999, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsTime(r, key, time.RFC3339)
			},
		},
		// #27
		{
			Values:         url.Values{"key": []string{"2018-12-31"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsTime(r, key, "2006-01-02")
			},
		},
		// -- MustParamsTimeWithDefault --
		// #28
		{
			Values:         url.Values{},
			ExpectedErr:    nil,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsTimeWithDefault(r, key, time.RFC3339, defaultTime)
			},
		},
		// #29
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    nil,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsTimeWithDefault(r, key, time.RFC3339, defaultTime)
			},
		},
		// #30
		{
			Values:         url.Values{"key": []string{"2018-12-31T18:47:59.999999999Z"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 18, 47, 59, 999999999, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsTimeWithDefault(r, key, time.RFC3339Nano, defaultTime)
			},
		},
		// #31
		{
			Values:         url.Values{"key": []string{"2018-12-31"}},
			ExpectedErr:    nil,
			ExpectedOutput: time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC),
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsTimeWithDefault(r, key, "2006-01-02", defaultTime)
			},
		},
		// -- MustParamsDuration --
		// #32
		{
			Values:         url.Values{},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsDuration(r, key)
			},
		},
		// #33
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    BadRequestError{Message: `Missing parameter "key"`},
			ExpectedOutput: zeroDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsDuration(r, key)
			},
		},
		// #34
		{
			Values:         url.Values{"key": []string{"1m12s"}},
			ExpectedErr:    nil,
			ExpectedOutput: 72 * time.Second,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsDuration(r, key)
			},
		},
		// -- MustParamsDurationWithDefault --
		// #35
		{
			Values:         url.Values{},
			ExpectedErr:    nil,
			ExpectedOutput: defaultDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsDurationWithDefault(r, key, defaultDuration)
			},
		},
		// #36
		{
			Values:         url.Values{"key": []string{""}},
			ExpectedErr:    nil,
			ExpectedOutput: defaultDuration,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsDurationWithDefault(r, key, defaultDuration)
			},
		},
		// #37
		{
			Values:         url.Values{"key": []string{"1m12s"}},
			ExpectedErr:    nil,
			ExpectedOutput: 72 * time.Second,
			Invoke: func(r *http.Request, key string) interface{} {
				return MustParamsDurationWithDefault(r, key, defaultDuration)
			},
		},
	}
	for i, tt := range tests {
		var err error
		var got interface{}
		f := func(req *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					e, ok := r.(error)
					if !ok {
						t.Fatalf("expected an error type, got %T", r)
					}
					err = e
				}
			}()
			got = tt.Invoke(req, "key")
		}
		req, _ := http.NewRequest("GET", "/", http.NoBody)
		vars := make(map[string]string)
		for k, v := range tt.Values {
			if len(v) > 0 {
				vars[k] = v[0]
			} else {
				vars[k] = ""
			}
		}
		req = mux.SetURLVars(req, vars)
		f(req)
		if err != nil {
			if tt.ExpectedErr != nil {
				if err.Error() != tt.ExpectedErr.Error() {
					t.Fatalf("#%d. expected Error = %v; got %v", i, tt.ExpectedErr, err)
				}
			} else {
				t.Fatalf("#%d. expected no error; got %v", i, err)
			}
		} else {
			if tt.ExpectedErr != nil {
				t.Fatalf("#%d. expected error %v; got %v", i, tt.ExpectedErr, err)
			} else {
				if got != tt.ExpectedOutput {
					t.Fatalf("#%d. expected Output = %v; got %v", i, tt.ExpectedOutput, got)
				}
			}
		}
	}
}

func testParamsParamsXXX(t *testing.T) {
	tests := []struct {
		Values         url.Values
		DefaultValue   interface{}
		ExpectedOutput interface{}
		Invoke         func(r *http.Request, key string, defaultValue interface{}) interface{}
	}{
		// -- ParamsString --
		// #0
		{
			Values:         url.Values{},
			DefaultValue:   "",
			ExpectedOutput: "",
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsString(r, key, defaultValue.(string))
			},
		},
		// #1
		{
			Values:         url.Values{"key": []string{""}},
			DefaultValue:   "",
			ExpectedOutput: "",
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsString(r, key, defaultValue.(string))
			},
		},
		// #2
		{
			Values:         url.Values{},
			DefaultValue:   "Mary",
			ExpectedOutput: "Mary",
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsString(r, key, defaultValue.(string))
			},
		},
		// #3
		{
			Values:         url.Values{"key": []string{""}},
			DefaultValue:   "Mary",
			ExpectedOutput: "Mary",
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsString(r, key, defaultValue.(string))
			},
		},
		// -- ParamsBool --
		// #4
		{
			Values:         url.Values{},
			DefaultValue:   false,
			ExpectedOutput: false,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsBool(r, key, defaultValue.(bool))
			},
		},
		// #5
		{
			Values:         url.Values{"key": []string{"true"}},
			DefaultValue:   false,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsBool(r, key, defaultValue.(bool))
			},
		},
		// #6
		{
			Values:         url.Values{},
			DefaultValue:   true,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsBool(r, key, defaultValue.(bool))
			},
		},
		// #7
		{
			Values:         url.Values{"key": []string{"true"}},
			DefaultValue:   false,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsBool(r, key, defaultValue.(bool))
			},
		},
		// #8
		{
			Values:         url.Values{"key": []string{"invalid-bool"}},
			DefaultValue:   true,
			ExpectedOutput: true,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsBool(r, key, defaultValue.(bool))
			},
		},
		// -- ParamsInt --
		// #9
		{
			Values:         url.Values{},
			DefaultValue:   0,
			ExpectedOutput: 0,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt(r, key, defaultValue.(int))
			},
		},
		// #10
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   0,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt(r, key, defaultValue.(int))
			},
		},
		// #11
		{
			Values:         url.Values{},
			DefaultValue:   42,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt(r, key, defaultValue.(int))
			},
		},
		// #12
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   87,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt(r, key, defaultValue.(int))
			},
		},
		// #13
		{
			Values:         url.Values{"key": []string{"invalid-int"}},
			DefaultValue:   42,
			ExpectedOutput: 42,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt(r, key, defaultValue.(int))
			},
		},
		// -- ParamsInt32 --
		// #14
		{
			Values:         url.Values{},
			DefaultValue:   int32(0),
			ExpectedOutput: int32(0),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt32(r, key, defaultValue.(int32))
			},
		},
		// #15
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   int32(0),
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt32(r, key, defaultValue.(int32))
			},
		},
		// #16
		{
			Values:         url.Values{},
			DefaultValue:   int32(42),
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt32(r, key, defaultValue.(int32))
			},
		},
		// #17
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   int32(87),
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt32(r, key, defaultValue.(int32))
			},
		},
		// #18
		{
			Values:         url.Values{"key": []string{"invalid-int"}},
			DefaultValue:   int32(42),
			ExpectedOutput: int32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt32(r, key, defaultValue.(int32))
			},
		},
		// -- ParamsInt64 --
		// #19
		{
			Values:         url.Values{},
			DefaultValue:   int64(0),
			ExpectedOutput: int64(0),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt64(r, key, defaultValue.(int64))
			},
		},
		// #20
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   int64(0),
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt64(r, key, defaultValue.(int64))
			},
		},
		// #21
		{
			Values:         url.Values{},
			DefaultValue:   int64(42),
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt64(r, key, defaultValue.(int64))
			},
		},
		// #22
		{
			Values:         url.Values{"key": []string{"42"}},
			DefaultValue:   int64(87),
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt64(r, key, defaultValue.(int64))
			},
		},
		// #23
		{
			Values:         url.Values{"key": []string{"invalid-int"}},
			DefaultValue:   int64(42),
			ExpectedOutput: int64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsInt64(r, key, defaultValue.(int64))
			},
		},
		// -- ParamsFloat32 --
		// #24
		{
			Values:         url.Values{},
			DefaultValue:   float32(0),
			ExpectedOutput: float32(0),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsFloat32(r, key, defaultValue.(float32))
			},
		},
		// #25
		{
			Values:         url.Values{"key": []string{"42.7"}},
			DefaultValue:   float32(0),
			ExpectedOutput: float32(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsFloat32(r, key, defaultValue.(float32))
			},
		},
		// #26
		{
			Values:         url.Values{},
			DefaultValue:   float32(42),
			ExpectedOutput: float32(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsFloat32(r, key, defaultValue.(float32))
			},
		},
		// #27
		{
			Values:         url.Values{"key": []string{"42.7"}},
			DefaultValue:   float32(87.3),
			ExpectedOutput: float32(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsFloat32(r, key, defaultValue.(float32))
			},
		},
		// #28
		{
			Values:         url.Values{"key": []string{"invalid-value"}},
			DefaultValue:   float32(42.7),
			ExpectedOutput: float32(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsFloat32(r, key, defaultValue.(float32))
			},
		},
		// -- ParamsFloat64 --
		// #29
		{
			Values:         url.Values{},
			DefaultValue:   float64(0),
			ExpectedOutput: float64(0),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsFloat64(r, key, defaultValue.(float64))
			},
		},
		// #30
		{
			Values:         url.Values{"key": []string{"42.7"}},
			DefaultValue:   float64(0),
			ExpectedOutput: float64(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsFloat64(r, key, defaultValue.(float64))
			},
		},
		// #31
		{
			Values:         url.Values{},
			DefaultValue:   float64(42),
			ExpectedOutput: float64(42),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsFloat64(r, key, defaultValue.(float64))
			},
		},
		// #32
		{
			Values:         url.Values{"key": []string{"42.7"}},
			DefaultValue:   float64(87.3),
			ExpectedOutput: float64(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsFloat64(r, key, defaultValue.(float64))
			},
		},
		// #33
		{
			Values:         url.Values{"key": []string{"invalid-value"}},
			DefaultValue:   float64(42.7),
			ExpectedOutput: float64(42.7),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsFloat64(r, key, defaultValue.(float64))
			},
		},
		// -- ParamsTime --
		// #34
		{
			Values:         url.Values{},
			DefaultValue:   zeroTime,
			ExpectedOutput: zeroTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsTime(r, key, time.RFC3339, defaultValue.(time.Time))
			},
		},
		// #35
		{
			Values:         url.Values{"key": []string{defaultTime.Truncate(time.Second).Format(time.RFC3339)}},
			DefaultValue:   zeroTime,
			ExpectedOutput: defaultTime.Truncate(time.Second),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsTime(r, key, time.RFC3339, defaultValue.(time.Time))
			},
		},
		// #36
		{
			Values:         url.Values{},
			DefaultValue:   defaultTime,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsTime(r, key, time.RFC3339, defaultValue.(time.Time))
			},
		},
		// #37
		{
			Values:         url.Values{"key": []string{defaultTime.Format(time.RFC3339Nano)}},
			DefaultValue:   zeroTime,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsTime(r, key, time.RFC3339Nano, defaultValue.(time.Time))
			},
		},
		// #38
		{
			Values:         url.Values{"key": []string{"invalid-value"}},
			DefaultValue:   defaultTime,
			ExpectedOutput: defaultTime,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsTime(r, key, time.RFC3339, defaultValue.(time.Time))
			},
		},
		// #39
		{
			Values:         url.Values{"key": []string{"2018-12-31"}},
			DefaultValue:   defaultTime,
			ExpectedOutput: time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC),
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsTime(r, key, "2006-01-02", defaultValue.(time.Time))
			},
		},
		// -- ParamsDuration --
		// #34
		{
			Values:         url.Values{},
			DefaultValue:   zeroDuration,
			ExpectedOutput: zeroDuration,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsDuration(r, key, defaultValue.(time.Duration))
			},
		},
		// #35
		{
			Values:         url.Values{"key": []string{"1m12s"}},
			DefaultValue:   zeroDuration,
			ExpectedOutput: 72 * time.Second,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsDuration(r, key, defaultValue.(time.Duration))
			},
		},
		// #36
		{
			Values:         url.Values{},
			DefaultValue:   12 * time.Second,
			ExpectedOutput: 12 * time.Second,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsDuration(r, key, defaultValue.(time.Duration))
			},
		},
		// #37
		{
			Values:         url.Values{"key": []string{"invalid-value"}},
			DefaultValue:   12 * time.Second,
			ExpectedOutput: 12 * time.Second,
			Invoke: func(r *http.Request, key string, defaultValue interface{}) interface{} {
				return ParamsDuration(r, key, defaultValue.(time.Duration))
			},
		},
	}
	for i, tt := range tests {
		req, _ := http.NewRequest("GET", "/", http.NoBody)
		vars := make(map[string]string)
		for k, v := range tt.Values {
			if len(v) > 0 {
				vars[k] = v[0]
			} else {
				vars[k] = ""
			}
		}
		req = mux.SetURLVars(req, vars)
		got := tt.Invoke(req, "key", tt.DefaultValue)
		if !cmp.Equal(got, tt.ExpectedOutput) {
			t.Fatalf("#%d. expected Output = %v; got %v", i, tt.ExpectedOutput, got)
		}
	}
}
