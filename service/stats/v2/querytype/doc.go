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

// Package querytype defines the QueryType enumeration.
package querytype

//go:generate codegen --output=query_type.go --source=$GOFILE ../../enum.template type=querytype.QueryType values[]=Unknown,Requests,Responses,Success,Error,Failure,LatencyP50,LatencyP99,SuccessRate,ResponsesForCode,DownstreamRequests,DownstreamResponses,DownstreamSuccess,DownstreamError,DownstreamFailure,DownstreamLatencyP50,DownstreamLatencyP99,DownstreamSuccessRate,DownstreamResponsesForCode

//go:generate codegen --output=query_type_test.go --source=$GOFILE ../../enum_test.template type=querytype.QueryType values[]=Unknown,Requests,Responses,Success,Error,Failure,LatencyP50,LatencyP99,SuccessRate,ResponsesForCode,DownstreamRequests,DownstreamResponses,DownstreamSuccess,DownstreamError,DownstreamFailure,DownstreamLatencyP50,DownstreamLatencyP99,DownstreamSuccessRate,DownstreamResponsesForCode
