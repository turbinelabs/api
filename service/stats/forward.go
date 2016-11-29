package stats

// A Stat is a named, timestamped, and tagged data point.
type Stat struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"` // microseconds since the Unix epoch, UTC
	Tags      map[string]string `json:"tags,omitempty"`
}

// Payload is the payload of a stats update call.
type Payload struct {
	Source string `json:"source"`
	Stats  []Stat `json:"stats"`
}

// ForwardResult is a JSON-encodable struct that encapsulates the result of
// forwarding metrics.
type ForwardResult struct {
	NumAccepted int `json:"numAccepted"`
}
