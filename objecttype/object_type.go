package objecttype

import (
	"fmt"
	"strings"
)

var all []ObjectType

func init() {
	for i := 1; ; i++ {
		ot, err := FromID(i)
		if err != nil {
			break
		}
		all = append(all, ot)
	}
}

// All returns a slice of all object types.
func All() []ObjectType {
	r := make([]ObjectType, 0, len(all))
	copy(r, all)
	return r
}

// AllNames returns a slice containing the names of all object types.
func AllNames() []string {
	r := make([]string, 0, len(all))
	for _, ot := range all {
		r = append(r, ot.Name)
	}
	return r
}

// ObjectType is representation of an object that can have the changes made
// made to it tracked in a persistant changelog.
type ObjectType struct {
	Name string `json:"object_type"`
	id   int64
}

// We define an ObjectType enum below. New values may be added
// but the existing value SHOULD NOT be changed or removed.
var (
	Org         = ObjectType{"org", 1}
	User        = ObjectType{"user", 2}
	Zone        = ObjectType{"zone", 3}
	Proxy       = ObjectType{"proxy", 4}
	Domain      = ObjectType{"domain", 5}
	Route       = ObjectType{"route", 6}
	Cluster     = ObjectType{"cluster", 7}
	SharedRules = ObjectType{"shared_rules", 8}
)

func (ot ObjectType) ID() int64 {
	return ot.id
}

var UnrecognizedObjectTypeError = fmt.Errorf("unrecognized object type")

func FromName(s string) (ObjectType, error) {
	s = strings.ToLower(s)

	switch s {
	case Org.Name:
		return Org, nil
	case User.Name:
		return User, nil
	case Zone.Name:
		return Zone, nil
	case Proxy.Name:
		return Proxy, nil
	case Domain.Name:
		return Domain, nil
	case Route.Name:
		return Route, nil
	case Cluster.Name:
		return Cluster, nil
	case SharedRules.Name:
		return SharedRules, nil
	}

	return ObjectType{}, UnrecognizedObjectTypeError
}

func FromID(i int) (ObjectType, error) {
	switch int64(i) {
	case Org.id:
		return Org, nil
	case User.id:
		return User, nil
	case Zone.id:
		return Zone, nil
	case Proxy.id:
		return Proxy, nil
	case Domain.id:
		return Domain, nil
	case Route.id:
		return Route, nil
	case Cluster.id:
		return Cluster, nil
	case SharedRules.id:
		return SharedRules, nil
	}

	return ObjectType{}, UnrecognizedObjectTypeError
}
