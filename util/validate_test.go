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

package util

import "testing"

var nameTests = []struct {
	name    string
	isValid bool
}{
	{"tape", true},
	// 1. MUST NOT start with “-“
	{"-tape", false},
	// 2. MUST NOT contain any “/” characters
	{"t/ape", false},
	// 3. MUST NOT be “.” or “..”
	{"..", false},
	{".", false},
	{"-", false},
}

func TestIsValid(t *testing.T) {
	for _, tt := range nameTests {
		isValid := IsValid(tt.name)
		if isValid != tt.isValid {
			t.Errorf("IsValid(%q) = %t, want %t", tt.name, isValid, tt.isValid)
		}
	}
}
