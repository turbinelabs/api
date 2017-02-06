package api

import (
	"fmt"
	"sort"
)

// Instances is a slice of Instance
type Instances []Instance

func (i Instances) equality(o Instances, checkfn func(Instance, Instance) bool) bool {
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
		} else if !checkfn(iInst, oInst) {
			return false
		}
	}

	return true
}

var equalsFnAdapter func(Instance, Instance) bool = func(i1, i2 Instance) bool {
	return i1.Equals(i2)
}

func (i Instances) Equals(o Instances) bool {
	return i.equality(o, equalsFnAdapter)
}

var equivFnAdapter func(Instance, Instance) bool = func(i1, i2 Instance) bool {
	return i1.Equivalent(i2)
}

func (i Instances) Equivalent(o Instances) bool {
	return i.equality(o, equivFnAdapter)
}

// Checks a collection of instances to ensure all are valid
func (i Instances) IsValid(precreation bool) *ValidationError {
	errs := &ValidationError{}

	for _, e := range i {
		errs.MergePrefixed(e.IsValid(precreation), "")
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

// Checks for approximate object equivalence. This requires Instance host and port are
// the same and its metadata should be equivalent as well.
func (i Instance) Equivalent(o Instance) bool {
	return i.hostPortCheck(o) && i.Metadata.Equivalent(o.Metadata)
}

// checks for host and port data as both are required for an instance to be
// well defined
func (i Instance) IsValid(precreation bool) *ValidationError {
	ecase := func(f, m string) ErrorCase {
		return ErrorCase{fmt.Sprintf("instance[%s].%s", i.Key(), f), m}
	}

	errs := &ValidationError{}

	validHost := i.Host != ""
	validPort := i.Port != 0

	if !validHost {
		errs.AddNew(ecase("host", "must not be empty"))
	}

	if !validPort {
		errs.AddNew(ecase("port", "must be non-zero"))
	}

	return errs.OrNil()
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
