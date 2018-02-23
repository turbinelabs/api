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

package service

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"time"

	"github.com/turbinelabs/api"
	"github.com/turbinelabs/api/service/changelog"
)

type History interface {
	// GET /v1.0/changelog/adhoc[?filter=<json>&start=time_utc_usec>&stop=<time_utc_usec>]
	//
	// Index provides a generalized search interface for changes that the API has
	// recorded. The filter to be applied should be a FilterOrs as described in
	// changelog docs. It may be URL encoded JSON and set as the value to the
	// 'filter' query param. Humane encoding is also available; if used the query
	// param names and values are derived from the encoded struct FilterOrs and
	// are based on the `form` struct tag.
	//
	// If no filters are specified when making the HTTP call all changes in the
	// org of the requesting user will be returned within the time window.
	//
	// If neither start nor end time is specified the default will be treated as
	// the previous three hours to current time. If only one endpoint is provided
	// the other will be inferred such that a three hour window is examined.
	Index(
		filters changelog.FilterExpr,
		start,
		end time.Time,
	) ([]api.ChangeDescription, error)

	// GET /v1.0/changelog/domain-graph/<key>[?start=<time_utc_usec>&stop=<time_utc_usec>]
	//
	// DomainGraph returns any changes within the object graph of the domain
	// specified by domainKey. Specifically this includes changes to:
	//
	//  1. the domain itself
	//  2. the set of proxies which route to the domain (note: doesn't include
	//     changes to the proxies themselves only the routing set)
	//  3. the routes that were a part of the domain during any part of the
	//     window (including routes which were initially added then removed
	//     during the window)
	//  4. any cluster referenced by any rules that was a part of the domain
	//     during the window.
	//
	// If one of the window sides are unset it will be filled to the default
	// window size. If both window edges are zero-value they will be set to
	// a window of the default size ending at the current time.
	//
	// If the window (stop - start) exceeds the maximum size or stop < start
	// an error will be returned.
	//
	// The maximum duration is 24 hours. The default window size is 1 hour.
	//
	// TODO: pull in org key hint (for query efficiency)
	//   https://github.com/turbinelabs/tbn/issues/1022
	DomainGraph(
		domainKey api.DomainKey,
		start,
		stop time.Time,
	) ([]api.ChangeDescription, error)

	// GET /v1.0/changelog/route-graph/<key>[?start=<time_utc_usec>&stop=<time_utc_usec>]
	//
	// RouteGraph returns any changes within a set window on a route or the clusters
	// within that route.
	//
	// If the window (stop - start) exceeds the maximum size or stop < start
	// an error will be returned.
	//
	// The maximum duration is 24 hours. The default window size is 1 hour.
	RouteGraph(
		routeKey api.RouteKey,
		start,
		stop time.Time,
	) ([]api.ChangeDescription, error)

	// GET /v1.0/changelog/shared-rules-graph/<key>[?start=<time_utc_usec>&stop=<time_utc_usec>]
	//
	// SharedRulesGraph returns any changes to a SharedRules object within a set window.
	//
	// If the window (stop - start) exceeds the maximum size or stop < start an
	// error will be returned.
	//
	// The maximum duration is 24 hours. The default window size is 1 hour.
	SharedRulesGraph(
		sharedRulesKey api.SharedRulesKey,
		start,
		stop time.Time,
	) ([]api.ChangeDescription, error)

	// GET /v1.0/changelog/cluster-graph/<key>[?start=<time_utc_usec>&stop=<time_utc_usec>]
	//
	// ClusterGraph returns any changes to a cluster and any domains that have
	// started, or stopped, routing to a cluster within a set window.
	//
	// If the window (stop - start) exceeds the maximum size or stop < start
	// an error will be returned.
	//
	// The maximum duration of 24 hours. The default window size is 1 hour.
	ClusterGraph(
		clusterKey api.ClusterKey,
		start,
		stop time.Time,
	) ([]api.ChangeDescription, error)

	// GET /v1.0/changelog/zone/<key>[?start=<time_utc_usec>&stop=<time_utc_usec>]
	//
	// Zone returns any changes within a Zone that occurred during a window.
	//
	// If the window (stop - start) exceeds the maximum size or stop < start
	// an error will be returned.
	//
	// The maximum duration of 24 hours. The default window size is 1 hour.
	Zone(
		zoneKey api.ZoneKey,
		start,
		stop time.Time,
	) ([]api.ChangeDescription, error)
}
