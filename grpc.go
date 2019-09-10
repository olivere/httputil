// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil

/*
import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// grpcToHttpMap maps a gRPC status code to a HTTP status code following the
// mapping described at https://cloud.google.com/apis/design/errors#handling_errors.
var grpcToHttpMap = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           499,
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusBadRequest,
	codes.Aborted:            http.StatusConflict,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DataLoss:           http.StatusInternalServerError,
	codes.Unauthenticated:    http.StatusUnauthorized,
}

// HTTPCodeFromGRPCCode returns a HTTP status code from a gRPC status code following the
// mapping described at https://cloud.google.com/apis/design/errors#handling_errors.
//
//   gRPC               | HTTP
//   -------------------+-----------
//   OK                 | 200 OK
//   InvalidArgument    | 400 Bad Request
//   FailedPrecondition | 400 Bad Request
//   OutOfRange         | 400 Bad Request
//   Unauthenticated    | 401 Unauthorized
//   PermissionDenied   | 403 Forbidden
//   NotFound           | 404 Not Found
//   Aborted            | 409 Conflict
//   AlreadyExists      | 409 Conflict
//   ResourceExhausted  | 429 Too Many Requests
//   Cancelled          | 499 Cancelled
//   DataLoss           | 500 Internal Server Error
//   Unknown            | 500 Internal Server Error
//   Internal           | 500 Internal Server Error
//   Unimplemented      | 501 Not Implemented
//   Unavailable        | 503 Service Unavailable
//   DeadlineExceeded   | 504 Gateway Timeout
//   Other              | 500 Internal Server Error
//
func HTTPCodeFromGRPCCode(code codes.Code) int {
	if c, ok := grpcToHttpMap[code]; ok {
		return c
	}
	return http.StatusInternalServerError
}

// httpToGrpcMap maps a HTTP status code to a gRPC status code following the
// mapping described at https://cloud.google.com/apis/design/errors#handling_errors.
var httpToGrpcMap = map[int]codes.Code{
	http.StatusOK:                  codes.OK,
	http.StatusBadRequest:          codes.InvalidArgument,
	http.StatusUnauthorized:        codes.Unauthenticated,
	http.StatusForbidden:           codes.PermissionDenied,
	http.StatusNotFound:            codes.NotFound,
	http.StatusConflict:            codes.AlreadyExists,
	http.StatusTooManyRequests:     codes.ResourceExhausted,
	499:                            codes.Canceled,
	http.StatusInternalServerError: codes.Internal,
	http.StatusNotImplemented:      codes.Unimplemented,
	http.StatusServiceUnavailable:  codes.Unavailable,
	http.StatusGatewayTimeout:      codes.DeadlineExceeded,
}

// GRPCCodeFromHTTPCode returns a gRPC status code from a HTTP status code following the
// mapping described at https://cloud.google.com/apis/design/errors#handling_errors.
//
//   HTTP                      | gRPC
//   --------------------------+-----------------
//   200 OK                    | OK
//   400 Bad Request           | InvalidArgument
//   401 Unauthorized          | Unauthenticated
//   403 Forbidden             | PermissionDenied
//   404 Not Found             | NotFound
//   409 Conflict              | AlreadyExists
//   429 Too Many Requests     | ResourceExhausted
//   499 Cancelled             | Cancelled
//   500 Internal Server Error | Internal
//   501 Not Implemented       | Unimplemented
//   503 Service Unavailable   | Unavailable
//   504 Gateway Timeout       | DeadlineExceeded
//   Other                     | Unknown
//
func GRPCCodeFromHTTPCode(code int) codes.Code {
	if c, ok := httpToGrpcMap[code]; ok {
		return c
	}
	return codes.Unknown
}
*/
