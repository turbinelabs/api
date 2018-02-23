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

//go:generate genny -in object.genny -out gen_cluster.go gen "__type__=cluster __Type__=Cluster __snake__=cluster __Coll__=Clusters __root__="
//go:generate genny -in object.genny -out gen_domain.go gen "__type__=domain __Type__=Domain __snake__=domain __Coll__=Domains __root__="
//go:generate genny -in object.genny -out gen_proxy.go gen "__type__=proxy __Type__=Proxy __snake__=proxy __Coll__=Proxies __root__="
//go:generate genny -in object.genny -out gen_route.go gen "__type__=route __Type__=Route __snake__=route __Coll__=Routes __root__="
//go:generate genny -in object.genny -out gen_shared_rules.go gen "__type__=sharedRules __snake__=shared_rules __Type__=SharedRules __Coll__=SharedRulesSlice __root__="
//go:generate genny -in object.genny -out gen_zone.go gen "__type__=zone __Type__=Zone __snake__=zone __Coll__=Zones __root__="

// admin-rooted objects

//go:generate genny -in object.genny -out gen_user.go gen "__type__=user __Type__=User __snake__=user __Coll__=Users __root__=/admin"
// access token not generated because it doesn't expose a general Modify endpoint
