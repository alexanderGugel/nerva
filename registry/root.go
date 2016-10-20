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
	"github.com/alexanderGugel/nerva/storage"
	"github.com/alexanderGugel/nerva/util"
	"net/http"
	"path"
)

// Root maps package names to package root descriptors. In this context,
// package root descriptors are URLs to package root documents.
// The root URL is the base of the package registry. Given this url, a name, and
// a version, a package can be uniquely identified, assuming it exists in the
// registry.
// When requested, the registry root URL SHOULD return a list of packages in the
// registry in the form of a hash of package names to package root descriptors.
// The package root descriptor MUST be either: an Object that would be valid for
// the “package root url” contents for the named package, or a string URL that
// should be used as the package root url.
// In the case of a string URL, it MAY refer to a different registry. In that
// case, a request for {registry root url}/{package name} SHOULD be EITHER a 301
// or 302 redirect to the same URL as named in the string value, OR a valid
// “package root url” response.
// See http://wiki.commonjs.org/wiki/Packages/Registry#registry_root_url
type Root map[string]string

// NewRoot creates a new CommonJS registry root document from a given
// storage directory by reading in the repositories that are available in the
// storage dir.
func NewRoot(storage *storage.Storage, url string) (*Root, error) {
	root := Root{}
	names, _ := storage.Ls()
	for _, name := range names {
		root[name] = path.Join(url, name)
	}
	return &root, nil
}

// HandleRoot handles requests to the registry root URL.
// The root URL is the base of the package registry. Given this url, a name, and
// a version, a package can be uniquely identified, assuming it exists in the
// registry.
// When requested, the registry root URL SHOULD return a list of packages in the
// registry in the form of a hash of package names to package root descriptors.
// The package root descriptor MUST be either: an Object that would be valid for
// the “package root url” contents for the named package, or a string URL that
// should be used as the package root url.
// In the case of a string URL, it MAY refer to a different registry. In that
// case, a request for {registry root url}/{package name} SHOULD be EITHER a 301
// or 302 redirect to the same URL as named in the string value, OR a valid
// “package root url” response.
// See http://wiki.commonjs.org/wiki/Packages/Registry#registry_root_url
func (r *Registry) HandleRoot(w http.ResponseWriter, req *http.Request) error {
	res, err := NewRoot(r.storage, r.config.FrontAddr)
	if err != nil {
		return err
	}
	return util.RespondJSON(w, 200, res)
}
