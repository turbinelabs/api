// Contains the Response envelope which the http service and api server use to
// encapsulate server behavior.
package envelope

import httperr "github.com/turbinelabs/api/http/error"

// Response is constructed at API render time to enable a predictable way to
// transmit error and request payload to a HTTP client. It is received by the
// HTTP client and unpacked into the appropriate types depending on the call
// being made.
type Response struct {
	Error   *httperr.Error `json:"error,omitempty"`
	Payload interface{}    `json:"result,omitempty"`
}
