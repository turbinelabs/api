package api

import (
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
	assert.True(t, m2.Equivalent(m1))
	assert.True(t, m1.Equivalent(m2))
}

func TestMetadatumValueChanged(t *testing.T) {
	m1 := Metadatum{"Key", "Value2"}
	m2 := Metadatum{"Key", "Value"}

	assert.False(t, m1.Equals(m2))
	assert.False(t, m2.Equals(m1))
	assert.False(t, m1.Equivalent(m2))
	assert.False(t, m2.Equivalent(m1))
}

func TestMetadatumKeyChanged(t *testing.T) {
	m1 := Metadatum{"Key", "Value"}
	m2 := Metadatum{"Key2", "Value"}

	assert.False(t, m1.Equals(m2))
	assert.False(t, m2.Equals(m1))
	assert.False(t, m1.Equivalent(m2))
	assert.False(t, m2.Equivalent(m1))
}

// Metadata
func TestMetadataZeroNil(t *testing.T) {
	m1 := Metadata{}
	var m2 Metadata

	assert.True(t, m1.Equals(m2))
	assert.True(t, m2.Equals(m1))
	assert.True(t, m1.Equivalent(m2))
	assert.True(t, m2.Equivalent(m1))
}

func TestMetadataEqualsNilNil(t *testing.T) {
	var m1, m2 Metadata
	assert.True(t, m1.Equals(m2))
	assert.True(t, m2.Equals(m1))
	assert.True(t, m1.Equivalent(m2))
	assert.True(t, m2.Equivalent(m1))
}

func TestMetadataZeroZero(t *testing.T) {
	m1 := Metadata{}
	m2 := Metadata{}

	assert.True(t, m1.Equals(m2))
	assert.True(t, m2.Equals(m1))
	assert.True(t, m1.Equivalent(m2))
	assert.True(t, m2.Equivalent(m1))
}

func TestMetadataOutOfOrder(t *testing.T) {
	ma := Metadatum{"Key1", "Value1"}
	mb := Metadatum{"Key2", "Value2"}
	m1 := Metadata{ma, mb}
	m2 := Metadata{mb, ma}

	assert.True(t, m1.Equals(m2))
	assert.True(t, m2.Equals(m1))
	assert.True(t, m1.Equivalent(m2))
	assert.True(t, m2.Equivalent(m1))
}

func TestMetadataExtraElement(t *testing.T) {
	ma := Metadatum{"Key1", "Value1"}
	mb := Metadatum{"Key2", "Value2"}
	mc := Metadatum{"Key3", "Value3"}
	m1 := Metadata{ma, mb, mc}
	m2 := Metadata{mb, ma}

	assert.False(t, m1.Equals(m2))
	assert.False(t, m2.Equals(m1))
	assert.False(t, m1.Equivalent(m2))
	assert.False(t, m2.Equivalent(m1))
}

func TestMetadataFromMap(t *testing.T) {
	m := map[string]string{
		"foo": "bar",
		"baz": "quix",
	}

	meta := MetadataFromMap(m)
	assert.HasSameElements(t, meta, Metadata{{"foo", "bar"}, {"baz", "quix"}})
}
