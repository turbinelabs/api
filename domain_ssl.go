package api

import (
	"fmt"
	"strings"
)

// SSLProtocol is a name of a SSL protocol that may be used by a domain.
type SSLProtocol string

// SSLConfig handles configuring support for SSL termination on a domain.
//
// CertKeyPairs is represented as an array but we currently support only one
// certificate being specified.
//
// At present I'm unsure how much flexibility will be needed so I'm supporting a
// basic set of knobs. Not yet exposed are things like DH params, EC selection,
// etc.)
type SSLConfig struct {
	CipherFilter string            `json:"cipher_filter"`
	Protocols    []SSLProtocol     `json:"protocols"`
	CertKeyPairs []CertKeyPathPair `json:"cert_key_pairs"`
}

// CertKeyPathPair is a container that binds a certificate path to a key path.
// Both of these must be specified.
type CertKeyPathPair struct {
	CertificatePath string `json:"certificate_path"`
	KeyPath         string `json:"key_path"`
}

const (
	SSL2   SSLProtocol = "SSLv2"
	SSL3   SSLProtocol = "SSLv3"
	TLS1   SSLProtocol = "TLSv1"
	TLS1_1 SSLProtocol = "TLSv1.1"
	TLS1_2 SSLProtocol = "TLSv1.2"

	// DefaultCipherFilter chooses the default set of ciphers that may be used for
	// communicating with a domain.
	DefaultCipherFilter = "EECDH+AESGCM:EDH+AESGCM:AES256+EECDH:AES256+EDH"
)

var sslProtoName = map[SSLProtocol]string{
	SSL2:   "SSLv2",
	SSL3:   "SSLv3",
	TLS1:   "TLSv1",
	TLS1_1: "TLSv1.1",
	TLS1_2: "TLSv1.2",
}

// DefaultProtocols indicates which protocols will be supported if none
// are specified in the SSLConfig object for a domain.
var DefaultProtocols = []SSLProtocol{TLS1, TLS1_1, TLS1_2}

func (c SSLConfig) Equals(o SSLConfig) bool {
	if strings.TrimSpace(c.CipherFilter) != strings.TrimSpace(o.CipherFilter) ||
		len(c.Protocols) != len(o.Protocols) ||
		len(c.CertKeyPairs) != len(o.CertKeyPairs) {
		return false
	}

	protos := map[SSLProtocol]bool{}

	for _, p := range c.Protocols {
		protos[p] = true
	}

	for _, p := range o.Protocols {
		if !protos[p] {
			return false
		}
	}

	certs := map[string]bool{}
	keys := map[string]bool{}

	for _, ckp := range c.CertKeyPairs {
		certs[ckp.CertificatePath] = true
		keys[ckp.KeyPath] = true
	}

	for _, ckp := range o.CertKeyPairs {
		if !certs[ckp.CertificatePath] ||
			!keys[ckp.KeyPath] {
			return false
		}
	}

	return true
}

func (s SSLConfig) IsValid() *ValidationError {
	errs := &ValidationError{}

	for _, e := range s.Protocols {
		if _, ok := sslProtoName[e]; !ok {
			errs.AddNew(ErrorCase{
				"ssl_config.protocols",
				fmt.Sprintf("invalid protocol specified %v", e)})
		}
	}

	// if we ever support multiple certs we should verify each cert appears
	// only once
	if len(s.CertKeyPairs) != 1 {
		errs.AddNew(ErrorCase{
			"ssl_config.cert_key_pairs",
			"a single SSL certificate and key pair must be specified"})
	} else {
		kp := s.CertKeyPairs[0]
		parent := fmt.Sprintf(
			"ssl_config.cert_key_pairs[%v].",
			kp.CertificatePath)

		errCheckIndex(kp.CertificatePath, errs, parent+"certificate_path")

		if strings.TrimSpace(kp.KeyPath) == "" {
			errs.AddNew(ErrorCase{parent + "key_path", "may not be empty"})
		}
	}

	return errs.OrNil()
}
