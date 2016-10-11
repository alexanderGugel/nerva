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
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestRespondJSONCode(t *testing.T) {
    w := httptest.NewRecorder()
    wantCode := http.StatusTeapot
    RespondJSON(w, wantCode, nil)

    if gotCode := w.Code; gotCode != wantCode {
        t.Errorf("w.Code = %v want %v", gotCode, wantCode)
    }
}

func TestRespondJSONContentType(t *testing.T) {
    w := httptest.NewRecorder()
    RespondJSON(w, http.StatusTeapot, nil)

    wantContentType := "application/json; charset=utf-8"
    if gotContentType := w.Header().Get("Content-Type"); gotContentType != wantContentType {
        t.Errorf("w.Header().Get(\"Content-Type\") = %v want %v", gotContentType, wantContentType)
    }
}

func TestRespondJSONSuccess(t *testing.T) {
    w := httptest.NewRecorder()
    if err := RespondJSON(w, http.StatusTeapot, nil); err != nil {
        t.Errorf("RespondJSON(%q, %q, %q) should not err: %v", w, http.StatusTeapot, nil, err)
    }
}

func TestRespondJSONNonStringKeys(t *testing.T) {
    w := httptest.NewRecorder()
    data := map[int]string{1: ""}
    if err := RespondJSON(w, http.StatusTeapot, data); err == nil {
        t.Errorf("RespondJSON(%q, %q, %q) should err", w, http.StatusTeapot, nil)
    }
}

func TestRespondJSONBody(t *testing.T) {
    w := httptest.NewRecorder()
    data := map[string]string{"hello": "world"}
    wantBody := "{\"hello\":\"world\"}"

    RespondJSON(w, http.StatusTeapot, data)

    if gotBody := w.Body.String(); gotBody != wantBody {
        t.Errorf("w.Body.String() = %v want %v", gotBody, wantBody)
    }
}
