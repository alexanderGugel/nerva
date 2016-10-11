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

import "github.com/libgit2/git2go"

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
