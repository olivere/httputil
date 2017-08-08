// Copyright 2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package httputil

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// -- FormValue --

// MustFormString checks if the request r has a Form value with
// the specified key of type string. If is doesn't, it will panic.
func MustFormString(r *http.Request, key string) string {
	v := r.FormValue(key)
	if v == "" {
		panic(MissingParameterError(key))
	}
	return v
}

// MustFormBool checks if the request r has a Form value with
// the specified key that can be converted to a bool.
// If is doesn't, it will panic.
func MustFormBool(r *http.Request, key string) bool {
	v := r.FormValue(key)
	if v == "" {
		panic(MissingParameterError(key))
	}
	f, err := strconv.ParseBool(v)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return f
}

// MustFormInt checks if the request r has a Form value with
// the specified key that can be converted to an int.
// If is doesn't, it will panic.
func MustFormInt(r *http.Request, key string) int {
	v := r.FormValue(key)
	if v == "" {
		panic(MissingParameterError(key))
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return i
}

// MustFormInt32 checks if the request r has a Form value with
// the specified key that can be converted to an int32.
// If is doesn't, it will panic.
func MustFormInt32(r *http.Request, key string) int32 {
	v := r.FormValue(key)
	if v == "" {
		panic(MissingParameterError(key))
	}
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return int32(i)
}

// MustFormInt64 checks if the request r has a Form value with
// the specified key that can be converted to an int64.
// If is doesn't, it will panic.
func MustFormInt64(r *http.Request, key string) int64 {
	v := r.FormValue(key)
	if v == "" {
		panic(MissingParameterError(key))
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return i
}

// MustFormFloat64 checks if the request r has a Form value with
// the specified key that can be converted to a float64.
// If is doesn't, it will panic.
func MustFormFloat64(r *http.Request, key string) float64 {
	v := r.FormValue(key)
	if v == "" {
		panic(MissingParameterError(key))
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return f
}

// FormString checks if the request r has a Form value with
// the specified key. If is doesn't, it will return defaultValue.
func FormString(r *http.Request, key string, defaultValue string) string {
	if v := r.FormValue(key); v != "" {
		return v
	}
	return defaultValue
}

// FormBool checks if the request r has a Form value with
// the specified key that can be converted to a bool.
// If is doesn't, it will return defaultValue.
func FormBool(r *http.Request, key string, defaultValue bool) bool {
	v := r.FormValue(key)
	if v == "" {
		return defaultValue
	}
	f, err := strconv.ParseBool(v)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return f
}

// FormInt checks if the request r has a Form value with
// the specified key that can be converted to an int.
// If is doesn't, it will return defaultValue.
func FormInt(r *http.Request, key string, defaultValue int) int {
	v := r.FormValue(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return i
}

// FormInt32 checks if the request r has a Form value with
// the specified key that can be converted to an int32.
// If is doesn't, it will return defaultValue.
func FormInt32(r *http.Request, key string, defaultValue int32) int32 {
	v := r.FormValue(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return int32(i)
}

// FormInt64 checks if the request r has a Form value with
// the specified key that can be converted to an int64.
// If is doesn't, it will return defaultValue.
func FormInt64(r *http.Request, key string, defaultValue int64) int64 {
	v := r.FormValue(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return i
}

// FormFloat64 checks if the request r has a Form value with
// the specified key that can be converted to a float64.
// If is doesn't, it will return defaultValue.
func FormFloat64(r *http.Request, key string, defaultValue float64) float64 {
	v := r.FormValue(key)
	if v == "" {
		return defaultValue
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return f
}

// -- Query string --

// MustQueryString checks if the request r has a query string with
// the specified key. If is doesn't, it will panic.
func MustQueryString(r *http.Request, key string) string {
	v := r.URL.Query().Get(key)
	if v == "" {
		panic(MissingParameterError(key))
	}
	return v
}

// MustQueryBool checks if the request r has a query string with
// the specified key that can be converted to a bool.
// If is doesn't, it will panic.
func MustQueryBool(r *http.Request, key string) bool {
	v := r.URL.Query().Get(key)
	if v == "" {
		panic(MissingParameterError(key))
	}
	f, err := strconv.ParseBool(v)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return f
}

// MustQueryInt checks if the request r has a query string with
// the specified key that can be converted to an int.
// If is doesn't, it will panic.
func MustQueryInt(r *http.Request, key string) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		panic(MissingParameterError(key))
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return i
}

// MustQueryInt32 checks if the request r has a query string with
// the specified key that can be converted to an int32.
// If is doesn't, it will panic.
func MustQueryInt32(r *http.Request, key string) int32 {
	v := r.URL.Query().Get(key)
	if v == "" {
		panic(MissingParameterError(key))
	}
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return int32(i)
}

// MustQueryInt64 checks if the request r has a query string with
// the specified key that can be converted to an int64.
// If is doesn't, it will panic.
func MustQueryInt64(r *http.Request, key string) int64 {
	v := r.URL.Query().Get(key)
	if v == "" {
		panic(MissingParameterError(key))
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return i
}

// MustQueryFloat64 checks if the request r has a query string with
// the specified key that can be converted to a float64.
// If is doesn't, it will panic.
func MustQueryFloat64(r *http.Request, key string) float64 {
	v := r.URL.Query().Get(key)
	if v == "" {
		panic(MissingParameterError(key))
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return f
}

// MustQueryTime checks if the request r has a query string with
// the specified key that can be converted to a time.Time, based on
// the given layout format.
// If is doesn't, it will return defaultValue or a zero time.
func MustQueryTime(r *http.Request, key, layout string) time.Time {
	v := r.URL.Query().Get(key)
	if v == "" {
		var t time.Time
		return t
	}
	t, err := time.Parse(layout, v)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return t
}

// MustQueryTimeWithDefault checks if the request r has a query string with
// the specified key that can be converted to a time.Time, based on
// the given layout format.
// If is doesn't, it will return defaultValue or a zero time.
func MustQueryTimeWithDefault(r *http.Request, key, layout string, defaultValue time.Time) time.Time {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultValue
	}
	t, err := time.Parse(layout, v)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return t
}

// MustQueryDuration checks if the request r has a query string with
// the specified key that can be converted to a time.Duration.
// If is doesn't, it will return defaultValue or a zero time.
func MustQueryDuration(r *http.Request, key string) time.Duration {
	v := r.URL.Query().Get(key)
	if v == "" {
		var d time.Duration
		return d
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return d
}

// MustQueryDurationWithDefault checks if the request r has a query string with
// the specified key that can be converted to a time.Duration.
// If is doesn't, it will return defaultValue or a zero time.
func MustQueryDurationWithDefault(r *http.Request, key string, defaultValue time.Duration) time.Duration {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return d
}

// QueryString checks if the request r has a query string with
// the specified key. If is doesn't, it will return defaultValue.
func QueryString(r *http.Request, key string, defaultValue string) string {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultValue
	}
	return v
}

// QueryStringArray checks if the request r has a query string with
// the specified key. If is doesn't, it will return defaultValue.
// Otherwise it'll split the string by a comma and return the resulting array.
func QueryStringArray(r *http.Request, key string, defaultValue []string) []string {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultValue
	}
	return strings.Split(v, ",")
}

// QueryBool checks if the request r has a query string with
// the specified key that can be converted to a bool.
// If is doesn't, it will return defaultValue.
func QueryBool(r *http.Request, key string, defaultValue bool) bool {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultValue
	}
	f, err := strconv.ParseBool(v)
	if err != nil {
		return defaultValue
	}
	return f
}

// QueryInt checks if the request r has a query string with
// the specified key that can be converted to an int.
// If is doesn't, it will return defaultValue.
func QueryInt(r *http.Request, key string, defaultValue int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return i
}

// QueryInt32 checks if the request r has a query string with
// the specified key that can be converted to an int32.
// If is doesn't, it will return defaultValue.
func QueryInt32(r *http.Request, key string, defaultValue int32) int32 {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int32(i)
}

// QueryInt64 checks if the request r has a query string with
// the specified key that can be converted to an int64.
// If is doesn't, it will return defaultValue.
func QueryInt64(r *http.Request, key string, defaultValue int64) int64 {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

// QueryFloat64 checks if the request r has a query string with
// the specified key that can be converted to a float64.
// If is doesn't, it will return defaultValue.
func QueryFloat64(r *http.Request, key string, defaultValue float64) float64 {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

// QueryTime checks if the request r has a query string with
// the specified key that can be converted to a time.Time, based on
// the given layout format.
// If is doesn't, it will return defaultValue or a zero time.
func QueryTime(r *http.Request, key, layout string) time.Time {
	v := r.URL.Query().Get(key)
	if v == "" {
		var t time.Time
		return t
	}
	t, err := time.Parse(layout, v)
	if err != nil {
		var t time.Time
		return t
	}
	return t
}

// QueryTimeWithDefault checks if the request r has a query string with
// the specified key that can be converted to a time.Time, based on
// the given layout format.
// If is doesn't, it will return defaultValue or a zero time.
func QueryTimeWithDefault(r *http.Request, key, layout string, defaultValue time.Time) time.Time {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultValue
	}
	t, err := time.Parse(layout, v)
	if err != nil {
		return defaultValue
	}
	return t
}

// QueryDuration checks if the request r has a query string with
// the specified key that can be converted to a time.Duration.
// If is doesn't, it will return defaultValue or a zero duration.
func QueryDuration(r *http.Request, key string) time.Duration {
	v := r.URL.Query().Get(key)
	if v == "" {
		var d time.Duration
		return d
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		var d time.Duration
		return d
	}
	return d
}

// QueryDurationWithDefault checks if the request r has a query string with
// the specified key that can be converted to a time.Duration.
// If is doesn't, it will return defaultValue or a zero duration.
func QueryDurationWithDefault(r *http.Request, key string, defaultValue time.Duration) time.Duration {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return defaultValue
	}
	return d
}

// -- Router parameters --

// MustParamsString checks if the request r has a routing component with
// the specified key. If is doesn't, it will panic.
func MustParamsString(r *http.Request, key string) string {
	vars := mux.Vars(r)
	v, found := vars[key]
	if !found || v == "" {
		panic(MissingParameterError(key))
	}
	return v
}

// MustParamsBool checks if the request r has a routing component with
// the specified key that can be converted to a bool.
// If is doesn't, it will panic.
func MustParamsBool(r *http.Request, key string) bool {
	vars := mux.Vars(r)
	v, found := vars[key]
	if !found || v == "" {
		panic(MissingParameterError(key))
	}
	f, err := strconv.ParseBool(v)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return f
}

// MustParamsInt checks if the request r has a routing component with
// the specified key that can be converted to an int.
// If is doesn't, it will panic.
func MustParamsInt(r *http.Request, key string) int {
	vars := mux.Vars(r)
	v, found := vars[key]
	if !found || v == "" {
		panic(MissingParameterError(key))
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return i
}

// MustParamsInt32 checks if the request r has a routing component with
// the specified key that can be converted to an int32.
// If is doesn't, it will panic.
func MustParamsInt32(r *http.Request, key string) int32 {
	vars := mux.Vars(r)
	v, found := vars[key]
	if !found || v == "" {
		panic(MissingParameterError(key))
	}
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return int32(i)
}

// MustParamsInt64 checks if the request r has a routing component with
// the specified key that can be converted to an int64.
// If is doesn't, it will panic.
func MustParamsInt64(r *http.Request, key string) int64 {
	vars := mux.Vars(r)
	v, found := vars[key]
	if !found || v == "" {
		panic(MissingParameterError(key))
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return i
}

// MustParamsFloat64 checks if the request r has a routing component with
// the specified key that can be converted to a float64.
// If is doesn't, it will panic.
func MustParamsFloat64(r *http.Request, key string) float64 {
	vars := mux.Vars(r)
	v, found := vars[key]
	if !found || v == "" {
		panic(MissingParameterError(key))
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return f
}

// ParamsString checks if the request r has a routing component with
// the specified key. If is doesn't, it will return defaultValue.
func ParamsString(r *http.Request, key string, defaultValue string) string {
	vars := mux.Vars(r)
	v, found := vars[key]
	if !found || v == "" {
		return defaultValue
	}
	return v
}

// ParamsInt checks if the request r has a routing component with
// the specified key that can be converted to an int.
// If is doesn't, it will return defaultValue.
func ParamsInt(r *http.Request, key string, defaultValue int) int {
	vars := mux.Vars(r)
	v, found := vars[key]
	if !found || v == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return i
}

// ParamsInt32 checks if the request r has a routing component with
// the specified key that can be converted to an int32.
// If is doesn't, it will return defaultValue.
func ParamsInt32(r *http.Request, key string, defaultValue int32) int32 {
	vars := mux.Vars(r)
	v, found := vars[key]
	if !found || v == "" {
		return defaultValue
	}
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return int32(i)
}

// ParamsInt64 checks if the request r has a routing component with
// the specified key that can be converted to an int64.
// If is doesn't, it will return defaultValue.
func ParamsInt64(r *http.Request, key string, defaultValue int64) int64 {
	vars := mux.Vars(r)
	v, found := vars[key]
	if !found || v == "" {
		return defaultValue
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return i
}

// ParamsFloat64 checks if the request r has a routing component with
// the specified key that can be converted to a float64.
// If is doesn't, it will return defaultValue.
func ParamsFloat64(r *http.Request, key string, defaultValue float64) float64 {
	vars := mux.Vars(r)
	v, found := vars[key]
	if !found || v == "" {
		return defaultValue
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic(InvalidParameterError(key))
	}
	return f
}
