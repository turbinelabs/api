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

package api

import (
	"time"

	"github.com/turbinelabs/api/changetype"
	"github.com/turbinelabs/api/objecttype"
	tbntime "github.com/turbinelabs/nonstdlib/time"
)

// ChangeDescription combines the metadata about a change and the specific
// object fields and value.
type ChangeDescription struct {
	ChangeMeta

	// Diffs is an itemized slice of attributes that were changed. Each modified
	// field on a modified object should correspond to one ChangeEntry.
	Diffs []ChangeEntry `json:"diffs"`
}

// ChangeMeta is a collection of data about a change.
type ChangeMeta struct {
	// AtMs indicates when the change was made in milliseconds since the Unix
	// epoch.
	AtMs int64 `json:"at"`

	// Txn is an id that can be used to tie multiple object changes to a single
	// logical action.
	Txn string `json:"txn"`

	// OrgKey is the owning organization of the object(s) changed.
	OrgKey OrgKey `json:"-"`

	// ActorKey is the user key of the person making the recorded change.
	ActorKey UserKey `json:"actor_key"`

	// Comment records any information provided as part of this change.
	Comment string `json:"comment"`
}

// At returns the time that the Change was recorded at.
func (cm ChangeMeta) At() time.Time {
	if cm.AtMs == 0 {
		return time.Time{}
	}

	return tbntime.FromUnixMilli(cm.AtMs)
}

func (cm *ChangeMeta) SetAt(t time.Time) {
	cm.AtMs = tbntime.ToUnixMilli(t)
}

// ChangeEntry records a single change made on a specific object.
type ChangeEntry struct {
	// ObjectType is the kind of object that was modified.
	objecttype.ObjectType

	// ObjectKey is the key of the object that was changed.
	ObjectKey string `json:"object_key"`

	// ZoneKey is the zone containing ChangeEntry references an object in this
	// zone.  If the object changed is not bound to a specific zone the value
	// will be ZoneKey("").
	ZoneKey ZoneKey `json:"zone_key"`

	// ChangeType is the kind of change (i.e. was data added or removed) that is
	// being recorded. Each modification will have two changes recorded: a removal
	// of the old content and addition of the new content.
	changetype.ChangeType

	// Path is the path to the data that is changed.
	Path string `json:"path"`

	// Value is the value that was updated as indicated by ChangeType.
	Value string `json:"value"`
}
