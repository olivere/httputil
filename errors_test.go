// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	allErrors = []struct {
		Error            error
		ExpectedHTTPCode int
		ExpectedMessage  string
		ExpectedDetails  []string
	}{
		// 400
		{
			Error:            BadRequestError{},
			ExpectedHTTPCode: http.StatusBadRequest,
			ExpectedMessage:  http.StatusText(http.StatusBadRequest),
		},
		{
			Error:            BadRequestError{Message: "Validation failed"},
			ExpectedHTTPCode: http.StatusBadRequest,
			ExpectedMessage:  "Validation failed",
		},
		{
			Error:            BadRequestError{Message: "Validation failed", Details: []string{"Name required", "Price required"}},
			ExpectedHTTPCode: http.StatusBadRequest,
			ExpectedMessage:  "Validation failed",
			ExpectedDetails:  []string{"Name required", "Price required"},
		},
		// 401
		{
			Error:            UnauthorizedError{},
			ExpectedHTTPCode: http.StatusUnauthorized,
			ExpectedMessage:  http.StatusText(http.StatusUnauthorized),
		},
		// 402
		{
			Error:            PaymentRequiredError{},
			ExpectedHTTPCode: http.StatusPaymentRequired,
			ExpectedMessage:  http.StatusText(http.StatusPaymentRequired),
		},
		// 403
		{
			Error:            ForbiddenError{},
			ExpectedHTTPCode: http.StatusForbidden,
			ExpectedMessage:  http.StatusText(http.StatusForbidden),
		},
		// 404
		{
			Error:            NotFoundError{},
			ExpectedHTTPCode: http.StatusNotFound,
			ExpectedMessage:  http.StatusText(http.StatusNotFound),
		},
		// 405
		{
			Error:            MethodNotAllowedError{},
			ExpectedHTTPCode: http.StatusMethodNotAllowed,
			ExpectedMessage:  http.StatusText(http.StatusMethodNotAllowed),
		},
		// 406
		{
			Error:            NotAcceptableError{},
			ExpectedHTTPCode: http.StatusNotAcceptable,
			ExpectedMessage:  http.StatusText(http.StatusNotAcceptable),
		},
		// 407
		{
			Error:            ProxyAuthRequiredError{},
			ExpectedHTTPCode: http.StatusProxyAuthRequired,
			ExpectedMessage:  http.StatusText(http.StatusProxyAuthRequired),
		},
		// 408
		{
			Error:            RequestTimeoutError{},
			ExpectedHTTPCode: http.StatusRequestTimeout,
			ExpectedMessage:  http.StatusText(http.StatusRequestTimeout),
		},
		// 409
		{
			Error:            ConflictError{},
			ExpectedHTTPCode: http.StatusConflict,
			ExpectedMessage:  http.StatusText(http.StatusConflict),
		},
		// 410
		{
			Error:            GoneError{},
			ExpectedHTTPCode: http.StatusGone,
			ExpectedMessage:  http.StatusText(http.StatusGone),
		},
		// 411
		{
			Error:            LengthRequiredError{},
			ExpectedHTTPCode: http.StatusLengthRequired,
			ExpectedMessage:  http.StatusText(http.StatusLengthRequired),
		},
		// 412
		{
			Error:            PreconditionFailedError{},
			ExpectedHTTPCode: http.StatusPreconditionFailed,
			ExpectedMessage:  http.StatusText(http.StatusPreconditionFailed),
		},
		// 413
		{
			Error:            RequestEntityTooLargeError{},
			ExpectedHTTPCode: http.StatusRequestEntityTooLarge,
			ExpectedMessage:  http.StatusText(http.StatusRequestEntityTooLarge),
		},
		// 414
		{
			Error:            RequestURITooLongError{},
			ExpectedHTTPCode: http.StatusRequestURITooLong,
			ExpectedMessage:  http.StatusText(http.StatusRequestURITooLong),
		},
		// 415
		{
			Error:            UnsupportedMediaTypeError{},
			ExpectedHTTPCode: http.StatusUnsupportedMediaType,
			ExpectedMessage:  http.StatusText(http.StatusUnsupportedMediaType),
		},
		// 416
		{
			Error:            RequestedRangeNotSatisfiableError{},
			ExpectedHTTPCode: http.StatusRequestedRangeNotSatisfiable,
			ExpectedMessage:  http.StatusText(http.StatusRequestedRangeNotSatisfiable),
		},
		// 417
		{
			Error:            ExpectationFailedError{},
			ExpectedHTTPCode: http.StatusExpectationFailed,
			ExpectedMessage:  http.StatusText(http.StatusExpectationFailed),
		},
		// 418
		{
			Error:            TeapotError{},
			ExpectedHTTPCode: http.StatusTeapot,
			ExpectedMessage:  http.StatusText(http.StatusTeapot),
		},
		// 421
		{
			Error:            MisdirectedRequestError{},
			ExpectedHTTPCode: http.StatusMisdirectedRequest,
			ExpectedMessage:  http.StatusText(http.StatusMisdirectedRequest),
		},
		// 422
		{
			Error:            UnprocessableEntityError{},
			ExpectedHTTPCode: http.StatusUnprocessableEntity,
			ExpectedMessage:  http.StatusText(http.StatusUnprocessableEntity),
		},
		// 423
		{
			Error:            LockedError{},
			ExpectedHTTPCode: http.StatusLocked,
			ExpectedMessage:  http.StatusText(http.StatusLocked),
		},
		// 424
		{
			Error:            FailedDependencyError{},
			ExpectedHTTPCode: http.StatusFailedDependency,
			ExpectedMessage:  http.StatusText(http.StatusFailedDependency),
		},
		// 425
		{
			Error:            TooEarlyError{},
			ExpectedHTTPCode: http.StatusTooEarly,
			ExpectedMessage:  http.StatusText(http.StatusTooEarly),
		},
		// 426
		{
			Error:            UpgradeRequiredError{},
			ExpectedHTTPCode: http.StatusUpgradeRequired,
			ExpectedMessage:  http.StatusText(http.StatusUpgradeRequired),
		},
		// 428
		{
			Error:            PreconditionRequiredError{},
			ExpectedHTTPCode: http.StatusPreconditionRequired,
			ExpectedMessage:  http.StatusText(http.StatusPreconditionRequired),
		},
		// 429
		{
			Error:            TooManyRequestsError{},
			ExpectedHTTPCode: http.StatusTooManyRequests,
			ExpectedMessage:  http.StatusText(http.StatusTooManyRequests),
		},
		// 431
		{
			Error:            RequestHeaderFieldsTooLargeError{},
			ExpectedHTTPCode: http.StatusRequestHeaderFieldsTooLarge,
			ExpectedMessage:  http.StatusText(http.StatusRequestHeaderFieldsTooLarge),
		},
		// 451
		{
			Error:            UnavailableForLegalReasonsError{},
			ExpectedHTTPCode: http.StatusUnavailableForLegalReasons,
			ExpectedMessage:  http.StatusText(http.StatusUnavailableForLegalReasons),
		},
		// 500
		{
			Error:            InternalServerError{},
			ExpectedHTTPCode: http.StatusInternalServerError,
			ExpectedMessage:  http.StatusText(http.StatusInternalServerError),
		},
		// 501
		{
			Error:            NotImplementedError{},
			ExpectedHTTPCode: http.StatusNotImplemented,
			ExpectedMessage:  http.StatusText(http.StatusNotImplemented),
		},
		// 502
		{
			Error:            BadGatewayError{},
			ExpectedHTTPCode: http.StatusBadGateway,
			ExpectedMessage:  http.StatusText(http.StatusBadGateway),
		},
		// 503
		{
			Error:            ServiceUnavailableError{},
			ExpectedHTTPCode: http.StatusServiceUnavailable,
			ExpectedMessage:  http.StatusText(http.StatusServiceUnavailable),
		},
		// 504
		{
			Error:            GatewayTimeoutError{},
			ExpectedHTTPCode: http.StatusGatewayTimeout,
			ExpectedMessage:  http.StatusText(http.StatusGatewayTimeout),
		},
		// 505
		{
			Error:            HTTPVersionNotSupportedError{},
			ExpectedHTTPCode: http.StatusHTTPVersionNotSupported,
			ExpectedMessage:  http.StatusText(http.StatusHTTPVersionNotSupported),
		},
		// 506
		{
			Error:            VariantAlsoNegotiatesError{},
			ExpectedHTTPCode: http.StatusVariantAlsoNegotiates,
			ExpectedMessage:  http.StatusText(http.StatusVariantAlsoNegotiates),
		},
		// 507
		{
			Error:            InsufficientStorageError{},
			ExpectedHTTPCode: http.StatusInsufficientStorage,
			ExpectedMessage:  http.StatusText(http.StatusInsufficientStorage),
		},
		// 508
		{
			Error:            LoopDetectedError{},
			ExpectedHTTPCode: http.StatusLoopDetected,
			ExpectedMessage:  http.StatusText(http.StatusLoopDetected),
		},
		// 510
		{
			Error:            NotExtendedError{},
			ExpectedHTTPCode: http.StatusNotExtended,
			ExpectedMessage:  http.StatusText(http.StatusNotExtended),
		},
		// 511
		{
			Error:            NetworkAuthenticationRequiredError{},
			ExpectedHTTPCode: http.StatusNetworkAuthenticationRequired,
			ExpectedMessage:  http.StatusText(http.StatusNetworkAuthenticationRequired),
		},
	}
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

func TestWriteJSONErrorWithHTTPCoder(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		WriteJSONError(w, BadRequestError{Message: `Invalid parameter "pin"`})
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

func TestWriteJSONErrorWithHTTPErrorDetails(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		err := UnprocessableEntityError{
			Details: []string{"A has failed", "B is invalid"},
		}
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
	if want, have := http.StatusUnprocessableEntity, fail.Error.Code; want != have {
		t.Errorf("expected error code = %d; got: %d", want, have)
	}
	if want, have := http.StatusText(http.StatusUnprocessableEntity), fail.Error.Message; want != have {
		t.Errorf("expected error message = %q; got: %q", want, have)
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

func TestHTTPErrors(t *testing.T) {
	for _, tt := range allErrors {
		h := func(w http.ResponseWriter, _ *http.Request) {
			WriteJSONError(w, tt.Error)
		}
		w := httptest.NewRecorder()
		h(w, nil)

		type failure struct {
			Error struct {
				Code    int      `json:"code"`
				Message string   `json:"message"`
				Details []string `json:"details"`
			} `json:"error"`
		}
		var fail failure
		err := json.NewDecoder(w.Body).Decode(&fail)
		if err != nil {
			t.Fatal(err)
		}
		if want, have := tt.ExpectedHTTPCode, fail.Error.Code; want != have {
			t.Errorf("with %#v: expected HTTP status code = %d; got: %d", tt.Error, want, have)
		}
		if want, have := tt.ExpectedMessage, fail.Error.Message; want != have {
			t.Errorf("with %#v: expected error message = %q; got: %q", tt.Error, want, have)
		}
		if want, have := tt.ExpectedDetails, fail.Error.Details; !cmp.Equal(want, have) {
			t.Errorf("with %#v: expected different error details:\n%s", tt.Error, cmp.Diff(want, have))
		}
	}
}
