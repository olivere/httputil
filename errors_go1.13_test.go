// +build go1.13

package httputil

import (
	"errors"
	"io"
	"testing"
)

func TestUnwrapping(t *testing.T) {
	innerErr := io.EOF
	outerErr := TeapotError{Message: "I used to brew coffee", Err: innerErr}
	if !errors.Is(outerErr, innerErr) {
		t.Error("expected error to be an io.EOF")
	}
	if err := errors.Unwrap(outerErr); err != innerErr {
		t.Errorf("expected the unwrapped error to be %T; got %T", innerErr, err)
	}
}

func TestMissingParameterErrorIsBadRequestError(t *testing.T) {
	err := MissingParameterError("name")
	if !errors.Is(err, BadRequestError{Message: `Missing parameter "name"`}) {
		t.Error("expected MissingParameterError to be a BadRequestError")
	}
}

func TestInvalidParameterErrorIsBadRequestError(t *testing.T) {
	err := InvalidParameterError("name")
	if !errors.Is(err, BadRequestError{Message: `Invalid parameter "name"`}) {
		t.Error("expected InvalidParameterError to be a BadRequestError")
	}
}
