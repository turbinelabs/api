package header

import (
	"net/http"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestAllHeadersAreCanonical(t *testing.T) {
	for _, h := range headers {
		assert.Group(h, t, func(g *assert.G) {
			assert.Equal(g, http.CanonicalHeaderKey(h), h)
		})
	}
}
