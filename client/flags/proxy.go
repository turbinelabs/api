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
package flags

import (
	"github.com/turbinelabs/api/service"
	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/nonstdlib/flag/usage"
)

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

// ProxyFromFlags installs flags to obtain a proxy from a proxy name
type ProxyFromFlags interface {
	Name() string
	Ref(service.ZoneRef) service.ProxyRef
}

type pff struct {
	proxyName string
}

// NewProxyFromFlags installs a ProxyFromFlags into the given flagset
func NewProxyFromFlags(flagset tbnflag.FlagSet) ProxyFromFlags {
	ff := &pff{}

	flagset.StringVar(
		&ff.proxyName,
		"proxy-name",
		"",
		usage.Required("The name of the Proxy to configure"),
	)

	return ff
}

func (ff *pff) Name() string { return ff.proxyName }

func (ff *pff) Ref(zoneRef service.ZoneRef) service.ProxyRef {
	return service.NewProxyNameProxyRef(ff.proxyName, zoneRef)
}
