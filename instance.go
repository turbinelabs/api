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

package api

import (
	"fmt"
	"regexp"
	"sort"
)

const (
	// HostPatternString represents the pattern that an Instance hostanme must
	// match.
	HostPatternString = "^[a-zA-Z0-9_.-]+$"

	// HostPatternMatchFailure is the error message returned when in invalid name
	// is provided.
	HostPatternMatchFailure = "host must match " + HostPatternString
)

var hostPattern = regexp.MustCompile(HostPatternString)

// Instances is a slice of Instance
type Instances []Instance

func (i Instances) Equals(o Instances) bool {
	if len(i) != len(o) {
		return false
	}

	oMap := make(map[string]Instance)

	for _, inst := range o {
		oMap[inst.Key()] = inst
	}

	for _, iInst := range i {
		if oInst, oOK := oMap[iInst.Key()]; !oOK {
			return false
		} else if !iInst.Equals(oInst) {
			return false
		}
	}

	return true
}

// Checks a collection of instances to ensure all are valid
func (i Instances) IsValid() *ValidationError {
	errs := &ValidationError{}

	seen := map[string]bool{}
	for _, e := range i {
		if seen[e.Key()] {
			errs.AddNew(ErrorCase{"instances", fmt.Sprintf("multiple instances of key %v", e.Key())})
		}
		seen[e.Key()] = true
		errs.MergePrefixed(e.IsValid(), fmt.Sprintf("instances[%v]", e.Key()))
	}

	return errs.OrNil()
}

// An Instance is a hostname/port pair with Metadata
type Instance struct {
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	Metadata Metadata `json:"metadata"`
}

func (i Instance) IsNil() bool {
	return i.Equals(Instance{})
}

func (i Instance) Key() string {
	return fmt.Sprintf("%s:%d", i.Host, i.Port)
}

func (i Instance) hostPortCheck(i2 Instance) bool {
	return !(i.Host != i2.Host || i.Port != i2.Port)
}

// Checks for exact object equality. This requires Instance host and port are
// equal as well as its metadata.
func (i Instance) Equals(o Instance) bool {
	return i.hostPortCheck(o) && i.Metadata.Equals(o.Metadata)
}

// checks for host and port data as both are required for an instance to be
// well defined
func (i Instance) IsValid() *ValidationError {
	iscope := fmt.Sprintf("instances[%s]", i.Key())

	ecase := func(f, m string) ErrorCase {
		return ErrorCase{iscope + "." + f, m}
	}

	errs := &ValidationError{}

	if i.Host == "" {
		errs.AddNew(ecase("host", "must not be empty"))
	} else if !hostPattern.MatchString(i.Host) {
		errs.AddNew(ecase("host", HostPatternMatchFailure))
	}

	if i.Port == 0 {
		errs.AddNew(ecase("port", "must be non-zero"))
	}

	errs.MergePrefixed(InstanceMetadataIsValid(i.Metadata), iscope)

	return errs.OrNil()
}

func InstanceMetadataIsValid(md Metadata) *ValidationError {
	return MetadataValid("metadata", md, func(kv Metadatum) *ValidationError {
		if !AllowedIndexPattern.MatchString(kv.Key) {
			return &ValidationError{
				[]ErrorCase{{"key", AllowedIndexPatternMatchFailure}},
			}
		}

		return nil
	})
}

// Sort a Instancds by Host and Port.
// Eg: sort.Sort(InstancesByHostPort(instances))
type InstancesByHostPort Instances

var _ sort.Interface = InstancesByHostPort{}

func (b InstancesByHostPort) Len() int      { return len(b) }
func (b InstancesByHostPort) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b InstancesByHostPort) Less(i, j int) bool {
	return b[i].Host < b[j].Host && b[i].Port < b[j].Port
}
