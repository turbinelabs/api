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

// Package timegranularity defines the TimeGranularity enumeration.
//
// Note that the Seconds TimeGranularity is not currently supported.
package timegranularity

//go:generate codegen --output=time_granularity.go --source=$GOFILE ../../enum.template type=timegranularity.TimeGranularity values[]=Seconds,Minutes,Hours,Unknown

//go:generate codegen --output=time_granularity_test.go --source=$GOFILE ../../enum_test.template type=timegranularity.TimeGranularity values[]=Seconds,Minutes,Hours,Unknown
