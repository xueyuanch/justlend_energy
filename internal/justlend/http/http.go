package http

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"justlend/internal/derrors"
	"net/http"
	"strconv"
	"strings"
)

type Failer interface {
	Failed() error
	Data() interface{}
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(Failer); ok && e.Failed() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		return e.Failed()
	} else {
		c, _ := derrors.ToCode(nil)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		return json.NewEncoder(w).Encode(map[string]interface{}{"code": c, "data": e.Data()})
	}
}

// encodeError is the common method to encode response types that
// errors occurs during handling the request.
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	c, message := derrors.ToCode(err) // Only message is required.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(derrors.ToHTTPSta(err))
	json.NewEncoder(w).Encode(map[string]interface{}{"code": c, "error": message})
}

func encodeResponsePass(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(Failer); ok && e.Failed() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		return e.Failed()
	} else {
		c, _ := derrors.ToCode(nil)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		return json.NewEncoder(w).Encode(map[string]interface{}{"code": c, "data": e.Data()})
	}
}

// safeExtractUint safely extract the variable of type uint64 from the request
// route whose name is given in param, or return 0 as the default value if
// an error occurs during extracting.
func safeExtractUint(r *http.Request, name string) (v uint64) {
	v, _ = strconv.ParseUint(mux.Vars(r)[name], 10, 64)
	return
}

// maxQueryFieldValueLen signifies the maximum string field value length.
const maxQueryFieldValueLen = 128

// safeExtractQueryString safely extract the variable of type string from the
// request query fields whose name is given in param, or return empty string
// as the default value if the field not present in the request query of if
// an error occurs during extracting.
//
// If the field value character length exceeds maxQueryFieldValueLen, the
// field value will be truncated and the first maxQueryFieldValueLen
// characters will be returned.
func safeExtractQueryString(r *http.Request, name string) (v string) {
	if v = r.FormValue(name); len(v) > maxQueryFieldValueLen {
		v = v[:maxQueryFieldValueLen]
	}
	return
}

func safeExtractIds(r *http.Request, name string) []uint64 {
	return parseIDs(r.FormValue(name))
}

func parseIDs(roomIDsStr string) []uint64 {
	roomIDs := strings.Split(roomIDsStr, ",")
	var ids []uint64
	for _, idStr := range roomIDs {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return nil
		}
		ids = append(ids, id)
	}
	return ids
}

// safeExtractQueryBoolPtr safely extract the variable of type pointer to a bool from the
// request query fields whose name is given in param. It returns a nil if the incoming
// value is neither ‘t’ nor ‘f’ which is equivalent to the field not being set.
func safeExtractQueryBoolPtr(r *http.Request, name string) *bool {
	if v, b := r.FormValue(name), true; v == "t" {
		return &b
	} else if b = false; v == "f" {
		return &b
	}
	return nil
}

// safeExtractQueryUint safely extract the variable of type uint64 from the request
// query fields whose name is given in param, or return the default value
// of 0 if an error occurs during extracting.
func safeExtractQueryUint(r *http.Request, name string) (v uint64) {
	v, _ = strconv.ParseUint(r.FormValue(name), 10, 64)
	return
}

func safeExtractQueryInt(r *http.Request, name string) (v int64) {
	v, _ = strconv.ParseInt(r.FormValue(name), 10, 64)
	return
}

func safeExtractInt(r *http.Request, name string) (v int64) {
	// id == 0 if error occurs during parsing.
	v, _ = strconv.ParseInt(mux.Vars(r)[name], 10, 64)
	return
}

// safeExtractQueryFloat extracts a float value from the query parameters of an HTTP request.
// It returns the extracted float value and ignores any parsing errors.
func safeExtractQueryFloat(r *http.Request, name string) (v float64) {
	v, _ = strconv.ParseFloat(r.FormValue(name), 64)
	return
}

// safeExtractQueryUintPtr safely extract the variable of type a pointer to a uint64
// from the request query fields whose name is given in param. It returns a nil if the
// incoming value does not meet the expected format or not present in query fields.
func safeExtractQueryUintPtr(r *http.Request, name string) *uint64 {
	v, err := strconv.ParseUint(r.FormValue(name), 10, 64)
	// Return nil if received value is not a valid uint64 which
	// signifies extracted field not set.
	if err != nil {
		return nil
	}
	return &v
}
