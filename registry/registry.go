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
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/alexanderGugel/nerva/storage"
	"github.com/alexanderGugel/nerva/util"
	"github.com/julienschmidt/httprouter"
	"github.com/libgit2/git2go"
	"net/http"
	"path"
	"regexp"
	"runtime"
)

// Registry represents an Common JS registry server. A Registry does exposes a
// router, which can be bound to an arbitrary socket.
type Registry struct {
	Router   *httprouter.Router
	Storage  *storage.Storage
	Upstream *Upstream
	shaCache *ShaCache
}

// PackageVersion represents a specific version of a package, typically its
// package.json file.
// The Package Version Object is almost identical to the Package Descriptor
// object described in the CommonJS Packages specification. For the purposes of
// the package registry, the following fields are required. Note that some of
// these do not exist in the Packages specification.
// name: The package name. This MUST match the {package name} portion of the
// URL.
// version: The package version. This MUST match the {package version} portion
// of the URL.
// dist: An object hash with urls of where the package archive can be found. The
// key is the type of archive. At the moment the following archive types are
// supported, but more may be added in the future:
// tarball: A url to a gzipped tar archive containing a single folder with the
// package contents (including the package.json file in the root of said
// folder).
// See http://wiki.commonjs.org/wiki/Packages/Registry#Package_Version_Object
type PackageVersion map[string]interface{}

// PackageRootVersions represents the versions map in a package root document
// that will be served when clients request a specific package, but no version.
type PackageRootVersions map[string]*PackageVersion

// PackageRoot represents a CommonJS package root document containing all
// available versions of a given package available in its local repository.
// The root object that describes all versions of a package MUST be a JSON
// object with the following fields:
// name: The name of the package. When both are decoded, this MUST match the
// “package name” portion of the URL. That is, packages with irregular
// characters in their names would be URL-Encoded in the request, and
// JSON-encoded in the data. So, a request to
// /%C3%A7%C2%A5%C3%A5%C3%B1%C3%AE%E2%88%82%C3%A9 would show a package root
// object with “\u00e7\u00a5\u00e5\u00f1\u00ee\u2202\u00e9” as the name, and
// would refer to the “ç¥åñî∂é” project.
// versions: An object hash of version identifiers to valid “package version
// url” responses: either URL strings or package descriptor objects.
// See http://wiki.commonjs.org/wiki/Packages/Registry#Package_Root_Object
type PackageRoot struct {
	Name     string               `json:"name"`
	DistTags *PackageDistTags     `json:"dist-tags"`
	Versions *PackageRootVersions `json:"versions"`
}

// PackageDistTags represents the dist-tags of a package root object. It maps
// Common JS tags, such latest, to specific versions.
type PackageDistTags map[string]string

// PackageDist describes how a package can be downloaded.
type PackageDist struct {
	Tarball string `json:"tarball"`
	Shasum  string `json:"shasum"`
}

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

// Ping represents the response to a ping request. Ping is simply an empty
// object used by npm ping.
type Ping struct{}

// NewPing creates a new ping object.
func NewPing() *Ping {
	return &Ping{}
}

// NewRoot creates a new CommonJS registry root document from a given
// storage directory by reading in the repositories that are available in the
// storage dir.
func NewRoot(storage *storage.Storage, url string) (*Root, error) {
	Root := Root{}
	names, _ := storage.Ls()
	for _, name := range names {
		Root[name] = path.Join(url, name)
	}
	return &Root, nil
}

var versionTagRef = regexp.MustCompile("^refs\\/tags\\/v(.*)$")

// NewPackageRoot creates a new CommonJS package root document.
func NewPackageRoot(name string, url string, repo *git.Repository, shaCache *ShaCache) (*PackageRoot, error) {
	versions := PackageRootVersions{}
	contextLog := log.WithFields(log.Fields{"name": name})

	latest := ""

	repo.Tags.Foreach(func(tagRef string, id *git.Oid) error {
		contextLog := contextLog.WithFields(log.Fields{"tagRef": tagRef})

		if !versionTagRef.MatchString(tagRef) {
			contextLog.Debug("skipping non-version tag")
			return nil
		}

		packageVersion, err := NewPackageVersion(repo, id)
		contextLog = contextLog.WithFields(log.Fields{"packageVersion": packageVersion})
		if err != nil || packageVersion == nil {
			util.LogErr(contextLog, err, "failed to generate package version")
			return nil
		} else if version, ok := (*packageVersion)["version"].(string); ok {
			contextLog = contextLog.WithFields(log.Fields{"version": version})
			if versions[version] != nil {
				contextLog.Warn("duplicate version")
			}
			tarball := url + "/" + name + "/-/" + id.String()

			shasum, ok := shaCache.Get(*id)
			if !ok {
				hasher := sha1.New()
				d, err := storage.NewDownload(repo, id)
				if err != nil || d == nil {
					util.LogErr(contextLog, err, "failed to create download")
					return nil
				}
				if err := d.Start(hasher); err != nil {
					util.LogErr(contextLog, err, "failed to start download")
				}
				shasum = hex.EncodeToString(hasher.Sum(nil))
			}

			shaCache.Add(*id, shasum)
			(*packageVersion)["dist"] = &PackageDist{tarball, shasum}
			versions[version] = packageVersion
			latest = version
		}
		return nil
	})

	distTags := PackageDistTags{}
	if latest != "" {
		distTags["latest"] = latest
	}
	packageRoot := &PackageRoot{name, &distTags, &versions}
	return packageRoot, nil
}

// NewPackageVersion creates a package root object (package.json) from a given
// Git Object id.
func NewPackageVersion(repo *git.Repository, id *git.Oid) (*PackageVersion, error) {
	tree, err := storage.PeelTree(repo, id)
	if err != nil || tree == nil {
		return nil, err
	}

	entry := tree.EntryByName("package.json")
	blob, err := repo.LookupBlob(entry.Id)
	if err != nil || blob == nil {
		return nil, err
	}

	contents := blob.Contents()
	packageVersion := &PackageVersion{}
	if err = json.Unmarshal(contents, packageVersion); err != nil {
		return nil, err
	}

	return packageVersion, nil
}

// PackageStats contains information about the underlying git repo of a package.
type PackageStats struct {
	Remotes []*PackageRemote `json:"remotes"`
}

// PackageRemote is the equivalent to `git remote -v`.
type PackageRemote struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// NewPackageStats create a package stats object, which contains information
// about the underlying git repository.
func NewPackageStats(repo *git.Repository) (*PackageStats, error) {
	names, err := repo.Remotes.List()
	if err != nil {
		return nil, err
	}
	remotes := []*PackageRemote{}
	for _, name := range names {
		remote, err := repo.Remotes.Lookup(name)
		if err != nil {
			return nil, err
		}
		url := remote.Url()
		remotes = append(remotes, &PackageRemote{name, url})
	}
	stats := &PackageStats{remotes}
	return stats, nil
}

// New create a new CommonJS registry.
func New(storageDir string, upstreamURL string, shaCacheSize int) (*Registry, error) {
	router := httprouter.New()
	storage := storage.New(storageDir)
	upstream, err := NewUpstream(upstreamURL)
	if err != nil {
		return nil, err
	}

	shaCache, err := NewShaCache(shaCacheSize)
	if err != nil {
		return nil, err
	}

	registry := &Registry{router, storage, upstream, shaCache}

	router.GET("/", util.ErrHandler(registry.HandleRoot))

	pkgRoot := util.ValidatePropHandler("name", registry.repoHandler(registry.HandlePackageRoot))
	router.GET("/:name", util.ErrHandler(pkgRoot))

	download := util.ValidatePropHandler("name", registry.repoHandler(registry.HandlePackageDownload))
	router.GET("/:name/-/:version", util.ErrHandler(download))

	router.GET("/:name/ping", util.ErrHandler(registry.HandlePing))

	router.GET("/:name/stats", util.ErrHandler(registry.HandleStats))

	return registry, nil
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
func (r *Registry) HandleRoot(w http.ResponseWriter, req *http.Request,
	_ httprouter.Params) error {
	res, err := NewRoot(r.Storage, req.Host)
	if err != nil {
		return err
	}
	return util.RespondJSON(w, 200, res)
}

// HandlePing responds with an empty JSON object. npm's ping command hits this
// endpoint.
func (r *Registry) HandlePing(w http.ResponseWriter, req *http.Request,
	_ httprouter.Params) error {
	res := NewPing()
	return util.RespondJSON(w, 200, res)
}

// HandlePackageStats retrieves the current memory stats.
func (r *Registry) HandlePackageStats(repo *git.Repository,
	w http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
	res, err := NewPackageStats(repo)
	if err != nil {
		return err
	}
	return util.RespondJSON(w, 200, res)
}

// NewMemStats aggregates and returns memory stats.
func NewMemStats() *runtime.MemStats {
	var m = new(runtime.MemStats)
	runtime.ReadMemStats(m)
	return m
}

// HandleMemStats retrieves the current memory stats.
func (r *Registry) HandleMemStats(w http.ResponseWriter, req *http.Request,
	ps httprouter.Params) error {
	res := NewMemStats()
	return util.RespondJSON(w, 200, res)
}

// HandleStats retrieves the current memory stats.
func (r *Registry) HandleStats(w http.ResponseWriter, req *http.Request,
	ps httprouter.Params) error {
	name := ps.ByName("name")
	if name == "-" {
		return r.HandleMemStats(w, req, ps)
	}
	return util.ValidatePropHandler("name", r.repoHandler(r.HandlePackageStats))(w, req, ps)
}

// HandlePackageDownload handles package downloads.
func (r *Registry) HandlePackageDownload(repo *git.Repository,
	w http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
	version := ps.ByName("version")

	id, err := git.NewOid(version)
	if err != nil || id == nil {
		res := &util.ErrorResponse{"bad request", "version is not a valid git object id"}
		return util.RespondJSON(w, http.StatusBadRequest, res)
	}

	d, err := storage.NewDownload(repo, id)
	if err != nil || d == nil {
		gitErr, ok := err.(*git.GitError)

		if !ok || gitErr.Class != git.ErrClassOdb {
			return err
		}
		res := &util.ErrorResponse{"not found", "package not found"}
		return util.RespondJSON(w, http.StatusNotFound, res)
	}
	return d.Start(w)
}

// HandlePackageRoot handles requests to the package root URL.
// The package root url is the base URL where a client can get top-level
// information about a package and all of the versions known to the registry.
// A valid “package root url” response MUST be returned when the client requests
// {registry root url}/{package name}.
// See http://wiki.commonjs.org/wiki/Packages/Registry#package_root_url
func (r *Registry) HandlePackageRoot(repo *git.Repository,
	w http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
	name := ps.ByName("name")
	res, err := NewPackageRoot(name, req.Host, repo, r.shaCache)
	if err != nil {
		return err
	}
	return util.RespondJSON(w, 200, res)
}

type repoHandle func(repo *git.Repository, w http.ResponseWriter,
	req *http.Request, ps httprouter.Params) error

func (r *Registry) repoHandler(handle repoHandle) util.ErrHandle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
		name := ps.ByName("name")

		repo, err := r.Storage.GetRepo(name)

		if err != nil {
			gitErr, ok := err.(*git.GitError)
			if !ok || gitErr.Class != git.ErrClassOs {
				return err
			}

			if r.Upstream == nil || r.Upstream.URL == nil {
				res := &util.ErrorResponse{"not found", "package not found"}
				return util.RespondJSON(w, http.StatusNotFound, res)
			}

			r.Upstream.RedirectPackageRoot(name, w, req)
			return nil
		}

		return handle(repo, w, req, ps)
	}
}
