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
	"net"
)

type ListenerProtocol string

const (
	HttpListenerProtocol     ListenerProtocol = "http"
	Http2ListenerProtocol    ListenerProtocol = "http2"
	HttpAutoListenerProtocol ListenerProtocol = "http_auto"
	TCPListenerProtocol      ListenerProtocol = "tcp"
)

func ListenerProtocolFromString(s string) (ListenerProtocol, error) {
	lp := ListenerProtocol(s)
	switch lp {
	case HttpListenerProtocol,
		Http2ListenerProtocol,
		TCPListenerProtocol,
		HttpAutoListenerProtocol:
		return lp, nil
	}
	return ListenerProtocol(""), fmt.Errorf("unknown ListenerProtocol: %s", s)
}

func (lp ListenerProtocol) IsValid() bool {
	_, err := ListenerProtocolFromString(string(lp))
	return err == nil
}

type ListenerKey string

// A Listener represents a port Envoy will listen on
type Listener struct {
	ListenerKey   ListenerKey      `json:"listener_key"` // overwritten for create
	ZoneKey       ZoneKey          `json:"zone_key"`
	Name          string           `json:"name"`
	IP            string           `json:"ip"`
	Port          int              `json:"port"`
	Protocol      ListenerProtocol `json:"protocol"`
	DomainKeys    []DomainKey      `json:"domain_keys"`
	TracingConfig *TracingConfig   `json:"tracing_config"`
	OrgKey        OrgKey           `json:"-"`
	Checksum
}

func (l Listener) GetZoneKey() ZoneKey   { return l.ZoneKey }
func (l Listener) GetOrgKey() OrgKey     { return l.OrgKey }
func (l Listener) Key() string           { return string(l.ListenerKey) }
func (l Listener) GetChecksum() Checksum { return l.Checksum }

func (l Listener) IsNil() bool {
	return l.Equals(Listener{})
}

type Listeners []Listener

// Checks for validity of a listener. A listener is considered valid if it has a:
//  1. ListenerKey OR is being checked before creation
//  2. non-empty ZoneKey
//  3. non-empty Name
//  4. non-empty IP
//  5. non-zero Port
//  6. non-empty Protocol, one of http, http2, http_auto or tcp
func (l Listener) IsValid() *ValidationError {
	scope := func(s string) string { return "listener." + s }
	ecase := func(f, m string) ErrorCase {
		return ErrorCase{scope(f), m}
	}

	errs := &ValidationError{}

	if !l.Protocol.IsValid() {
		errs.AddNew(ecase(
			"protocol",
			fmt.Sprintf("%s is not a valid listener protocol", string(l.Protocol))))
	}

	errCheckKey(string(l.ListenerKey), errs, scope("listener_key"))
	errCheckKey(string(l.ZoneKey), errs, scope("zone_key"))
	errCheckIndex(l.Name, errs, scope("name"))
	errCheckKey(string(l.OrgKey), errs, scope("org_key"))

	seenDomain := map[string]bool{}
	for _, dk := range l.DomainKeys {
		sdk := string(dk)
		if seenDomain[sdk] {
			errs.AddNew(ErrorCase{scope("domain_keys"), fmt.Sprintf("duplicate domain key '%v'", sdk)})
		}
		seenDomain[sdk] = true
		errCheckKey(sdk, errs, fmt.Sprintf("listener.domain_keys[%v]", sdk))
	}

	if len(l.IP) < 1 {
		errs.AddNew(ecase("ip", "must be specified"))
	}

	if net.ParseIP(l.IP) == nil {
		errs.AddNew(ecase("host", fmt.Sprintf("%s is not a valid ip", l.IP)))
	}

	if l.Port <= 0 {
		errs.AddNew(ecase("port", "must be positive"))
	}

	if !l.Protocol.IsValid() {
		errs.AddNew(ecase("protocol", "must be one of http, http2, http_auto, or tcp"))
	}

	if l.TracingConfig != nil {
		errs.MergePrefixed(l.TracingConfig.IsValid(), scope("tracing_config"))
	}
	return errs.OrNil()
}

// Check if all fields of this listener are exactly equal to fields of another
// listener.
func (l Listener) Equals(o Listener) bool {
	lTCNil := l.TracingConfig == nil
	oTCNil := o.TracingConfig == nil

	if lTCNil != oTCNil {
		return false
	}
	tcEq := oTCNil || l.TracingConfig.Equals(*o.TracingConfig)

	if len(l.DomainKeys) != len(o.DomainKeys) {
		return false
	}

	hasDomain := make(map[DomainKey]bool)

	for _, dk := range l.DomainKeys {
		hasDomain[dk] = true
	}

	for _, dk := range o.DomainKeys {
		if !hasDomain[dk] {
			return false
		}
	}

	return l.ListenerKey == o.ListenerKey &&
		l.ZoneKey == o.ZoneKey &&
		l.Name == o.Name &&
		l.IP == o.IP &&
		l.Port == o.Port &&
		l.Protocol == o.Protocol &&
		l.Checksum.Equals(o.Checksum) &&
		l.OrgKey == o.OrgKey &&
		tcEq
}
