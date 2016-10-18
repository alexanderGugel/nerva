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
	"github.com/libgit2/git2go"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func createStorage(dir string, t *testing.T) *Storage {
	storage, err := New(dir)
	if err != nil {
		t.Errorf("New(%v) failed: %v", dir, err)
	}
	return storage
}

func TestNew(t *testing.T) {
	dir := "./some/directory"
	storage := createStorage(dir, t)
	if storage.Dir != dir {
		t.Errorf("New(%v).Dir = %v; want %v", dir, storage.Dir, dir)
	}
}

func TestLs(t *testing.T) {
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	dirA := filepath.Join(dir, "a")
	dirB := filepath.Join(dir, "b")

	createTestRepo(dirA, t)
	createTestRepo(dirB, t)

	storage := createStorage(dir, t)

	got, err := storage.Ls()
	want := []string{"a", "b"}
	if err != nil {
		t.Errorf("storage.Ls() failed: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("storage.Ls() = %v; want %v", got, want)
	}
}

func TestLsFailedReadDir(t *testing.T) {
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	nonExistingDir := filepath.Join(dir, "non_existing")
	storage := createStorage(nonExistingDir, t)
	if err := os.RemoveAll(nonExistingDir); err != nil {
		t.Errorf("os.RemoveAll(%v) failed: %v", nonExistingDir, err)
	}

	names, err := storage.Ls()
	if err == nil {
		t.Errorf("storage.Ls() did not fail")
	}
	if names != nil {
		t.Errorf("storage.Ls() = %v; want %v", names, nil)
	}
}

func TestLsFiles(t *testing.T) {
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	dirA := filepath.Join(dir, "a")
	createTestRepo(dirA, t)

	err := ioutil.WriteFile(filepath.Join(dir, "file"), nil, 0644)
	if err != nil {
		t.Fatalf("failed to write %v: %v", dir, err)
	}

	storage := createStorage(dir, t)

	got, err := storage.Ls()
	want := []string{"a"}
	if err != nil {
		t.Errorf("storage.Ls() failed: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("storage.Ls() = %v; want %v", got, want)
	}
}

func TestGetRepo(t *testing.T) {
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	name := "a"
	dirA := filepath.Join(dir, name)
	createTestRepo(dirA, t)

	storage := createStorage(dir, t)
	repo, err := storage.GetRepo(name)
	if err != nil {
		t.Errorf("storage.GetRepo(%v) failed: %v", name, err)
	}
	if repo == nil {
		t.Errorf("storage.GetRepo(%v) = %v; want not nil", name, repo)
	}
}

func TestPeelTreeFailedLookup(t *testing.T) {
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	dirA := filepath.Join(dir, "a")
	repo := createTestRepo(dirA, t)

	var zeroID git.Oid
	if _, err := PeelTree(repo, &zeroID); err == nil {
		t.Errorf("PeelTree(repo, %v) did not fail", zeroID)
	}
}

func createTempDir(t *testing.T) string {
	dir := "storage_test"
	path, err := ioutil.TempDir("", dir)
	if err != nil {
		t.Fatalf("ioutil.TempDir(%v, %v) failed: %v", "", dir, err)
	}
	return path
}

func createTestRepo(path string, t *testing.T) *git.Repository {
	repo, err := git.InitRepository(path, false)
	if err != nil {
		t.Fatalf("git.InitRepository(%v, %t) failed: %v", path, false, err)
	}

	filename := path + "/" + "README"
	data := []byte("foo\n")
	perm := os.FileMode(0644)
	if err = ioutil.WriteFile(filename, data, perm); err != nil {
		t.Fatalf("ioutil.WriteFile(%v, %v, %v) failed: %v", filename, data, perm, err)
	}

	return repo
}
