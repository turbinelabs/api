package api

// A type commonly embeded in other domain objects to ensure modifications
// are being on an underlying object in the expected state.
type Checksum struct {
	Checksum string `json:"checksum"` // may be overwritten
}

func (c *Checksum) IsNil() bool {
	return c.Equals(Checksum{})
}

// An empty checksum is equivalent to an unset checksum.
func (c *Checksum) IsEmpty() bool {
	return len(c.Checksum) == 0
}

func (c Checksum) Equals(o Checksum) bool {
	return c.Checksum == o.Checksum
}
