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
	"net/url"
	"path"
)

// Upstream represents an external registry.
type Upstream struct {
	URL *url.URL
}

// NewUpstream instantiates a new registry proxy.
func NewUpstream(rootURL string) (*Upstream, error) {
	urlURL, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}
	return &Upstream{urlURL}, nil
}

// RedirectPackageRoot redirects the client to the package root of the package
// with the specified name.
func (u *Upstream) RedirectPackageRoot(name string, w http.ResponseWriter,
	req *http.Request) {
	url := *u.URL
	url.Path = path.Join(url.Path, name)
	http.Redirect(w, req, url.String(), 301)
}
