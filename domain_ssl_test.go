/*
Copyright 2018 Turbine Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"fmt"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func getSSLConfig() SSLConfig {
	return SSLConfig{
		CipherFilter: DefaultCipherFilter,
		Protocols:    DefaultProtocols,
		CertKeyPairs: []CertKeyPathPair{{
			"/path/to/cert.pem",
			"/path/to/key.pem",
		}},
	}
}

func TestSSLConfigEquals(t *testing.T) {
	s1 := getSSLConfig()
	s2 := getSSLConfig()

	assert.True(t, s1.Equals(s2))
	assert.True(t, s2.Equals(s1))
}

func TestSSLConfigEqualsCiphers(t *testing.T) {
	s1 := getSSLConfig()
	s2 := getSSLConfig()
	s2.CipherFilter = "other filter set"

	assert.False(t, s1.Equals(s2))
	assert.False(t, s2.Equals(s1))
}

func TestSSLConfigEqualsProtocols(t *testing.T) {
	s1 := getSSLConfig()
	s2 := getSSLConfig()
	s2.Protocols = append(s2.Protocols, SSL3)
	assert.False(t, s1.Equals(s2))
	assert.False(t, s2.Equals(s1))
}

func TestSSLConfigEqualsCKP(t *testing.T) {
	s1 := getSSLConfig()
	s2 := getSSLConfig()
	s2.CertKeyPairs = []CertKeyPathPair{}
	assert.False(t, s1.Equals(s2))
	assert.False(t, s2.Equals(s1))
}

func TestSSLConfigEqualsCKPCertPath(t *testing.T) {
	s1 := getSSLConfig()
	s2 := getSSLConfig()
	s2.CertKeyPairs[0].CertificatePath = "aoeu"
	assert.False(t, s1.Equals(s2))
	assert.False(t, s2.Equals(s1))
}

func TestSSLConfigEqualsCKPKeyPath(t *testing.T) {
	s1 := getSSLConfig()
	s2 := getSSLConfig()
	s2.CertKeyPairs[0].KeyPath = "aoeu"
	assert.False(t, s1.Equals(s2))
	assert.False(t, s2.Equals(s1))
}

func TestSSLConfigIsValid(t *testing.T) {
	s1 := getSSLConfig()
	assert.Nil(t, s1.IsValid())
}

func TestSSLConfigIsValidNoCiphers(t *testing.T) {
	s1 := getSSLConfig()
	s1.CipherFilter = ""
	assert.Nil(t, s1.IsValid())
}

func TestSSLConfigIsValidNoProtocols(t *testing.T) {
	s1 := getSSLConfig()
	s1.Protocols = []SSLProtocol{}
	assert.Nil(t, s1.IsValid())
}

func TestSSLConfigIsValidNilProtocols(t *testing.T) {
	s1 := getSSLConfig()
	s1.Protocols = nil
	assert.Nil(t, s1.IsValid())
}

func TestSSLConfigEqualsNoCertPath(t *testing.T) {
	s1 := getSSLConfig()
	s1.CertKeyPairs[0].CertificatePath = ""
	assert.DeepEqual(t, s1.IsValid(), &ValidationError{[]ErrorCase{
		{"ssl_config.cert_key_pairs[].certificate_path", "may not be empty"},
	}})
}

func TestSSLConfigEqualsNoKeyPath(t *testing.T) {
	s1 := getSSLConfig()
	kp := s1.CertKeyPairs[0]
	s1.CertKeyPairs[0].KeyPath = ""
	assert.DeepEqual(t, s1.IsValid(), &ValidationError{[]ErrorCase{{
		fmt.Sprintf("ssl_config.cert_key_pairs[%v].key_path", kp.CertificatePath),
		"may not be empty",
	}}})
}

func TestSSLConfigEqualsMultipleCKPairs(t *testing.T) {
	s1 := getSSLConfig()
	s1.CertKeyPairs = append(s1.CertKeyPairs, CertKeyPathPair{"1234", "234"})
	assert.DeepEqual(t, s1.IsValid(), &ValidationError{[]ErrorCase{
		{"ssl_config.cert_key_pairs", "a single SSL certificate and key pair must be specified"},
	}})
}
