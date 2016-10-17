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
	"errors"
	log "github.com/Sirupsen/logrus"
)

// Config represents the configuration options of registry.
type Config struct {
	StorageDir   string
	UpstreamURL  string
	ShaCacheSize int
	Addr         string
	CertFile     string
	KeyFile      string
	FrontAddr    string
	Logger       *log.Logger `json:"-"`
}

// DefaultConfig create a default configuration with sane defaults.
func DefaultConfig() *Config {
	return &Config{
		StorageDir:   "./packages",
		UpstreamURL:  "http://registry.npmjs.com",
		ShaCacheSize: 500,
		Addr:         ":8200",
		CertFile:     "",
		KeyFile:      "",
		FrontAddr:    "http://127.0.0.1:8200",
		Logger:       log.StandardLogger(),
	}
}

// shouldUseTLS checks if TLS is (partially) configured.
func (c *Config) shouldUseTLS() bool {
	return c.CertFile != "" || c.KeyFile != ""
}

// Validate checks if the supplied config is valid.
func (c *Config) Validate() error {
	if c.Addr == "" {
		return errors.New("missing Addr")
	}
	if c.shouldUseTLS() && c.CertFile == "" {
		return errors.New("missing CertFile")
	}
	if c.shouldUseTLS() && c.KeyFile == "" {
		return errors.New("missing KeyFile")
	}
	if c.FrontAddr == "" {
		return errors.New("missing FrontAddr")
	}
	if c.Logger == nil {
		return errors.New("missing Logger")
	}
	return nil
}
