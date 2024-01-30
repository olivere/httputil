// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BadRequestError returns HTTP status 400 and an error message as HTML.
func BadRequestError(w http.ResponseWriter, errorMessage string, args ...interface{}) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "<h1>Bad request</h1>")
}

// ForbiddenError returns HTTP status 403 and an error message as HTML.
func ForbiddenError(w http.ResponseWriter, errorMessage string, args ...interface{}) {
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprintf(w, "<h1>Forbidden</h1>")
}

// InternalServerError returns HTTP status 500 and an error message as HTML.
func InternalServerError(w http.ResponseWriter, r *http.Request, err interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "<h1>Server error</h1>")
}

// WriteError writes an error message for display in a HTML page.
func WriteError(w http.ResponseWriter, err interface{}) {
	code := 500
	if i, ok := err.(httpCoder); ok {
		code = i.HTTPCode()
	}
	msg := fmt.Sprint(err)
	w.WriteHeader(code)
	fmt.Fprintf(w, "<h1>%s</h1>", msg)
}

// WriteJSONError writes error information, serialized in a JSON structure.
// Example:
//
//	{
//	  "error":{
//	    "code":    500,
//	    "message": "Something went wrong",
//	    "details": ["A was bad", "B is missing"]
//	  }
//	}
//
// If err implements the httpCoder interface, it can specify the HTTP code
// to return. If err implements the httpErrorDetails interface, its
// ErrorDetails func is used to collect the error details; otherwise,
// the "details" field is missing in the error returned.
func WriteJSONError(w http.ResponseWriter, err interface{}) {
	code := 500
	if i, ok := err.(httpCoder); ok {
		code = i.HTTPCode()
	}
	var details []string
	if i, ok := err.(httpErrorDetails); ok {
		details = i.ErrorDetails()
	}
	msg := fmt.Sprint(err)
	innerErr := map[string]interface{}{
		"code":    code,
		"message": msg,
	}
	if len(details) > 0 {
		innerErr["details"] = details
	}

	WriteJSONCode(w, code, map[string]interface{}{
		"error": innerErr,
	})
}

// httpCoder provides an interface to return the HTTP status code
// in an error. See InvalidMethodError for an example.
type httpCoder interface {
	HTTPCode() int
}

// httpErrorDetails provides an interface to return a list of error
// details. See UnprocessableEntityError for an example.
type httpErrorDetails interface {
	ErrorDetails() []string
}

// InvalidMethodError indicates that an invalid HTTP method is being used.
type InvalidMethodError struct{}

// Error returns the error in text form.
func (InvalidMethodError) Error() string { return "Invalid HTTP method" }

// HTTPCode returns the HTTP status code of the error.
func (InvalidMethodError) HTTPCode() int { return http.StatusMethodNotAllowed }

// UnauthorizedError indicates that credentials are either missing or invalid.
type UnauthorizedError struct{}

// Error returns the error in text form.
func (UnauthorizedError) Error() string { return "Missing or invalid credentials" }

// HTTPCode returns the HTTP status code of the error.
func (UnauthorizedError) HTTPCode() int { return http.StatusUnauthorized }

// NotFoundError indicates that a record or resource does not exist.
type NotFoundError struct{}

// Error returns the error in text form.
func (NotFoundError) Error() string { return "Record not found" }

// HTTPCode returns the HTTP status code of the error.
func (NotFoundError) HTTPCode() int { return http.StatusNotFound }

// InvalidJSONError indicates that the JSON data are invalid.
type InvalidJSONError struct {
	error
}

// HTTPCode returns the HTTP status code of the error.
func (InvalidJSONError) HTTPCode() int { return http.StatusBadRequest }

// MissingParameterError indicates that a required parameter is missing or blank.
type MissingParameterError string

// Error returns the error in text form.
func (p MissingParameterError) Error() string { return fmt.Sprintf("Missing parameter %q", string(p)) }

// HTTPCode returns the HTTP status code of the error.
func (MissingParameterError) HTTPCode() int { return http.StatusBadRequest }

// InvalidParameterError indicates that a parameter is invalid.
type InvalidParameterError string

// Error returns the error in text form.
func (p InvalidParameterError) Error() string { return fmt.Sprintf("Invalid parameter %q", string(p)) }

// HTTPCode returns the HTTP status code of the error.
func (InvalidParameterError) HTTPCode() int { return http.StatusBadRequest }

// InvalidXSRFToken indicates that the user has not provided a valid XSRF token.
type InvalidXSRFToken struct{}

// Error returns the error in text form.
func (InvalidXSRFToken) Error() string { return "Invalid or missing XSRF token" }

// HTTPCode returns the HTTP status code of the error.
func (InvalidXSRFToken) HTTPCode() int { return http.StatusBadRequest }

// UnprocessableEntityError indicates that there was a semantic error in
// parsing a request, e.g. a record with validation errors.
type UnprocessableEntityError struct {
	Errors []string
}

// Error returns the error in text form.
func (UnprocessableEntityError) Error() string { return "Record has semantic errors" }

// HTTPCode returns the HTTP status code of the error.
func (UnprocessableEntityError) HTTPCode() int { return 422 }

// ErrorDetails returns additional information about the error.
func (p UnprocessableEntityError) ErrorDetails() []string { return p.Errors }

// TimeoutError indicates that the request has timed out.
type TimeoutError struct{}

// Error returns the error in text form.
func (TimeoutError) Error() string { return "Request has timed out" }

// HTTPCode returns the HTTP status code of the error.
func (TimeoutError) HTTPCode() int { return http.StatusGatewayTimeout }

// ServerError indicates any kind of internal server problem.
type ServerError string

// Error returns the error in text form.
func (e ServerError) Error() string { return string(e) }

// HTTPCode returns the HTTP status code of the error.
func (ServerError) HTTPCode() int { return http.StatusInternalServerError }

// NotImplementedError indicates that an endpoint has yet to be implemented.
type NotImplementedError struct{}

// Error returns the error in text form.
func (NotImplementedError) Error() string { return "Not implemented" }

// HTTPCode returns the HTTP status code of the error.
func (NotImplementedError) HTTPCode() int { return http.StatusNotImplemented }

// GrpcError is a placeholder for a gRPC error, and will turn it into a HTTP error.
type GrpcError struct {
	Err error
}

// Error returns the error message.
func (e GrpcError) Error() string {
	if s, ok := status.FromError(e.Err); ok {
		return s.Message()
	}
	return "Internal server error"
}

// HTTPCode returns the HTTP status code of the gRPC error.
func (e GrpcError) HTTPCode() int {
	switch status.Code(e.Err) {
	case codes.OK:
		return http.StatusOK
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
