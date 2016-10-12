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

// Package registry implements a CommonJS compliant package registry.
// See http://wiki.commonjs.org/wiki/Packages/Registry
package registry

import (
	"github.com/alexanderGugel/nerva/storage"
	"github.com/alexanderGugel/nerva/util"
	"github.com/julienschmidt/httprouter"
	"github.com/libgit2/git2go"
	"net/http"
)

// Registry represents an Common JS registry server. A Registry does exposes a
// router, which can be bound to an arbitrary socket.
type Registry struct {
	Router   *httprouter.Router
	Storage  *storage.Storage
	Upstream *Upstream
	ShaCache *storage.ShaCache
	Config   Config
}

// New create a new CommonJS registry.
func New(config Config) (*Registry, error) {
	upstream, err := NewUpstream(config.UpstreamURL)
	if err != nil {
		return nil, err
	}

	shaCache, err := storage.NewShaCache(config.ShaCacheSize)
	if err != nil {
		return nil, err
	}

	registry := &Registry{
		Router:   httprouter.New(),
		Storage:  storage.New(config.StorageDir),
		Upstream: upstream,
		ShaCache: shaCache,
		Config:   config,
	}
	registry.attachRoutes()

	return registry, nil
}

func (r *Registry) isTLSEnabled() bool {
	return r.Config.CertFile != "" && r.Config.KeyFile != ""
}

func (r *Registry) getScheme() string {
	if r.isTLSEnabled() {
		return "https"
	}
	return "http"
}

// Start starts the registry.
func (r *Registry) Start() error {
	c := r.Config
	if r.isTLSEnabled() {
		return http.ListenAndServeTLS(c.Addr, c.CertFile, c.KeyFile, r.Router)
	}
	return http.ListenAndServe(c.Addr, r.Router)
}

func (r *Registry) attachRoutes() {
	r.Router.GET("/", util.ErrHandler(r.HandleRoot))

	r.Router.GET("/:name", util.ErrHandler(r.repoHandler(r.HandlePackageRoot)))
	r.Router.GET("/:name/-/:version", util.ErrHandler(r.repoHandler(r.HandlePackageDownload)))

	r.Router.GET("/:name/ping", util.ErrHandler(r.HandlePing))
	r.Router.GET("/:name/stats", util.ErrHandler(r.HandleStats))
}

type repoHandle func(repo *git.Repository, w http.ResponseWriter,
	req *http.Request, ps httprouter.Params) error

func (r *Registry) repoHandler(handle repoHandle) util.ErrHandle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
		name := ps.ByName("name")

		if !util.IsValid(name) {
			res := &util.ErrorResponse{"bad request", "invalid name"}
			return util.RespondJSON(w, http.StatusBadRequest, res)
		}

		repo, err := r.Storage.GetRepo(name)
		if err == nil {
			return handle(repo, w, req, ps)
		}

		if gitErr, ok := err.(*git.GitError); !ok ||
			gitErr.Class != git.ErrClassOs {
			return err
		}

		if r.Upstream == nil || r.Upstream.URL == nil {
			res := &util.ErrorResponse{"not found", "package not found"}
			return util.RespondJSON(w, http.StatusNotFound, res)
		}

		r.Upstream.RedirectPackageRoot(name, w, req)
		return nil
	}
}
