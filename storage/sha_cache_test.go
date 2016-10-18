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
	"testing"
)

func TestNewShaCacheInvalidSize(t *testing.T) {
	for i := -1; i < 1; i++ {
		if _, err := NewShaCache(i); err == nil {
			t.Errorf("NewShaCache(%v) failed: %v", i, err)
		}
	}
}

func TestNewShaCachePositiveSize(t *testing.T) {
	for i := 1; i < 100; i++ {
		if _, err := NewShaCache(i); err != nil {
			t.Errorf("NewShaCache(%v) failed: %v", i, err)
		}
	}
}

func TestShaCacheAdd(t *testing.T) {
	c, err := NewShaCache(1)
	if err != nil {
		t.Errorf("NewShaCache(%v) failed: %v", 1, err)
	}
	id0 := *git.NewOidFromBytes([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	shasum0 := "shasum0"
	if ok := c.Add(id0, shasum0); ok {
		t.Errorf("c.Add(%v, %v) = %v; want %v", id0, shasum0, ok, false)
	}
	id1 := *git.NewOidFromBytes([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	shasum1 := "shasum1"
	if ok := c.Add(id1, shasum1); !ok {
		t.Errorf("c.Add(%v, %v) = %v; want %v", id1, shasum1, ok, true)
	}
}

func TestShaCacheGet(t *testing.T) {
	c, err := NewShaCache(1)
	if err != nil {
		t.Errorf("NewShaCache(%v) unexpected err: %v", 1, err)
	}
	id := *git.NewOidFromBytes([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	shasum := "shasum"
	if ok := c.Add(id, shasum); ok {
		t.Errorf("c.Add(%v, %v) = %v; want %v", id, shasum, ok, false)
	}
	result, ok := c.Get(id)
	if !ok || result != shasum {
		t.Errorf("c.Get(%v) = %v, %v; want %v, %v", id, result, ok, shasum, true)
	}
}
