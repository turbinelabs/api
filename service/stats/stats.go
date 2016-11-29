package stats

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"strings"
	"time"

	"github.com/turbinelabs/nonstdlib/stats"
	tbntime "github.com/turbinelabs/nonstdlib/time"
)

// StatsService forwards stats data to a remote stats-server.
type StatsService interface {
	// Forward the given stats payload.
	Forward(*Payload) (*ForwardResult, error)

	// Query for stats
	Query(*Query) (*QueryResult, error)

	// Closes the client and releases any resources it created.
	Close() error
}

// Creates a nonstdlib/stats.Stats that uses this StatsService to forward
// arbitrary stats with the given source and optional scopes.
func AsStats(svc StatsService, source string, scopes ...string) stats.Stats {
	resolvedScope := strings.Join(scopes, "/")

	return &asStats{
		svc:    svc,
		source: source,
		scope:  resolvedScope,
	}
}

type asStats struct {
	svc    StatsService
	source string
	scope  string
}

var _ stats.Stats = &asStats{}

func (s *asStats) stat(name string, value float64) error {
	if s.scope != "" {
		name = s.scope + "/" + name
	}

	payload := &Payload{
		Source: s.source,
		Stats: []Stat{
			{
				Name:      name,
				Value:     value,
				Timestamp: tbntime.ToUnixMicro(time.Now()),
			},
		},
	}
	_, err := s.svc.Forward(payload)
	return err
}

func (s *asStats) Inc(name string, v int64) error {
	return s.stat(name, float64(v))
}

func (s *asStats) Gauge(name string, v int64) error {
	return s.stat(name, float64(v))
}

func (s *asStats) TimingDuration(name string, d time.Duration) error {
	return s.stat(name, d.Seconds())
}
