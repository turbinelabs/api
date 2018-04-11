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

package client

//go:generate codegen --output=gen_cluster.go object.template Key=github.com/turbinelabs/api.ClusterKey Object=github.com/turbinelabs/api.Cluster ObjectArray=github.com/turbinelabs/api.Clusters Root=
//go:generate codegen --output=gen_domain.go object.template Key=github.com/turbinelabs/api.DomainKey Object=github.com/turbinelabs/api.Domain ObjectArray=github.com/turbinelabs/api.Domains Root=
//go:generate codegen --output=gen_proxy.go object.template Key=github.com/turbinelabs/api.ProxyKey Object=github.com/turbinelabs/api.Proxy ObjectArray=github.com/turbinelabs/api.Proxies Root=
//go:generate codegen --output=gen_route.go object.template Key=github.com/turbinelabs/api.RouteKey Object=github.com/turbinelabs/api.Route ObjectArray=github.com/turbinelabs/api.Routes Root=
//go:generate codegen --output=gen_shared_rules.go object.template Key=github.com/turbinelabs/api.SharedRulesKey Object=github.com/turbinelabs/api.SharedRules ObjectArray=github.com/turbinelabs/api.SharedRulesSlice Root=
//go:generate codegen --output=gen_zone.go object.template Key=github.com/turbinelabs/api.ZoneKey Object=github.com/turbinelabs/api.Zone ObjectArray=github.com/turbinelabs/api.Zones Root=

// admin-rooted objects

//go:generate codegen --output=gen_user.go object.template Key=github.com/turbinelabs/api.UserKey Object=github.com/turbinelabs/api.User ObjectArray=github.com/turbinelabs/api.Users Root=/admin
// access token not generated because it doesn't expose a general Modify endpoint
