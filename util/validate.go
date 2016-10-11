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

// IsValid checks if the given name is valid according to the CommonJS spec.
// Besides the addition of fields to the Package Version Object, this addition
// to the Packages spec imposes the following restrictions on the “name” and
// “version” fields:
// 1. MUST NOT start with “-“
// 2. MUST NOT contain any “/” characters
// 3. MUST NOT be “.” or “..”
// 4. SHOULD contain only URL-safe characters
// See http://wiki.commonjs.org/wiki/Packages/Registry#Changes_to_Packages_Spec
func IsValid(name string) bool {
    if name[0] == '-' || name == "." || name == ".." {
        return false
    }
    for _, c := range name {
        if c == '/' {
            return false
        }
    }
    return true
}
