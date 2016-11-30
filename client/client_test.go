package client

import (
	"fmt"
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
	apiheader "github.com/turbinelabs/api/http/header"
	"github.com/turbinelabs/api/service"
	"github.com/turbinelabs/test/assert"
)

const (
	clientTestApiKey     = "whee-whee-whee"
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

	apiKey := rr.Underlying().Header.Get(http.CanonicalHeaderKey(apiheader.APIKey))
	if apiKey != clientTestApiKey {
		rw.WriteHeader(400)
		rw.Write([]byte(
			fmt.Sprintf(
				"wrong api key header, got %s, want %s",
				apiKey,
				clientTestApiKey,
			),
		))
		return
	}

	clientId := rr.Underlying().Header.Get(http.CanonicalHeaderKey(apiheader.ClientID))
	if clientId != apiClientID {
		rw.WriteHeader(400)
		rw.Write([]byte(
			fmt.Sprintf(
				"wrong client id header: got %s, want %s",
				clientId,
				apiClientID,
			),
		))
		return
	}

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

func newTestEndpoint(host string, port int) apihttp.Endpoint {
	e, err := apihttp.NewEndpoint(apihttp.HTTP, host, port)
	if err != nil {
		log.Fatal(err)
	}
	return e
}

func newTestEndpointFromServer(server *httptest.Server) apihttp.Endpoint {
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
	serviceall, err := NewAll(endpoint, clientTestApiKey)
	if err != nil {
		log.Fatal(err)
	}
	return serviceall
}

func getAdminInterface(server *httptest.Server) service.Admin {
	endpoint := newTestEndpointFromServer(server)
	admin, err := NewAdmin(endpoint, clientTestApiKey)
	if err != nil {
		log.Fatal(err)
	}
	return admin
}
