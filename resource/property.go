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

package resource

// Property type represents a resource property, which can be
// evaluated and set if needed.
type Property interface {
	// Name returns the property name
	Name() string

	// Set sets the property to it's desired state.
	Set() error

	// IsSynced returns a boolean indicating whether the
	// resource property is in sync or not.
	IsSynced() (bool, error)
}

// ResourceProperty type implements the Property interface.
type ResourceProperty struct {
	// PropertySetFunc is the type of the function that is called when
	// setting a resource property to it's desired state.
	PropertySetFunc func() error

	// PropertyIsSyncedFunc is the type of the function that is called when
	// determining whether a resource property is in the desired state.
	PropertyIsSyncedFunc func() (bool, error)

	// PropertyName is the name of the property.
	PropertyName string
}

// Set sets the property to it's desired state.
func (rp *ResourceProperty) Set() error {
	return rp.PropertySetFunc()
}

// IsSynced returns a boolean indicating whether the
// resource property is in the desired state.
func (rp *ResourceProperty) IsSynced() (bool, error) {
	return rp.PropertyIsSyncedFunc()
}

// Name returns the property name.
func (rp *ResourceProperty) Name() string {
	return rp.PropertyName
}
