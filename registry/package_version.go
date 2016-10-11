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
    "encoding/json"
    "github.com/alexanderGugel/nerva/storage"
    "github.com/libgit2/git2go"
)

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
