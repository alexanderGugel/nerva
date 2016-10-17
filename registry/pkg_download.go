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
	"github.com/julienschmidt/httprouter"
	"github.com/libgit2/git2go"
	"net/http"
)

// HandlePkgDownload handles package downloads.
func (r *Registry) HandlePkgDownload(repo *git.Repository,
	w http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
	version := ps.ByName("version")

	id, err := git.NewOid(version)
	if err != nil || id == nil {
		code := http.StatusBadRequest
		res := &util.ErrorResponse{
			http.StatusText(code),
			"version is not a valid git object id",
		}
		return util.RespondJSON(w, code, res)
	}

	d, err := storage.NewDownload(repo, id)
	if err != nil || d == nil {
		if gitErr, ok := err.(*git.GitError); !ok ||
			gitErr.Class != git.ErrClassOdb {
			return err
		}
		code := http.StatusNotFound
		res := &util.ErrorResponse{
			http.StatusText(code),
			"package not found",
		}
		return util.RespondJSON(w, code, res)
	}
	return d.Start(w)
}
