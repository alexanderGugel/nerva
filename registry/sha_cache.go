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
  "github.com/hashicorp/golang-lru"
  "github.com/libgit2/git2go"
)

// ShaCache serves as an adapter for an immutable LRU cache. Git object ids
// are cryptographically unique, therefore there is no need to "manually" remove
// items from the underlying LRU cache.
type ShaCache struct {
  lru *lru.Cache
}

// NewShaCache creates a new LRU cache used for mapping Git object ids to
// respective SHA1 sums.
func NewShaCache(size int) (*ShaCache, error) {
  lru, err := lru.New(size)
  if err != nil {
    return nil, err
  }
  return &ShaCache{lru}, nil
}

// Add populates the cache with the given Git object id.
func (c *ShaCache) Add(id git.Oid, shasum string) bool {
  return c.lru.Add(id, shasum)
}

// Get retrieves the corresponding SHA sum for the supplied Git object id.
func (c *ShaCache) Get(id git.Oid) (string, bool) {
  shasum, ok := c.lru.Get(id)
  if !ok {
	  return "", false
  }
  return shasum.(string), ok
}
