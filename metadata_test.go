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

package api

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestMetadataAsMap(t *testing.T) {
	metadata := Metadata{{"foo", "bar"}, {"baz", "blegga"}}
	got := metadata.Map()
	want := map[string]string{"foo": "bar", "baz": "blegga"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("(%q).AsMap == %q, want %q", metadata, got, want)
	}
}

func TestMetadatumEqualsTrue(t *testing.T) {
	m1 := Metadatum{"Key", "Value"}
	m2 := Metadatum{"Key", "Value"}

	assert.True(t, m2.Equals(m1))
	assert.True(t, m1.Equals(m2))
}

func TestMetadatumValueChanged(t *testing.T) {
	m1 := Metadatum{"Key", "Value2"}
	m2 := Metadatum{"Key", "Value"}

	assert.False(t, m1.Equals(m2))
	assert.False(t, m2.Equals(m1))
}

func TestMetadatumKeyChanged(t *testing.T) {
	m1 := Metadatum{"Key", "Value"}
	m2 := Metadatum{"Key2", "Value"}

	assert.False(t, m1.Equals(m2))
	assert.False(t, m2.Equals(m1))
}

// Metadata
func TestMetadataZeroNil(t *testing.T) {
	m1 := Metadata{}
	var m2 Metadata

	assert.True(t, m1.Equals(m2))
	assert.True(t, m2.Equals(m1))
}

func TestMetadataEqualsNilNil(t *testing.T) {
	var m1, m2 Metadata
	assert.True(t, m1.Equals(m2))
	assert.True(t, m2.Equals(m1))
}

func TestMetadataCompare(t *testing.T) {
	type testcase struct {
		mda  Metadata
		mdb  Metadata
		want int
	}

	for _, tc := range []testcase{
		{want: 0},
		{mda: Metadata{}, want: 0},
		{mdb: Metadata{}, want: 0},
		{
			mda:  Metadata{{Key: "a", Value: "b"}},
			want: 1,
		},
		{
			mdb:  Metadata{{Key: "a", Value: "b"}},
			want: -1,
		},
		{
			mda:  Metadata{{Key: "a", Value: "b"}},
			mdb:  Metadata{{Key: "a", Value: "b"}},
			want: 0,
		},
		{
			mda:  Metadata{{Key: "a", Value: "b"}},
			mdb:  Metadata{{Key: "a", Value: "c"}},
			want: -1,
		},
		{
			mda:  Metadata{{Key: "a", Value: "c"}},
			mdb:  Metadata{{Key: "a", Value: "b"}},
			want: 1,
		},
		{
			mda:  Metadata{{Key: "a", Value: "b"}, {Key: "c", Value: "d"}},
			mdb:  Metadata{{Key: "c", Value: "d"}, {Key: "a", Value: "b"}},
			want: 0,
		},
		{
			mda:  Metadata{{Key: "a", Value: "b"}, {Key: "c", Value: "d"}},
			mdb:  Metadata{{Key: "a", Value: "b"}},
			want: 1,
		},
		{
			mda:  Metadata{{Key: "a", Value: "b"}},
			mdb:  Metadata{{Key: "c", Value: "d"}, {Key: "a", Value: "b"}},
			want: -1,
		},
		{
			mda:  Metadata{{Key: "e", Value: "f"}},
			mdb:  Metadata{{Key: "c", Value: "d"}, {Key: "a", Value: "b"}},
			want: 1,
		},
		{
			mda:  Metadata{{Key: "a", Value: "b"}, {Key: "c", Value: "d"}},
			mdb:  Metadata{{Key: "e", Value: "f"}},
			want: -1,
		},
		{
			mda:  Metadata{{Key: "a", Value: "b"}},
			mdb:  Metadata{{Key: "c", Value: "d"}, {Key: "a", Value: "b"}},
			want: -1,
		},
		{
			mda:  Metadata{{Key: "a", Value: "b"}, {Key: "c", Value: "d"}},
			mdb:  Metadata{{Key: "a", Value: "b"}, {Key: "c", Value: "e"}},
			want: -1,
		},
		{
			mda:  Metadata{{Key: "a", Value: "b"}, {Key: "c", Value: "e"}},
			mdb:  Metadata{{Key: "c", Value: "d"}, {Key: "a", Value: "b"}},
			want: 1,
		},
		{
			mda:  Metadata{{Key: "a", Value: "a"}, {Key: "b", Value: "b"}, {Key: "x", Value: "x"}},
			mdb:  Metadata{{Key: "a", Value: "a"}, {Key: "b", Value: "b"}, {Key: "c", Value: "c"}, {Key: "d", Value: "d"}},
			want: 1,
		},
	} {
		got := tc.mda.Compare(tc.mdb)
		if got != tc.want {
			assert.Failed(t, fmt.Sprintf(
				"Compare(%#v, %#v): want %d, got %d", tc.mda, tc.mdb, got, tc.want,
			))
		}
	}
}

func TestMetadataZeroZero(t *testing.T) {
	m1 := Metadata{}
	m2 := Metadata{}

	assert.True(t, m1.Equals(m2))
	assert.True(t, m2.Equals(m1))
}

func TestMetadataOutOfOrder(t *testing.T) {
	ma := Metadatum{"Key1", "Value1"}
	mb := Metadatum{"Key2", "Value2"}
	m1 := Metadata{ma, mb}
	m2 := Metadata{mb, ma}

	assert.True(t, m1.Equals(m2))
	assert.True(t, m2.Equals(m1))
}

func TestMetadataExtraElement(t *testing.T) {
	ma := Metadatum{"Key1", "Value1"}
	mb := Metadatum{"Key2", "Value2"}
	mc := Metadatum{"Key3", "Value3"}
	m1 := Metadata{ma, mb, mc}
	m2 := Metadata{mb, ma}

	assert.False(t, m1.Equals(m2))
	assert.False(t, m2.Equals(m1))
}

func TestMetadataFromMap(t *testing.T) {
	m := map[string]string{
		"foo": "bar",
		"baz": "quix",
	}

	meta := MetadataFromMap(m)
	assert.HasSameElements(t, meta, Metadata{{"foo", "bar"}, {"baz", "quix"}})
}

func getTestMD() Metadata {
	return Metadata{
		{"key1", "value1"},
		{"key2", "value2"},
	}
}

func mdPass(md Metadatum) *ValidationError { return nil }

func TestMetadataIsValid(t *testing.T) {
	md := getTestMD()
	assert.Nil(t, MetadataValid("meta", md, mdPass))
}

func TestMetadataIsValidNil(t *testing.T) {
	var md Metadata = nil
	assert.Nil(t, MetadataValid("meta", md, mdPass))
}

func TestMetadataIsValidEmpty(t *testing.T) {
	md := Metadata{}
	assert.Nil(t, MetadataValid("meta", md, mdPass))
}

func TestMetadataIsValidDupes(t *testing.T) {
	md := getTestMD()
	md = append(md, md[0])
	assert.DeepEqual(t, MetadataValid("meta", md, mdPass), &ValidationError{[]ErrorCase{
		{"meta", "duplicate meta key 'key1'"},
	}})
}

func TestMetadataIsValidFailCheck(t *testing.T) {
	md := getTestMD()
	assert.DeepEqual(
		t,
		MetadataValid("meta", md, func(d Metadatum) *ValidationError {
			if d.Key == "key2" {
				return &ValidationError{[]ErrorCase{{"whee", "whoo"}}}
			}
			return nil
		}),
		&ValidationError{[]ErrorCase{
			{"meta[key2].whee", "whoo"},
		}},
	)
}
