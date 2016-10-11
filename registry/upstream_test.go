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

package registry

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func createUpstream(url string, t *testing.T) *Upstream {
    upstream, err := NewUpstream(url)
    if err != nil {
        t.Errorf("NewUpstream(%q) unexpected err: %q", url, err)
    }
    return upstream
}

func TestNewUpstream(t *testing.T) {
    url := "http://registry.npmjs.com"
    upstream := createUpstream(url, t)
    if upstream.URL.String() != url {
        t.Errorf("upstream.URL = %q want %q", url, upstream.URL.String(), url)
    }
}

func TestUpstreamRedirectCode(t *testing.T) {
    url := "http://registry.npmjs.com"
    upstream := createUpstream(url, t)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "", nil)
    upstream.RedirectPackageRoot("tape", w, req)

    wantCode := http.StatusMovedPermanently
    if gotCode := w.Code; gotCode != wantCode {
        t.Errorf("w.Code = %q, want %q", gotCode, wantCode)
    }
}

func TestUpstreamRedirectHeader(t *testing.T) {
    url := "http://registry.npmjs.com"
    upstream := createUpstream(url, t)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "", nil)
    upstream.RedirectPackageRoot("tape", w, req)

    wantLocation := "http://registry.npmjs.com/tape"
    if gotLocation := w.Header().Get("Location"); gotLocation != wantLocation {
        t.Errorf("w.Header().Get(\"Location\") = %q, want %q", gotLocation, wantLocation)
    }
}
