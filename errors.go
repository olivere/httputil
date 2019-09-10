// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil

import (
	"fmt"
	"net/http"
)

// WriteError writes an error message for display in a HTML page.
func WriteError(w http.ResponseWriter, err interface{}) {
	code := http.StatusInternalServerError
	if i, ok := err.(httpCoder); ok {
		code = i.HTTPCode()
	}
	w.WriteHeader(code)
	fmt.Fprintf(w, "<h1>%s</h1>", fmt.Sprint(err))
}

// WriteJSONError writes error information, serialized in a JSON structure.
// Example:
//
//   {
//     "error":{
//       "code":    500,
//       "message": "Something went wrong",
//       "details": ["A was bad", "B is missing"]
//     }
//   }
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
// in an error.
type httpCoder interface {
	HTTPCode() int
}

// httpErrorDetails provides an interface to return a list of error
// details.
type httpErrorDetails interface {
	ErrorDetails() []string
}

// --

// Error represents a generic HTTP error with a HTTP status code,
// an error message, and (optional) error details.
type Error struct {
	Code    int
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e Error) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	if s := http.StatusText(e.Code); s != "" {
		return s
	}
	return "Unknown error"
}

// ErrorDetails returns additional details about the error.
func (e Error) ErrorDetails() []string {
	return e.Details
}

// Unwrap returns the inner error which might be nil.
func (e Error) Unwrap() error {
	return e.Err
}

// HTTPCode returns the HTTP status code of the error.
func (e Error) HTTPCode() int { return e.Code }

// --

// BadRequestError (400) indicates that the server cannot or will not
// process the request due to something that is perceived to be a client
// error (e.g., malformed request syntax, invalid request message
// framing, or deceptive request routing).
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.1 for details.
type BadRequestError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e BadRequestError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e BadRequestError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e BadRequestError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (BadRequestError) HTTPCode() int { return http.StatusBadRequest }

// --

// UnauthorizedError (401) indicates that the request has not
// been applied because it lacks valid authentication credentials for
// the target resource.  The server generating a 401 response MUST send
// a WWW-Authenticate header
// field (https://tools.ietf.org/html/rfc7235#section-4.1) containing
// at least one challenge applicable to the target resource.
//
// If the request included authentication credentials, then the 401
// response indicates that authorization has been refused for those
// credentials.  The user agent MAY repeat the request with a new or
// replaced Authorization
// header field (https://tools.ietf.org/html/rfc7235#section-4.2).
// If the 401 response contains the same challenge as the prior response,
// and the user agent has already attempted authentication at least once,
// then the user agent SHOULD present the enclosed representation to the
// user, since it usually contains relevant diagnostic information.
//
// See https://tools.ietf.org/html/rfc7235#section-3.1 for details.
type UnauthorizedError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e UnauthorizedError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e UnauthorizedError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e UnauthorizedError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (UnauthorizedError) HTTPCode() int { return http.StatusUnauthorized }

// --

// PaymentRequiredError (402) indicates is reserved for future use.
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.2 for details.
type PaymentRequiredError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e PaymentRequiredError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e PaymentRequiredError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e PaymentRequiredError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (PaymentRequiredError) HTTPCode() int { return http.StatusPaymentRequired }

// --

// ForbiddenError (403) indicates that the server understood the request but
// refuses to authorize it.  A server that wishes to make public why
// the request has been forbidden can describe that reason in the response
// payload (if any).
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.3 for details.
type ForbiddenError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e ForbiddenError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e ForbiddenError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e ForbiddenError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (ForbiddenError) HTTPCode() int { return http.StatusForbidden }

// --

// NotFoundError (404) indicates that the origin server did not find a
// current representation for the target resource or is not willing to
// disclose that one exists.
// A 404 status code does not indicate whether this
// lack of representation is temporary or permanent; the 410 (Gone)
// status code is preferred over 404 if the origin server knows,
// presumably through some configurable means, that the condition is
// likely to be permanent.
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.4 for details.
type NotFoundError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e NotFoundError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e NotFoundError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e NotFoundError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (NotFoundError) HTTPCode() int { return http.StatusNotFound }

// --

// InvalidMethodError (405) is an alias for MethodNotAllowedError.
//
// Deprecated: Use MethodNotAllowedError instead.
type InvalidMethodError = MethodNotAllowedError

// MethodNotAllowedError (405) indicates that the method received in the
// request-line is known by the origin server but not supported by
// the target resource.  The origin server MUST generate an
// Allow header field in a 405 response containing a list of the target
// resource's currently supported methods.
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.5 for details.
type MethodNotAllowedError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e MethodNotAllowedError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e MethodNotAllowedError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e MethodNotAllowedError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (MethodNotAllowedError) HTTPCode() int { return http.StatusMethodNotAllowed }

// --

// NotAcceptableError (406) indicates that the target resource does not have
// a current representation that would be acceptable to the user agent,
// according to the proactive negotiation header fields received in the
// request, and the server is unwilling to supply a default representation.
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.6 for details.
type NotAcceptableError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e NotAcceptableError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e NotAcceptableError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e NotAcceptableError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (NotAcceptableError) HTTPCode() int { return http.StatusNotAcceptable }

// --

// ProxyAuthRequiredError (407) is similar to 401
// (Unauthorized), but it indicates that the client needs to
// authenticate itself in order to use a proxy.  The proxy MUST send a
// Proxy-Authenticate header field (Section 4.3 of RFC7235) containing a
// challenge applicable to that proxy for the target resource.  The client MAY
// repeat the request with a new or replaced Proxy-Authorization header
// field (Section 4.4 of RFC7235).
//
// See https://tools.ietf.org/html/rfc7235#section-3.2 for details.
type ProxyAuthRequiredError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e ProxyAuthRequiredError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e ProxyAuthRequiredError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e ProxyAuthRequiredError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (ProxyAuthRequiredError) HTTPCode() int { return http.StatusProxyAuthRequired }

// --

// RequestTimeoutError (408) indicates that the server did not receive
// a complete request message within the time that it was prepared to wait.
// A server SHOULD send the "close" connection option
// (Section 6.1 of RFC7230) in the response, since 408 implies that
// the server has decided to close the connection rather than continue
// waiting.  If the client has an outstanding request in transit, the
// client MAY repeat that request on a new connection.
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.7 for details.
type RequestTimeoutError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e RequestTimeoutError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e RequestTimeoutError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e RequestTimeoutError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (RequestTimeoutError) HTTPCode() int { return http.StatusRequestTimeout }

// --

// ConflictError (409) indicates that the request could not
// be completed due to a conflict with the current state of the target
// resource.  This code is used in situations where the user might be
// able to resolve the conflict and resubmit the request.  The server
// SHOULD generate a payload that includes enough information for a user
// to recognize the source of the conflict.
//
// Conflicts are most likely to occur in response to a PUT request.  For
// example, if versioning were being used and the representation being
// PUT included changes to a resource that conflict with those made by
// an earlier (third-party) request, the origin server might use a 409
// response to indicate that it can't complete the request.  In this
// case, the response representation would likely contain information
// useful for merging the differences based on the revision history.
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.8 for details.
type ConflictError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e ConflictError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e ConflictError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e ConflictError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (ConflictError) HTTPCode() int { return http.StatusConflict }

// --

// GoneError (410) indicates that access to the target
// resource is no longer available at the origin server and that this
// condition is likely to be permanent.  If the origin server does not
// know, or has no facility to determine, whether or not the condition
// is permanent, the status code 404 (Not Found) ought to be used
// instead.
//
// The 410 response is primarily intended to assist the task of web
// maintenance by notifying the recipient that the resource is
// intentionally unavailable and that the server owners desire that
// remote links to that resource be removed.  Such an event is common
// for limited-time, promotional services and for resources belonging to
// individuals no longer associated with the origin server's site.  It
// is not necessary to mark all permanently unavailable resources as
// "gone" or to keep the mark for any length of time -- that is left to
// the discretion of the server owner.
//
// A 410 response is cacheable by default; i.e., unless otherwise
// indicated by the method definition or explicit cache controls (see
// Section 4.2.2 of RFC7234).
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.9 for details.
type GoneError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e GoneError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e GoneError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e GoneError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (GoneError) HTTPCode() int { return http.StatusGone }

// --

// LengthRequiredError (411) indicates that the server
// refuses to accept the request without a defined Content-Length
// (Section 3.3.2 of RFC7230).  The client MAY repeat the request if
// it adds a valid Content-Length header field containing the length of
// the message body in the request message.
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.10 for details.
type LengthRequiredError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e LengthRequiredError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e LengthRequiredError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e LengthRequiredError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (LengthRequiredError) HTTPCode() int { return http.StatusLengthRequired }

// --

// PreconditionFailedError (412) indicates that one or more
// conditions given in the request header fields evaluated to false when
// tested on the server.  This response code allows the client to place
// preconditions on the current resource state (its current
// representations and metadata) and, thus, prevent the request method
// from being applied if the target resource is in an unexpected state.
//
// See https://tools.ietf.org/html/rfc7232#section-4.2 for details.
type PreconditionFailedError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e PreconditionFailedError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e PreconditionFailedError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e PreconditionFailedError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (PreconditionFailedError) HTTPCode() int { return http.StatusPreconditionFailed }

// --

// RequestEntityTooLargeError (413) indicates that the server is
// refusing to process a request because the request payload is larger
// than the server is willing or able to process.  The server MAY close
// the connection to prevent the client from continuing the request.
//
// If the condition is temporary, the server SHOULD generate a
// Retry-After header field to indicate that it is temporary and after
// what time the client MAY try again.
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.11 for details.
type RequestEntityTooLargeError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e RequestEntityTooLargeError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e RequestEntityTooLargeError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e RequestEntityTooLargeError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (RequestEntityTooLargeError) HTTPCode() int { return http.StatusRequestEntityTooLarge }

// --

// RequestURITooLongError (414) indicates that the server is
// refusing to service the request because the request-target (Section
// 5.3 of RFC7230) is longer than the server is willing to interpret.
// This rare condition is only likely to occur when a client has
// improperly converted a POST request to a GET request with long query
// information, when the client has descended into a "black hole" of
// redirection (e.g., a redirected URI prefix that points to a suffix of
// itself) or when the server is under attack by a client attempting to
// exploit potential security holes.
//
// A 414 response is cacheable by default; i.e., unless otherwise
// indicated by the method definition or explicit cache controls (see
// Section 4.2.2 of RFC7234).
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.12 for details.
type RequestURITooLongError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e RequestURITooLongError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e RequestURITooLongError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e RequestURITooLongError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (RequestURITooLongError) HTTPCode() int { return http.StatusRequestURITooLong }

// --

// UnsupportedMediaTypeError (415) indicates that the
// origin server is refusing to service the request because the payload
// is in a format not supported by this method on the target resource.
// The format problem might be due to the request's indicated
// Content-Type or Content-Encoding, or as a result of inspecting the
// data directly.
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.13 for details.
type UnsupportedMediaTypeError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e UnsupportedMediaTypeError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e UnsupportedMediaTypeError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e UnsupportedMediaTypeError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (UnsupportedMediaTypeError) HTTPCode() int { return http.StatusUnsupportedMediaType }

// --

// RequestedRangeNotSatisfiableError (416) indicates that none of
// the ranges in the request's Range header field (Section 3.1 of RFC7233)
// overlap the current extent of the selected resource or that the set of
// ranges requested has been rejected due to invalid ranges or an excessive
// request of small or overlapping ranges.
//
// For byte ranges, failing to overlap the current extent means that the
// first-byte-pos of all of the byte-range-spec values were greater than
// the current length of the selected representation.  When this status
// code is generated in response to a byte-range request, the sender
// SHOULD generate a Content-Range header field specifying the current
// length of the selected representation (Section 4.2 of RFC7233).
//
// See https://tools.ietf.org/html/rfc7233#section-4.4 for details.
type RequestedRangeNotSatisfiableError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e RequestedRangeNotSatisfiableError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e RequestedRangeNotSatisfiableError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e RequestedRangeNotSatisfiableError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (RequestedRangeNotSatisfiableError) HTTPCode() int {
	return http.StatusRequestedRangeNotSatisfiable
}

// --

// ExpectationFailedError (417) indicates that the expectation given in the
// request's Expect header field (Section 5.1.1 of RFC7231) could not
// be met by at least one of the inbound servers.
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.14 for details.
type ExpectationFailedError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e ExpectationFailedError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e ExpectationFailedError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e ExpectationFailedError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (ExpectationFailedError) HTTPCode() int { return http.StatusExpectationFailed }

// --

// TeapotError (418) indicates how TEA-capable pots that are not provisioned to
// brew coffee may return either a status code of 503, indicating temporary
// unavailability of coffee, or a code of 418 as defined in the base
// HTCPCP specification to denote a more permanent indication that the
// pot is a teapot.
//
// See https://tools.ietf.org/html/rfc7168#section-2.3.3 for details.
type TeapotError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e TeapotError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e TeapotError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e TeapotError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (TeapotError) HTTPCode() int { return http.StatusTeapot }

// --

// MisdirectedRequestError (421) indicates that the request
// was directed at a server that is not able to produce a response.
// This can be sent by a server that is not configured to produce
// responses for the combination of scheme and authority that are
// included in the request URI.
//
// Clients receiving a 421 (Misdirected Request) response from a server
// MAY retry the request -- whether the request method is idempotent or
// not -- over a different connection.  This is possible if a connection
// is reused (Section 9.1.1 of RFC7540) or if an alternative service is
// selected.
//
// This status code MUST NOT be generated by proxies.
//
// A 421 response is cacheable by default, i.e., unless otherwise
// indicated by the method definition or explicit cache controls (see
// Section 4.2.2 of of RFC7234).
//
// See https://tools.ietf.org/html/rfc7540#section-9.1.2 for details.
type MisdirectedRequestError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e MisdirectedRequestError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e MisdirectedRequestError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e MisdirectedRequestError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (MisdirectedRequestError) HTTPCode() int { return http.StatusMisdirectedRequest }

// --

// UnprocessableEntityError (422) indicates that the server
// understands the content type of the request entity (hence a
// 415(Unsupported Media Type) status code is inappropriate), and the
// syntax of the request entity is correct (thus a 400 (Bad Request)
// status code is inappropriate) but was unable to process the contained
// instructions.  For example, this error condition may occur if an XML
// request body contains well-formed (i.e., syntactically correct), but
// semantically erroneous, XML instructions.
//
// See https://tools.ietf.org/html/rfc4918#section-11.2 for details.
type UnprocessableEntityError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e UnprocessableEntityError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e UnprocessableEntityError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e UnprocessableEntityError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (UnprocessableEntityError) HTTPCode() int { return http.StatusUnprocessableEntity }

// --

// LockedError (423) indicates that the source or destination resource
// of a method is locked.  This response SHOULD contain an appropriate
// precondition or postcondition code, such as 'lock-token-submitted' or
// 'no-conflicting-lock'.
//
// See https://tools.ietf.org/html/rfc4918#section-11.3 for details.
type LockedError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e LockedError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e LockedError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e LockedError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (LockedError) HTTPCode() int { return http.StatusLocked }

// --

// FailedDependencyError (424) indicates that the method could
// not be performed on the resource because the requested action
// depended on another action and that action failed.  For example, if a
// command in a PROPPATCH method fails, then, at minimum, the rest of
// the commands will also fail with 424 (Failed Dependency).
//
// See https://tools.ietf.org/html/rfc4918#section-11.4 for details.
type FailedDependencyError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e FailedDependencyError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e FailedDependencyError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e FailedDependencyError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (FailedDependencyError) HTTPCode() int { return http.StatusFailedDependency }

// --

// TooEarlyError (425) indicates that the server is unwilling
// to risk processing a request that might be replayed.
//
// User agents that send a request in early data are expected to retry
// the request when receiving a 425 (Too Early) response status code.  A
// user agent SHOULD retry automatically, but any retries MUST NOT be
// sent in early data.
//
// In all cases, an intermediary can forward a 425 (Too Early) status
// code.  Intermediaries MUST forward a 425 (Too Early) status code if
// the request that it received and forwarded contained an Early-Data
// header field.  Otherwise, an intermediary that receives a request in
// early data MAY automatically retry that request in response to a 425
// (Too Early) status code, but it MUST wait for the TLS handshake to
// complete on the connection where it received the request.
//
// The server cannot assume that a client is able to retry a request
// unless the request is received in early data or the Early-Data header
// field is set to "1".  A server SHOULD NOT emit the 425 status code
// unless one of these conditions is met.
//
// The 425 (Too Early) status code is not cacheable by default. Its
// payload is not the representation of any identified resource.
//
// See https://tools.ietf.org/html/rfc8470#section-5.2 for details.
type TooEarlyError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e TooEarlyError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e TooEarlyError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e TooEarlyError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (TooEarlyError) HTTPCode() int { return http.StatusTooEarly }

// --

// UpgradeRequiredError (426) indicates that the server refuses to perform
// the request using the current protocol but might be willing to do so
// after the client upgrades to a different protocol. The server MUST send
// an Upgrade header field in a 426 response to indicate the required
// protocol(s) (Section 6.7 of RFC7230).
//
// See https://tools.ietf.org/html/rfc7231#section-6.5.15 for details.
type UpgradeRequiredError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e UpgradeRequiredError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e UpgradeRequiredError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e UpgradeRequiredError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (UpgradeRequiredError) HTTPCode() int { return http.StatusUpgradeRequired }

// --

// PreconditionRequiredError (428) indicates that the origin server
// requires the request to be conditional.
//
// Its typical use is to avoid the "lost update" problem, where a client
// GETs a resource's state, modifies it, and PUTs it back to the server,
// when meanwhile a third party has modified the state on the server,
// leading to a conflict.  By requiring requests to be conditional, the
// server can assure that clients are working with the correct copies.
//
// Responses using this status code SHOULD explain how to resubmit the
// request successfully. For example:
//
//    HTTP/1.1 428 Precondition Required
//    Content-Type: text/html
//
//    <html>
//       <head>
//          <title>Precondition Required</title>
//       </head>
//       <body>
//          <h1>Precondition Required</h1>
//          <p>This request is required to be conditional;
//          try using "If-Match".</p>
//       </body>
//    </html>
//
// Responses with the 428 status code MUST NOT be stored by a cache.
//
// See https://tools.ietf.org/html/rfc6585#section-3 for details.
type PreconditionRequiredError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e PreconditionRequiredError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e PreconditionRequiredError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e PreconditionRequiredError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (PreconditionRequiredError) HTTPCode() int { return http.StatusPreconditionRequired }

// --

// TooManyRequestsError (429) indicates that the user has sent too many
// requests in a given amount of time ("rate limiting").
//
// The response representations SHOULD include details explaining the
// condition, and MAY include a Retry-After header indicating how long
// to wait before making a new request.
//
// For example:
//
//    HTTP/1.1 429 Too Many Requests
//    Content-Type: text/html
//    Retry-After: 3600
//
//    <html>
//       <head>
//          <title>Too Many Requests</title>
//       </head>
//       <body>
//          <h1>Too Many Requests</h1>
//          <p>I only allow 50 requests per hour to this Web site per
//             logged in user.  Try again soon.</p>
//       </body>
//    </html>
//
// Note that this specification does not define how the origin server
// identifies the user, nor how it counts requests.  For example, an
// origin server that is limiting request rates can do so based upon
// counts of requests on a per-resource basis, across the entire server,
// or even among a set of servers.  Likewise, it might identify the user
// by its authentication credentials, or a stateful cookie.
//
// Responses with the 429 status code MUST NOT be stored by a cache.
//
// See https://tools.ietf.org/html/rfc6585#section-4 for details.
type TooManyRequestsError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e TooManyRequestsError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e TooManyRequestsError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e TooManyRequestsError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (TooManyRequestsError) HTTPCode() int { return http.StatusTooManyRequests }

// --

// RequestHeaderFieldsTooLargeError (431) indicates that the server is
// unwilling to process the request because its header fields are too large.
// The request MAY be resubmitted after reducing the size of the request
// header fields.
//
// It can be used both when the set of request header fields in total is
// too large, and when a single header field is at fault.  In the latter
// case, the response representation SHOULD specify which header field
// was too large.
//
// For example:
//
//    HTTP/1.1 431 Request Header Fields Too Large
//    Content-Type: text/html
//
//    <html>
//       <head>
//          <title>Request Header Fields Too Large</title>
//       </head>
//       <body>
//          <h1>Request Header Fields Too Large</h1>
//          <p>The "Example" header was too large.</p>
//       </body>
//    </html>
//
// Responses with the 431 status code MUST NOT be stored by a cache.
//
// See https://tools.ietf.org/html/rfc6585#section-5 for details.
type RequestHeaderFieldsTooLargeError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e RequestHeaderFieldsTooLargeError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e RequestHeaderFieldsTooLargeError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e RequestHeaderFieldsTooLargeError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (RequestHeaderFieldsTooLargeError) HTTPCode() int { return http.StatusRequestHeaderFieldsTooLarge }

// --

// UnavailableForLegalReasonsError (451) indicates that the server is
// denying access to the resource as a consequence of a legal demand.
//
// The server in question might not be an origin server.  This type of
// legal demand typically most directly affects the operations of ISPs
// and search engines.
//
// Responses using this status code SHOULD include an explanation, in
// the response body, of the details of the legal demand: the party
// making it, the applicable legislation or regulation, and what classes
// of person and resource it applies to.  For example:
//
//    HTTP/1.1 451 Unavailable For Legal Reasons
//    Link: <https://spqr.example.org/legislatione>; rel="blocked-by"
//    Content-Type: text/html
//
//    <html>
//     <head><title>Unavailable For Legal Reasons</title></head>
//     <body>
//      <h1>Unavailable For Legal Reasons</h1>
//      <p>This request may not be serviced in the Roman Province
//      of Judea due to the Lex Julia Majestatis, which disallows
//      access to resources hosted on servers deemed to be
//      operated by the People's Front of Judea.</p>
//     </body>
//    </html>
//
// The use of the 451 status code implies neither the existence nor
// nonexistence of the resource named in the request.  That is to say,
// it is possible that if the legal demands were removed, a request for
// the resource still might not succeed.
//
// Note that in many cases clients can still access the denied resource
// by using technical countermeasures such as a VPN or the Tor network.
//
// A 451 response is cacheable by default, i.e., unless otherwise
// indicated by the method definition or explicit cache controls; see
// RFC7234.
//
// See https://tools.ietf.org/html/rfc7725#section-3 for details.
type UnavailableForLegalReasonsError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e UnavailableForLegalReasonsError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e UnavailableForLegalReasonsError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e UnavailableForLegalReasonsError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (UnavailableForLegalReasonsError) HTTPCode() int { return http.StatusUnavailableForLegalReasons }

// --

// ServerError (500) is an alias for InternalServerError.
//
// Deprecated: Use InternalServerError instead.
type ServerError = InternalServerError

// InternalServerError (500) indicates that the server
// encountered an unexpected condition that prevented it from fulfilling
// the request.
//
// See https://tools.ietf.org/html/rfc7231#section-6.6.1 for details.
type InternalServerError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e InternalServerError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e InternalServerError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e InternalServerError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (InternalServerError) HTTPCode() int { return http.StatusInternalServerError }

// --

// NotImplementedError (501) indicates that the server does
// not support the functionality required to fulfill the request.  This
// is the appropriate response when the server does not recognize the
// request method and is not capable of supporting it for any resource.
//
// A 501 response is cacheable by default; i.e., unless otherwise
// indicated by the method definition or explicit cache controls (see
// Section 4.2.2 of RFC7234).
//
// See https://tools.ietf.org/html/rfc7231#section-6.6.2 for details.
type NotImplementedError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e NotImplementedError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e NotImplementedError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e NotImplementedError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (NotImplementedError) HTTPCode() int { return http.StatusNotImplemented }

// --

// BadGatewayError (502) indicates that the server, while
// acting as a gateway or proxy, received an invalid response from an
// inbound server it accessed while attempting to fulfill the request.
//
// See https://tools.ietf.org/html/rfc7231#section-6.6.3 for details.
type BadGatewayError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e BadGatewayError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e BadGatewayError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e BadGatewayError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (BadGatewayError) HTTPCode() int { return http.StatusBadGateway }

// --

// ServiceUnavailableError (503) indicates that the server
// is currently unable to handle the request due to a temporary overload
// or scheduled maintenance, which will likely be alleviated after some
// delay.  The server MAY send a Retry-After header field
// (Section 7.1.3 of RFC7231) to suggest an appropriate amount of time
// for the client to wait before retrying the request.
//
// Note: The existence of the 503 status code does not imply that a
// server has to use it when becoming overloaded.  Some servers might
// simply refuse the connection.
//
// See https://tools.ietf.org/html/rfc7231#section-6.6.4 for details.
type ServiceUnavailableError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e ServiceUnavailableError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e ServiceUnavailableError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e ServiceUnavailableError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (ServiceUnavailableError) HTTPCode() int { return http.StatusServiceUnavailable }

// --

// GatewayTimeoutError (504) indicates that the server, while
// acting as a gateway or proxy, received an invalid response from an
// inbound server it accessed while attempting to fulfill the request.
//
// See https://tools.ietf.org/html/rfc7231#section-6.6.5 for details.
type GatewayTimeoutError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e GatewayTimeoutError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e GatewayTimeoutError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e GatewayTimeoutError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (GatewayTimeoutError) HTTPCode() int { return http.StatusGatewayTimeout }

// --

// HTTPVersionNotSupportedError (505) indicates that the server,
// while acting as a gateway or proxy, did not receive a timely response
// from an upstream server it needed to access in order to complete the
// request.
//
// See https://tools.ietf.org/html/rfc7231#section-6.6.6 for details.
type HTTPVersionNotSupportedError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e HTTPVersionNotSupportedError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e HTTPVersionNotSupportedError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e HTTPVersionNotSupportedError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (HTTPVersionNotSupportedError) HTTPCode() int { return http.StatusHTTPVersionNotSupported }

// --

// VariantAlsoNegotiatesError (506) indicates that the server has an internal
// configuration error: the chosen variant resource is configured to
// engage in transparent content negotiation itself, and is therefore
// not a proper end point in the negotiation process.
//
// See https://tools.ietf.org/html/rfc2295#section-8.1 for details.
type VariantAlsoNegotiatesError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e VariantAlsoNegotiatesError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e VariantAlsoNegotiatesError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e VariantAlsoNegotiatesError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (VariantAlsoNegotiatesError) HTTPCode() int { return http.StatusVariantAlsoNegotiates }

// --

// InsufficientStorageError (507) means the method could not
// be performed on the resource because the server is unable to store
// the representation needed to successfully complete the request.  This
// condition is considered to be temporary.  If the request that
// received this status code was the result of a user action, the
// request MUST NOT be repeated until it is requested by a separate user
// action.
//
// See https://tools.ietf.org/html/rfc4918#section-11.5 for details.
type InsufficientStorageError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e InsufficientStorageError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e InsufficientStorageError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e InsufficientStorageError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (InsufficientStorageError) HTTPCode() int { return http.StatusInsufficientStorage }

// --

// LoopDetectedError (508) indicates that the server
// terminated an operation because it encountered an infinite loop while
// processing a request with "Depth: infinity".  This status indicates
// that the entire operation failed.
//
// See https://tools.ietf.org/html/rfc5842#section-7.2 for details.
type LoopDetectedError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e LoopDetectedError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e LoopDetectedError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e LoopDetectedError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (LoopDetectedError) HTTPCode() int { return http.StatusLoopDetected }

// --

// NotExtendedError (510) indicates that the policy for accessing the resource
// has not been met in the request.
// The server should send back all the information necessary
// for the client to issue an extended request. It is outside the scope
// of this specification to specify how the extensions inform the
// client.
//
// If the 510 response contains information about extensions that were
// not present in the initial request then the client MAY repeat the
// request if it has reason to believe it can fulfill the extension
// policy by modifying the request according to the information provided
// in the 510 response. Otherwise the client MAY present any entity
// included in the 510 response to the user, since that entity may
// include relevant diagnostic information.
//
// See https://tools.ietf.org/html/rfc2774#section-7 for details.
type NotExtendedError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e NotExtendedError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e NotExtendedError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e NotExtendedError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (NotExtendedError) HTTPCode() int { return http.StatusNotExtended }

// --

// NetworkAuthenticationRequiredError (511) indicates that the client
// needs to authenticate to gain network access.
//
// The response representation SHOULD contain a link to a resource that
// allows the user to submit credentials (e.g., with an HTML form).
//
// Note that the 511 response SHOULD NOT contain a challenge or the
// login interface itself, because browsers would show the login
// interface as being associated with the originally requested URL,
// which may cause confusion.
//
// The 511 status SHOULD NOT be generated by origin servers; it is
// intended for use by intercepting proxies that are interposed as a
// means of controlling access to the network.
//
// Responses with the 511 status code MUST NOT be stored by a cache.
//
// See https://tools.ietf.org/html/rfc6585#section-6 for details.
type NetworkAuthenticationRequiredError struct {
	Message string
	Details []string
	Err     error
}

// Error returns the error in text form.
func (e NetworkAuthenticationRequiredError) Error() string {
	if s := e.Message; s != "" {
		return s
	}
	return http.StatusText(e.HTTPCode())
}

// Unwrap returns the inner error which might be nil.
func (e NetworkAuthenticationRequiredError) Unwrap() error {
	return e.Err
}

// ErrorDetails returns additional details about the error.
func (e NetworkAuthenticationRequiredError) ErrorDetails() []string {
	return e.Details
}

// HTTPCode returns the HTTP status code of the error.
func (NetworkAuthenticationRequiredError) HTTPCode() int {
	return http.StatusNetworkAuthenticationRequired
}

// --

// Following errors are opinionated and used by httputil functions.
// E.g. the MissingParameterError is a BadRequestError internally.

// --

// MissingParameterError indicates that a required parameter is missing
// or blank.
//
// Internally, this error is represented with HTTP status code
// 400 (Bad Request).
type MissingParameterError string

// Error returns the error in text form.
func (e MissingParameterError) Error() string {
	return fmt.Sprintf("Missing parameter %q", string(e))
}

// Unwrap ensures that the following condition holds true:
//
//   err := MissingParameterError("name")
//   errors.Is(err, BadRequestError{})
func (e MissingParameterError) Unwrap() error {
	return BadRequestError{Message: e.Error()}
}

// Is returns true if the given error is a descendant of MissingParameterError.
func (e MissingParameterError) Is(err error) bool {
	switch err.(type) {
	default:
		return false
	case BadRequestError, *BadRequestError:
		return true
	}
}

// HTTPCode returns the HTTP status code of the error.
func (MissingParameterError) HTTPCode() int { return http.StatusBadRequest }

// --

// InvalidParameterError indicates that a parameter is invalid.
//
// Internally, this error is represented with HTTP status code
// 400 (Bad Request).
type InvalidParameterError string

// Error returns the error in text form.
func (e InvalidParameterError) Error() string {
	return fmt.Sprintf("Invalid parameter %q", string(e))
}

// Unwrap ensures that the following condition holds true:
//
//   err := InvalidParameterError("name")
//   errors.Is(err, BadRequestError{})
func (e InvalidParameterError) Unwrap() error {
	return BadRequestError{Message: e.Error()}
}

// Is returns true if the given error is a descendant of InvalidParameterError.
func (e InvalidParameterError) Is(err error) bool {
	switch err.(type) {
	default:
		return false
	case BadRequestError, *BadRequestError:
		return true
	}
}

// HTTPCode returns the HTTP status code of the error.
func (InvalidParameterError) HTTPCode() int { return http.StatusBadRequest }
