package client

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/turbinelabs/api/fixture"
	apihttp "github.com/turbinelabs/api/http"
	"github.com/turbinelabs/api/http/envelope"
	"github.com/turbinelabs/api/service"
	tbnhttp "github.com/turbinelabs/client/http"
	"github.com/turbinelabs/test/assert"
)

const (
	clusterTestApiKey    = "whee-whee-whee"
	clusterCommonURL     = "/v1.0/cluster"
	domainCommonURL      = "/v1.0/domain"
	proxyCommonURL       = "/v1.0/proxy"
	routeCommonURL       = "/v1.0/route"
	sharedRulesCommonURL = "/v1.0/shared_rules"
	zoneCommonURL        = "/v1.0/zone"
	userCommonURL        = "/v1.0/admin/user"
)

var fixtures = fixture.DataFixtures

// Used for verifying http client tests. It does clever things to decide how
// to write out a response:
//
//   If response is a X the verifier handler writes Y:
//    string ------------- exactly those bytes
//    envelope.Response -- the marshaled version of that object
//    *envelope.Response - the marshaled version of that object
//    Something else ----- an envelope.Response with the response parameter as the "response" field of the envelope
type verifyingHandler struct {
	fn       func(apihttp.RichRequest)
	status   int
	response interface{}
}

func (w verifyingHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rr := apihttp.NewRichRequest(r)
	rrw := apihttp.RichResponseWriter{rw}

	w.fn(rr)

	if w.response != nil {
		switch t := w.response.(type) {
		case string:
			rw.WriteHeader(w.status)
			rw.Write([]byte(t))
		case envelope.Response:
			rrw.WriteEnvelope(t.Error, t.Payload)
		case *envelope.Response:
			rrw.WriteEnvelope(t.Error, t.Payload)
		default:
			rrw.WriteEnvelope(nil, w.response)
		}
	}
}

func assertURLPrefix(t *testing.T, url, prefix string) bool {
	if !assert.True(t, strings.HasPrefix(url, prefix)) {
		assert.Tracing(t).Errorf("%q is not prefixed by %q", url, prefix)
		return false
	}
	return true
}

func stripURLPrefix(url, prefix string) string {
	return url[len(prefix):]
}

func newTestEndpoint(host string, port int) tbnhttp.Endpoint {
	e, err := tbnhttp.NewEndpoint(tbnhttp.HTTP, host, port)
	if err != nil {
		log.Fatal(err)
	}
	return e
}

func newTestEndpointFromServer(server *httptest.Server) tbnhttp.Endpoint {
	u, e := url.Parse(server.URL)
	if e != nil {
		log.Fatal(e)
	}

	hostPortPair := strings.Split(u.Host, ":")
	host := hostPortPair[0]
	port, e := strconv.Atoi(hostPortPair[1])
	if e != nil {
		log.Fatal(e)
	}

	return newTestEndpoint(host, port)
}

func getAllInterface(server *httptest.Server) service.All {
	endpoint := newTestEndpointFromServer(server)
	serviceall, err := NewAll(endpoint, clusterTestApiKey, http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}
	return serviceall
}

func getAdminInterface(server *httptest.Server) service.Admin {
	endpoint := newTestEndpointFromServer(server)
	admin, err := NewAdmin(endpoint, clusterTestApiKey, http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}
	return admin
}
