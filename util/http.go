// Copyright © 2016 Alexander Gugel <alexander.gugel@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package util

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// ErrorResponse is a HTTP response.
type ErrorResponse struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

// RespondJSON encodes the passed in data as JSON-compatible string and returns
// it to the client.
func RespondJSON(w http.ResponseWriter, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	res, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	_, err = w.Write(res)
	return err
}

// ErrHandle is a custom HTTP handle that can optionally return an error.
type ErrHandle func(w http.ResponseWriter, r *http.Request,
	ps httprouter.Params) error

// ErrHandler allows custom handlers to return internal server errors.
func ErrHandler(handler ErrHandle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		if err := handler(w, req, ps); err != nil {
			contextLog := log.WithFields(GetRequestFields(req))

			LogErr(contextLog, err, "handler failed")
			res := &ErrorResponse{"internal server error", "unexpected internal error"}
			if err := RespondJSON(w, http.StatusInternalServerError, res); err != nil {
				LogErr(contextLog, err, "failed to write response")
			}
		}
	}
}

// ValidatePropHandler validates the passed in property names.
func ValidatePropHandler(name string, handle ErrHandle) ErrHandle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
		prop := ps.ByName(name)
		if !IsValid(prop) {
			fields := log.Fields{}
			fields[name] = prop
			res := &ErrorResponse{"bad request", "invalid " + name}
			return RespondJSON(w, http.StatusBadRequest, res)
		}
		return handle(w, req, ps)
	}
}
