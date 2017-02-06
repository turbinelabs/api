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

func ChangeTypeFromName(s string) (ChangeType, error) {
	switch s {
	case Addition.Name:
		return Addition, nil
	case Removal.Name:
		return Removal, nil
	}

	return ChangeType{}, UnrecognizedChangeTypeError
}

func ChangeTypeFromID(i int) (ChangeType, error) {
	switch int64(i) {
	case Addition.id:
		return Addition, nil
	case Removal.id:
		return Removal, nil
	}

	return ChangeType{}, UnrecognizedChangeTypeError
}
