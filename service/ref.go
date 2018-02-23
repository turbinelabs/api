/*
Copyright 2017 Turbine Labs, Inc.

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

package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/turbinelabs/api"
)

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

// ProxyRef encapsulates a lookup of a Proxy by Name and Zone Name
type ProxyRef interface {
	// Get returns the Proxy corresponding to the ProxyRef. The lookup is memoized.
	Get(All) (api.Proxy, error)
	Name() string
	ZoneRef() ZoneRef

	// MapKey returns a string suitable for keying the ProxyRef
	// in a map. ProxyRefs with the same MapKey refer to the same
	// Proxy.
	MapKey() string
}

type proxyRefMapKey struct {
	ProxyName string `json:"proxy_name"`
	ZoneName  string `json:"zone_name"`
}

type proxyRef struct {
	p       *api.Proxy
	name    string
	zoneRef ZoneRef
}

func (r *proxyRef) set(p *api.Proxy) {
	r.p = p
	r.name = p.Name
}

func (r *proxyRef) Get(svc All) (api.Proxy, error) {
	if r.p != nil {
		return *r.p, nil
	}
	if r.name == "" {
		return api.Proxy{}, errors.New("proxyName must be non-empty")
	}
	z, err := r.zoneRef.Get(svc)
	if err != nil {
		return api.Proxy{}, err
	}
	ps, err := svc.Proxy().Index(ProxyFilter{ZoneKey: z.ZoneKey, Name: r.name})
	if err != nil {
		return api.Proxy{}, err
	}
	if len(ps) == 0 {
		return api.Proxy{}, fmt.Errorf("no Proxy found for name %q", r.name)
	}
	r.set(&ps[0])
	return ps[0], nil
}

func (r *proxyRef) Name() string {
	return r.name
}

func (r *proxyRef) ZoneRef() ZoneRef {
	return r.zoneRef
}

func (r *proxyRef) MapKey() string {
	key := proxyRefMapKey{
		ProxyName: r.Name(),
		ZoneName:  r.zoneRef.Name(),
	}
	bs, err := json.Marshal(key)
	if err != nil {
		// this really should never ever happen
		panic(err)
	}
	return string(bs)
}

// NewProxyRef produces a ProxyRef from an api.Proxy and an api.Zone
func NewProxyRef(p api.Proxy, z api.Zone) ProxyRef {
	r := &proxyRef{zoneRef: NewZoneRef(z)}
	r.set(&p)
	return r
}

// NewProxyNameProxyRef returns a ProxyRef keyed by the given api.Proxy name and
// ZoneRef
func NewProxyNameProxyRef(name string, zRef ZoneRef) ProxyRef {
	return &proxyRef{
		name:    name,
		zoneRef: zRef,
	}
}

// NewProxyRefFromMapKey returns a ProxyRef keyed by the given map key
func NewProxyRefFromMapKey(keyStr string) (ProxyRef, error) {
	key := proxyRefMapKey{}
	if err := json.Unmarshal([]byte(keyStr), &key); err != nil {
		return nil, err
	}
	return NewProxyNameProxyRef(key.ProxyName, NewZoneNameZoneRef(key.ZoneName)), nil
}

// ZoneRef encapsulates a lookup of a Zone by Name
type ZoneRef interface {
	// Get returns the Zone corresponding to the ZoneRef. The lookup is memoized.
	Get(All) (api.Zone, error)

	Name() string
}

// NewZoneRef produces a ZoneRef from an api.Zone
func NewZoneRef(z api.Zone) ZoneRef {
	r := &zoneRef{}
	r.set(&z)
	return r
}

// NewZoneNameZoneRef produces a ZoneRef from an api.Zone name
func NewZoneNameZoneRef(name string) ZoneRef {
	return &zoneRef{name: name}
}

type zoneRef struct {
	z    *api.Zone
	name string
}

func (r *zoneRef) set(z *api.Zone) {
	r.z = z
	r.name = z.Name
}

func (r *zoneRef) Get(svc All) (api.Zone, error) {
	if r.z != nil {
		return *r.z, nil
	}
	if r.name == "" {
		return api.Zone{}, errors.New("zoneName must be non-empty")
	}
	zs, err := svc.Zone().Index(ZoneFilter{Name: r.name})
	if err != nil {
		return api.Zone{}, err
	}
	if len(zs) == 0 {
		return api.Zone{}, fmt.Errorf("no Zone found for name %q", r.name)
	}
	r.set(&zs[0])
	return zs[0], nil
}

func (r *zoneRef) Name() string {
	return r.name
}
