package stats

type ZeroFill string

const (
	// None is the default mode and doesn't fill any values in. This is the
	// default behavior.
	None ZeroFill = "none"

	// Partial fills in only series that are partially complete
	Partial ZeroFill = "partial"

	// Full fills in all series even if there was no data initially. If this is
	// the case EmptySeries will be set on the TimeSeries.
	Full ZeroFill = "full"
)

func (zf ZeroFill) IsNone() bool {
	return zf == None || !(zf.IsPartial() || zf.IsFull())
}

func (zf ZeroFill) IsPartial() bool {
	return zf == Partial
}

func (zf ZeroFill) IsFull() bool {
	return zf == Full
}
