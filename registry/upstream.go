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
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
	"path"
)

// Upstream represents an external registry. It provides a caching layer for
// frequently requested packages.
type Upstream struct {
	URL *url.URL
}

// NewUpstream instantiates a new registry proxy.
func NewUpstream(rootURL string) (*Upstream, error) {
	u, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}
	return &Upstream{u}, nil
}

// Proxy proxies a specific document hosted on the upstream registry.
func (u *Upstream) Proxy(w http.ResponseWriter, req *http.Request,
	ps httprouter.Params) error {

	url := *u.URL
	url.Path = path.Join(url.Path, req.URL.Path)

	return nil
}

// RedirectPackageRoot redirects the client to the package root of the package
// with the specified name.
func (u *Upstream) RedirectPackageRoot(name string, w http.ResponseWriter,
	req *http.Request) {
	url := *u.URL
	url.Path = path.Join(url.Path, name)
	http.Redirect(w, req, url.String(), http.StatusMovedPermanently)
}
