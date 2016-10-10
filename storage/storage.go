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

// Package storage allows access to git directories stored in a (usually
// registry-specific) storage directory.
package storage

import (
	"github.com/libgit2/git2go"
	"io/ioutil"
	"path"
)

// Storage manages a directory of Git repositories.
type Storage struct {
	Dir string
}

// New creates a new storage bound to the specified directory.
func New(dir string) *Storage {
	return &Storage{dir}
}

// GetRepo opens the repository in the sub-directory "name".
func (s *Storage) GetRepo(name string) (*git.Repository, error) {
	abs := path.Join(s.Dir, name)
	return git.OpenRepository(abs)
}

// Ls lists all available repository names.
func (s *Storage) Ls() ([]string, error) {
	files, err := ioutil.ReadDir(s.Dir)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, file := range files {
		name := file.Name()
		if file.IsDir() {
			names = append(names, name)
		}
	}

	return names, nil
}

// PeelTree recursively traverses the passed in Git object until a Git tree
// object is found.
func PeelTree(repo *git.Repository, id *git.Oid) (*git.Tree, error) {
	object, err := repo.Lookup(id)
	if err != nil || object == nil {
		return nil, err
	}

	treeObject, err := object.Peel(git.ObjectTree)
	if err != nil || treeObject == nil {
		return nil, err
	}

	tree, err := treeObject.AsTree()
	if err != nil || tree == nil {
		return nil, err
	}

	return tree, nil
}
