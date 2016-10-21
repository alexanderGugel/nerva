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
	"github.com/bmizerany/pat"
	"github.com/libgit2/git2go"
	"net/http"
)

// Registry represents an Common JS registry server. A Registry does exposes a
// router, which can be bound to an arbitrary socket.
type Registry struct {
	config   *Config
	mux      *pat.PatternServeMux
	storage  *storage.Storage
	upstream *Upstream
	shaCache *storage.ShaCache
}

// New create a new CommonJS registry.
func New(config *Config) (*Registry, error) {
	if config == nil {
		config = DefaultConfig()
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	registry := &Registry{config: config}
	if err := registry.init(); err != nil {
		return nil, err
	}

	return registry, nil
}

func (r *Registry) init() error {
	initFns := []func() error{
		r.initShaCache,
		r.initUpstream,
		r.initStorage,
		r.initRouter,
	}
	for _, f := range initFns {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

// Start starts the registry.
func (r *Registry) Start() error {
	r.config.Logger.WithFields(log.Fields{
		"config": *r.config,
	}).Info("starting registry")
	server := &http.Server{
		Addr:    r.config.Addr,
		Handler: r.mux,
	}
	if r.config.shouldUseTLS() {
		server.ListenAndServeTLS(
			r.config.CertFile,
			r.config.KeyFile,
		)
	}
	return server.ListenAndServe()
}

func (r *Registry) initShaCache() error {
	shaCache, err := storage.NewShaCache(r.config.ShaCacheSize)
	r.shaCache = shaCache
	return err
}

func (r *Registry) initUpstream() error {
	upstream, err := NewUpstream(r.config.UpstreamURL)
	r.upstream = upstream
	return err
}

func (r *Registry) initStorage() error {
	storage, err := storage.New(r.config.StorageDir)
	r.storage = storage
	return err
}

func (r *Registry) initRouter() error {
	r.mux = pat.New()

	r.mux.Get("/", makeRootEndpoint(r))

	r.mux.Get("/-/ping", makePingEndpoint(r))
	r.mux.Get("/-/ui", makeUIEndpoint(r))
	r.mux.Get("/-/stats", makeStatsEndpoint(r))
	r.mux.Get("/-/upstreams", makeUpstreamsEndpoint(r))

	r.mux.Get("/:name", makePkgRootEndpoint(r))
	r.mux.Get("/:name/-/:version.tgz", makePkgDownloadEndpoint(r))
	r.mux.Get("/:name/stats", makePkgStatsEndpoint(r))

	return nil
}

func makeRootEndpoint(r *Registry) http.HandlerFunc {
	return wrapErrHandle(r.HandleRoot, r.config.Logger)
}

func makePingEndpoint(r *Registry) http.HandlerFunc {
	return wrapErrHandle(r.HandlePing, r.config.Logger)
}

func makeUIEndpoint(r *Registry) http.HandlerFunc {
	return wrapErrHandle(r.HandleUI, r.config.Logger)
}

func makeStatsEndpoint(r *Registry) http.HandlerFunc {
	return wrapErrHandle(HandleMemStats, r.config.Logger)
}

func makeUpstreamsEndpoint(r *Registry) http.HandlerFunc {
	return wrapErrHandle(r.HandleUpstreams, r.config.Logger)
}

func makePkgRootEndpoint(r *Registry) http.HandlerFunc {
	return wrapErrHandle(
		wrapUpstreamHandle(
			wrapRepoHandle(r.HandlePackageRoot, r.storage),
			r.upstream,
		),
		r.config.Logger,
	)
}

func makePkgDownloadEndpoint(r *Registry) http.HandlerFunc {
	return wrapErrHandle(
		wrapUpstreamHandle(
			wrapRepoHandle(r.HandlePkgDownload, r.storage),
			r.upstream,
		),
		r.config.Logger,
	)
}

func makePkgStatsEndpoint(r *Registry) http.HandlerFunc {
	return wrapErrHandle(
		wrapRepoHandle(HandlePkgStats, r.storage),
		r.config.Logger,
	)
}

// errHandle is a custom HTTP handle that can optionally return an error.
type errHandle func(http.ResponseWriter, *http.Request) error

func wrapErrHandle(handler errHandle, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := handler(w, req)
		if err == nil {
			return
		}
		contextLog := logger.WithFields(util.GetRequestFields(req))
		util.LogErr(contextLog, err, "handler failed")
		code := http.StatusInternalServerError
		res := &util.ErrorResponse{
			http.StatusText(code),
			"unexpected internal error",
		}
		if err := util.RespondJSON(w, code, res); err != nil {
			util.LogErr(contextLog, err, "failed to write response")
		}
	}
}

type repoHandle func(*git.Repository, http.ResponseWriter, *http.Request) error

func wrapRepoHandle(handle repoHandle, storage *storage.Storage) errHandle {
	return func(w http.ResponseWriter, req *http.Request) error {
		name := req.URL.Query().Get(":name")
		repo, err := storage.GetRepo(name)
		if err != nil {
			return err
		}
		return handle(repo, w, req)
	}
}

func wrapUpstreamHandle(handle errHandle, upstream *Upstream) errHandle {
	return func(w http.ResponseWriter, req *http.Request) error {
		err := handle(w, req)
		if err == nil {
			return nil
		}
		if gitErr, ok := err.(*git.GitError); !ok ||
			gitErr.Class != git.ErrClassOs {
			return err
		}
		return upstream.HandleReq(w, req)
	}
}
