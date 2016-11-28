package timegranularity

import (
	"encoding/json"
	"fmt"
)

type TimeGranularity int

const (
	Seconds TimeGranularity = iota
	Minutes
	Hours
	Unknown
)

var _dummyGranularity = Seconds
var _ json.Marshaler = &_dummyGranularity
var _ json.Unmarshaler = &_dummyGranularity

const (
	seconds = "seconds"
	minutes = "minutes"
	hours   = "hours"
	unknown = "unknown"
)

var granularityNames = [...]string{
	seconds,
	minutes,
	hours,
}

var maxTimeGranularity = TimeGranularity(len(granularityNames) - 1)

func IsValid(i TimeGranularity) bool {
	return i >= 0 && i <= maxTimeGranularity
}

func FromName(s string) TimeGranularity {
	for idx, name := range granularityNames {
		if name == s {
			return TimeGranularity(idx)
		}
	}

	return Unknown
}

func ForEach(f func(TimeGranularity)) {
	for i := 0; i <= int(maxTimeGranularity); i++ {
		tg := TimeGranularity(i)
		f(tg)
	}
}

func (tg TimeGranularity) String() string {
	if !IsValid(tg) {
		return fmt.Sprintf("unknown(%d)", tg)
	}
	return granularityNames[tg]
}

func (tg *TimeGranularity) MarshalJSON() ([]byte, error) {
	if tg == nil {
		return nil, fmt.Errorf("cannot marshal unknown time granularity (nil)")
	}

	timeGran := *tg
	if !IsValid(timeGran) {
		return nil, fmt.Errorf("cannot marshal unknown time granularity (%d)", timeGran)
	}

	name := granularityNames[timeGran]
	b := make([]byte, 0, len(name)+2)
	b = append(b, '"')
	b = append(b, name...)
	return append(b, '"'), nil
}

func (tg *TimeGranularity) UnmarshalJSON(bytes []byte) error {
	if tg == nil {
		return fmt.Errorf("cannot unmarshal into nil TimeGranularity")
	}

	length := len(bytes)
	if length <= 2 || bytes[0] != '"' || bytes[length-1] != '"' {
		return fmt.Errorf("cannot unmarshal invalid JSON: `%s`", string(bytes))
	}

	unmarshalName := string(bytes[1 : length-1])
	timeGran := FromName(unmarshalName)
	if timeGran == Unknown {
		return fmt.Errorf("cannot unmarshal unknown time granularity `%s`", unmarshalName)
	}

	*tg = timeGran
	return nil
}

func (tg *TimeGranularity) UnmarshalForm(value string) error {
	if tg == nil {
		return fmt.Errorf("cannot unmarshal into nil TimeGranularity")
	}

	timeGran := FromName(value)
	if timeGran == Unknown {
		return fmt.Errorf("cannot unmarshal unknown time granularity `%s`", value)
	}

	*tg = timeGran
	return nil

}
