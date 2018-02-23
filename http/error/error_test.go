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

package error

import (
	"errors"
	"fmt"
	"testing"

	"github.com/turbinelabs/test/assert"
)

type errorCodeTestCase struct {
	f              func(string, ErrorCode) *Error
	expectedStatus int
}

func (tc *errorCodeTestCase) run(t *testing.T) {
	status := fmt.Sprintf("%d", tc.expectedStatus)
	assert.Group(status, t, func(g *assert.G) {
		err := tc.f("rat-a-tat-tat", UnknownUnclassifiedCode)
		assert.NonNil(g, err)
		assert.Equal(g, err.Message, "rat-a-tat-tat")
		assert.Equal(g, err.Code, UnknownUnclassifiedCode)
		assert.Equal(g, err.Status, tc.expectedStatus)
		assert.MatchesRegex(g, err.Error(), "rat-a-tat-tat")
	})
}

var testCases = []errorCodeTestCase{
	{New400, 400},
	{New404, 404},
	{New409, 409},
	{New500, 500},
}

func TestConstructors(t *testing.T) {
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestAuthorizationError(t *testing.T) {
	err := AuthorizationError()

	assert.MatchesRegex(t, err.Message, "not authorized")
	assert.Equal(t, err.Code, UnknownUnauthorizedCode)
	assert.Equal(t, err.Status, 403)
}

func TestFromErrorWithErrorPointer(t *testing.T) {
	x := &Error{Message: "piu piu", Code: UnknownUnclassifiedCode, Status: 503}

	err := FromError(x, UnknownTransportCode)
	assert.SameInstance(t, err, x)
}

func TestFromErrorWithGoError(t *testing.T) {
	x := errors.New("piu piu")

	err := FromError(x, UnknownTransportCode)
	assert.NonNil(t, err)

	expected := &Error{
		Message: x.Error(),
		Code:    UnknownTransportCode,
		Status:  500,
	}

	assert.DeepEqual(t, err, expected)
}

func TestFromErrorWithNil(t *testing.T) {
	assert.Nil(t, FromError(nil, UnknownTransportCode))
}
