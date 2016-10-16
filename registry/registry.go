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
	log "github.com/Sirupsen/logrus"
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
	Config   *Config
}

// New create a new CommonJS registry.
func New(config *Config) (*Registry, error) {
	if err := ValidateConfig(config); err != nil {
		return nil, err
	}
	registry := &Registry{Config: config}
	if err := registry.init(); err != nil {
		return nil, err
	}
	registry.initRouter()
	registry.attachRoutes()

	return registry, nil
}

// Start starts the registry.
func (r *Registry) Start() error {
	r.Config.Logger.WithFields(log.Fields{
		"config": *r.Config,
	}).Info("starting registry")
	if r.Config.ShouldUseTLS() {
		return http.ListenAndServeTLS(
			r.Config.Addr,
			r.Config.CertFile,
			r.Config.KeyFile,
			r.Router,
		)
	}
	return http.ListenAndServe(r.Config.Addr, r.Router)
}

func (r *Registry) init() error {
	if err := r.initShaCache(); err != nil {
		return err
	}
	if err := r.initUpstream(); err != nil {
		return err
	}
	r.initStorage()
	return nil
}

func (r *Registry) initShaCache() error {
	shaCache, err := storage.NewShaCache(r.Config.ShaCacheSize)
	r.ShaCache = shaCache
	return err
}

func (r *Registry) initUpstream() error {
	upstream, err := NewUpstream(r.Config.UpstreamURL)
	r.Upstream = upstream
	return err
}

func (r *Registry) initStorage() {
	r.Storage = storage.New(r.Config.StorageDir)
}

func (r *Registry) initRouter() {
	r.Router = httprouter.New()
}

func (r *Registry) attachRoutes() {
	r.Router.GET("/", r.wrapErrHandle(r.HandleRoot))

	r.Router.GET("/:name", r.wrapErrHandle(r.wrapRepoHandle(r.HandlePackageRoot)))
	r.Router.GET("/:name/-/:version", r.wrapErrHandle(r.wrapRepoHandle(r.HandlePkgDownload)))

	r.Router.GET("/:name/ping", r.wrapErrHandle(r.HandlePing))
	r.Router.GET("/:name/stats", r.wrapErrHandle(r.HandleStats))
	r.Router.GET("/:name/upstreams", r.wrapErrHandle(r.HandleUpstreams))
	r.Router.GET("/:name/ui", r.wrapErrHandle(r.HandleUI))
}

// errHandle is a custom HTTP handle that can optionally return an error.
type errHandle func(w http.ResponseWriter, req *http.Request,
	ps httprouter.Params) error

func (r *Registry) wrapErrHandle(handler errHandle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		err := handler(w, req, ps)
		if err == nil {
			return
		}
		contextLog := r.Config.Logger.WithFields(util.GetRequestFields(req))
		util.LogErr(contextLog, err, "handler failed")
		res := &util.ErrorResponse{"internal server error", "unexpected internal error"}
		if err := util.RespondJSON(w, http.StatusInternalServerError, res); err != nil {
			util.LogErr(contextLog, err, "failed to write response")
		}
	}
}

type repoHandle func(repo *git.Repository, w http.ResponseWriter,
	req *http.Request, ps httprouter.Params) error

func (r *Registry) wrapRepoHandle(handle repoHandle) errHandle {
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

		return r.Upstream.HandleReq(w, req, ps)
	}
}
