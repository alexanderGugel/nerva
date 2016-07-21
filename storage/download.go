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

package storage

import (
	"archive/tar"
	log "github.com/Sirupsen/logrus"
	"github.com/alexanderGugel/nerva/util"
	"github.com/libgit2/git2go"
	"io"
	"path"
)

// Download represents an ongoing download.
type Download struct {
	Repo *git.Repository
	Tree *git.Tree
}

// NewDownload creates a new download.
func NewDownload(repo *git.Repository, id *git.Oid) (*Download, error) {
	ref, err := repo.Lookup(id)
	if err != nil || ref == nil {
		return nil, err
	}

	treeObject, err := ref.Peel(git.ObjectTree)
	if err != nil || treeObject == nil {
		return nil, err
	}

	tree, err := treeObject.AsTree()
	if err != nil || tree == nil {
		return nil, err
	}

	return &Download{repo, tree}, nil
}

// Start recursively traverses the internal Git object tree and dynamically
// creates and compresses the corresponding tarball.
func (d *Download) Start(w io.Writer) error {
	tarWriter := tar.NewWriter(w)
	defer tarWriter.Close()

	return d.Tree.Walk(func(dir string, entry *git.TreeEntry) int {
		name := path.Join("package", dir, entry.Name)
		contextLog := log.WithFields(log.Fields{"name": name})

		switch entry.Type {
		case git.ObjectBlob:
			blob, err := d.Repo.LookupBlob(entry.Id)
			if err != nil {
				util.LogErr(contextLog, err, "failed to lookup blob")
				return 1
			}
			hdr := &tar.Header{
				Name: name,
				Mode: int64(entry.Filemode),
				Size: int64(blob.Size()),
			}
			if err := tarWriter.WriteHeader(hdr); err != nil {
				util.LogErr(contextLog, err, "failed to write headers")
				return 1
			}
			if _, err := tarWriter.Write(blob.Contents()); err != nil {
				util.LogErr(contextLog, err, "failed to write contents")
				return 1
			}
			if err := tarWriter.Flush(); err != nil {
				util.LogErr(contextLog, err, "failed to flush writer")
				return 1
			}
		}
		return 0
	})
}
