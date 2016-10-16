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
	"github.com/alexanderGugel/nerva/util"
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

	return &Upstream{
		URL:    urlURL,
		Client: &http.Client{},
	}, nil
}

// Ping checks it the upstream registry can be reached.
func (u *Upstream) Ping() error {
	res, err := http.Get(u.URL.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
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
	defer res.Body.Close()

	_, err = io.Copy(w, res.Body)
	return err
}

// UpstreamStatus represents the response to a request to the /upstreams
// endpoint.
type UpstreamStatus struct {
	URL    string
	Status string
}

// GetStatus returns the current status of the upstream registry.
func (u *Upstream) GetStatus() *UpstreamStatus {
	pingErr := u.Ping()
	var status string
	switch pingErr {
	case nil:
		status = "up"
	default:
		status = "down"
	}

	return &UpstreamStatus{
		URL:    u.URL.String(),
		Status: status,
	}
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

// HandleUpstreams retrieves the current memory stats.
func (r *Registry) HandleUpstreams(w http.ResponseWriter, req *http.Request,
	ps httprouter.Params) error {
	name := ps.ByName("name")

	if name != "-" {
		r.Router.NotFound.ServeHTTP(w, req)
		return nil
	}
	return util.RespondJSON(w, 200, r.Upstream.GetStatus())
}
