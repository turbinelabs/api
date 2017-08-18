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

// Package stats defines the interfaces representing the portion of the
// Turbine Labs public API prefixed by /v1.0/stats
package stats

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

// StatsService forwards stats data to a remote stats-server.
type StatsService interface {
	// Forward the given stats payload.
	Forward(*Payload) (*ForwardResult, error)

	// Query for stats
	Query(*Query) (*QueryResult, error)

	// Closes the client and releases any resources it created.
	Close() error
}

// StatsServiceV2 forwards stats data to a remote stats-server using
// the version 2 forwarding interface.
type StatsServiceV2 interface {
	// Forward the given stats payload.
	ForwardV2(*PayloadV2) (*ForwardResult, error)

	// Query for stats
	Query(*Query) (*QueryResult, error)

	// Closes the client and releases any resources it created.
	Close() error
}
