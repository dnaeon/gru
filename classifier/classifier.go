// Copyright (c) 2015-2017 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//  1. Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer
//     in this position and unchanged.
//  2. Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in the
//     documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package classifier

import "errors"

// ErrClassifierNotFound is returned if the requested classifier was
// not found in the classifier registry
var ErrClassifierNotFound = errors.New("Classifier key not found")

// Classifier type contains a key/value pair repsenting a classifier
type Classifier struct {
	// Classifier key
	Key string `json:"key"`

	// Classifier value
	Value string `json:"value"`
}

// Type of classifier providers
// A classifier provider is what does the
// actual evaluation of a classifier
type provider func() (string, error)

// Registry provides a global registry for all classifiers
var Registry = make(map[string]provider)

// Register registers a new classifier to the registry
func Register(key string, p provider) error {
	Registry[key] = p

	return nil
}

// Get retrieves a classifier from the registry by looking up its key
func Get(key string) (*Classifier, error) {
	c := new(Classifier)

	p, ok := Registry[key]
	if ok {
		// Evaluate the classifier provider
		value, err := p()
		c := &Classifier{
			Key:   key,
			Value: value,
		}
		return c, err
	}

	return c, ErrClassifierNotFound
}
