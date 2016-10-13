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
	"github.com/hashicorp/golang-lru"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"net/url"
	"path"
)

// Upstream represents an external registry. It provides a caching layer for
// frequently requested packages.
type Upstream struct {
	URL    *url.URL
	Client *http.Client
}

// NewUpstream instantiates a new registry proxy.
func NewUpstream(rootURL string) (*Upstream, error) {
	urlURL, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	return &Upstream{urlURL, client}, nil
}

// HandleReq redirects the client to the package root of the package
// with the specified name.
func (u *Upstream) HandleReq(w http.ResponseWriter, req *http.Request,
	ps httprouter.Params) error {

	url := *u.URL
	url.Path = path.Join(url.Path, req.URL.Path)

	req, err := http.NewRequest(req.Method, url.String(), req.Body)
	if err != nil {
		return err
	}

	res, err := u.Client.Do(req)
	if err != nil {
		return err
	}
	copyHeader(w.Header(), res.Header)

	_, err = io.Copy(w, res.Body)
	return err
}

// copyHeader copies header pairs from one header to another.
// See https://golang.org/src/net/http/httputil/reverseproxy.go
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
