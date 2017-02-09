
[//]: # ( Copyright 2017 Turbine Labs, Inc.                                   )
[//]: # ( you may not use this file except in compliance with the License.    )
[//]: # ( You may obtain a copy of the License at                             )
[//]: # (                                                                     )
[//]: # (     http://www.apache.org/licenses/LICENSE-2.0                      )
[//]: # (                                                                     )
[//]: # ( Unless required by applicable law or agreed to in writing, software )
[//]: # ( distributed under the License is distributed on an "AS IS" BASIS,   )
[//]: # ( WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or     )
[//]: # ( implied. See the License for the specific language governing        )
[//]: # ( permissions and limitations under the License.                      )

# turbinelabs/api

[![Apache 2.0](https://img.shields.io/hexpm/l/plug.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/turbinelabs/api?status.svg)](https://godoc.org/github.com/turbinelabs/api)
[![CircleCI](https://circleci.com/gh/turbinelabs/api.svg?style=shield)](https://circleci.com/gh/turbinelabs/api)
[![Go Report Card](https://goreportcard.com/badge/github.com/turbinelabs/api)](https://goreportcard.com/report/github.com/turbinelabs/apicli)

The api project provides Go bindings for the Turbine Labs public API. We use
these bindings in both our open-sourced projects and our private server
implementations; as such this project represents the source of truth for how
our public API is defined.

We also provide a swagger definition, which backs our
[public API documentation](https://docs.turbinelabs.io/docs/versions/1.0/) and
[API explorer](https://docs.turbinelabs.io/api-explorer/turbine-labs/versions/1.0/)
and may be used to generate additional clients in the future.

## Requirements

- Go 1.7.4 or later (previous versions may work, but we don't build or test against them)

## Dependencies

The api project depends on the
[nonstdlib](https://github.com/turbinelabs/nonstdlib) package. The tests depend
on our [test package](https://github.com/turbinelabs/test), and on
[gomock](https://github.com/golang/mock), and gomock-based Mocks of
most interfaces are provided.

## Install

```
go get -u github.com/turbinelabs/api/...
```

## Clone/Test

```
mkdir -p $GOPATH/src/turbinelabs
git clone https://github.com/turbinelabs/api.git > $GOPATH/src/turbinelabs/api
go test github.com/turbinelabs/api/...
```

## Godoc

[`api`](https://godoc.org/github.com/turbinelabs/api)

## Versioning

Please see [Versioning of Turbine Labs Open Source Projects](http://github.com/turbinelabs/developer/blob/master/README.md#versioning).

## Pull Requests

Patches accepted! Please see
[Contributing to Turbine Labs Open Source Projects](http://github.com/turbinelabs/developer/blob/master/README.md#contributing).

## Code of Conduct

All Turbine Labs open-sourced projects are released with a
[Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in our
projects you agree to abide by its terms, which will be carefully enforced.
