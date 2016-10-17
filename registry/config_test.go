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
	log "github.com/Sirupsen/logrus"
	"testing"
)

func TestDefaultConfigNil(t *testing.T) {
	config := DefaultConfig()
	if config == nil {
		t.Errorf("DefaultConfig() = %q, want config struct", config)
	}
}

var configValidateTests = []struct {
	config  Config
	isValid bool
}{
	{
		config:  Config{},
		isValid: false,
	},
	{
		config: Config{
			Addr: ":8200",
		},
		isValid: false,
	},
	{
		config: Config{
			Addr:      ":8200",
			FrontAddr: "http://127.0.0.1:8200",
		},
		isValid: false,
	},
	{
		config: Config{
			Addr:      ":8200",
			FrontAddr: "http://127.0.0.1:8200",
			Logger:    log.StandardLogger(),
		},
		isValid: true,
	},
	{
		config: Config{
			Addr:      ":8200",
			FrontAddr: "http://127.0.0.1:8200",
			Logger:    log.StandardLogger(),
			CertFile:  "./certfile",
		},
		isValid: false,
	},
	{
		config: Config{
			Addr:      ":8200",
			FrontAddr: "http://127.0.0.1:8200",
			Logger:    log.StandardLogger(),
			KeyFile:   "./keyfile",
		},
		isValid: false,
	},
	{
		config: Config{
			Addr:      ":8200",
			FrontAddr: "http://127.0.0.1:8200",
			Logger:    log.StandardLogger(),
			CertFile:  "./certFile",
			KeyFile:   "./keyfile",
		},
		isValid: true,
	},
}

func TestConfigValidate(t *testing.T) {
	for _, tt := range configValidateTests {
		err := tt.config.Validate()
		isValid := err == nil
		if isValid != tt.isValid {
			t.Errorf("%d.Validate() = %t, want %t", tt.config, isValid, tt.isValid)
		}
	}
}
