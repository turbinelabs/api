/*
Copyright 2017 Turbine Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package changetype

import (
	"fmt"
)

// ChangeType describes the types of changes that can be made to data within
// a struct.
type ChangeType struct {
	Name string `json:"change_type"`
	id   int64
}

// We define a ChangeType enum below. New values may be added but the
// existing value SHOULD NOT be changed or removed.
var (
	Addition = ChangeType{"addition", 1}
	Removal  = ChangeType{"removal", 2}
)

func (ct ChangeType) ID() int64 {
	return ct.id
}

var UnrecognizedChangeTypeError = fmt.Errorf("unrecognized change type")

func FromName(s string) (ChangeType, error) {
	switch s {
	case Addition.Name:
		return Addition, nil
	case Removal.Name:
		return Removal, nil
	}

	return ChangeType{}, UnrecognizedChangeTypeError
}

func FromID(i int) (ChangeType, error) {
	switch int64(i) {
	case Addition.id:
		return Addition, nil
	case Removal.id:
		return Removal, nil
	}

	return ChangeType{}, UnrecognizedChangeTypeError
}
