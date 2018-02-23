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

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/turbinelabs/api"
	"github.com/turbinelabs/test/assert"
)

func TestRefNameMapKeys(t *testing.T) {
	zNameRef := NewZoneNameZoneRef("that-zone-name")
	pNameRef := NewProxyNameProxyRef("that-proxy-name", zNameRef)
	got := pNameRef.MapKey()
	want := `{"proxy_name":"that-proxy-name","zone_name":"that-zone-name"}`
	assert.Equal(t, got, want)
}

func TestRefFromMapKey(t *testing.T) {
	keyStr := `{"proxy_name":"that-proxy-name","zone_name":"that-zone-name"}`
	zNameRef := NewZoneNameZoneRef("that-zone-name")
	want := NewProxyNameProxyRef("that-proxy-name", zNameRef)

	got, err := NewProxyRefFromMapKey(keyStr)
	assert.Nil(t, err)
	assert.Equal(t, got.Name(), want.Name())
	assert.Equal(t, got.ZoneRef().Name(), want.ZoneRef().Name())
}

func TestRefMapKey(t *testing.T) {
	z := api.Zone{Name: "that-zone-name"}
	p := api.Proxy{Name: "that-proxy-name"}
	pRef := NewProxyRef(p, z)

	got := pRef.MapKey()
	want := `{"proxy_name":"that-proxy-name","zone_name":"that-zone-name"}`
	assert.Equal(t, got, want)
}

func TestNewProxyRef(t *testing.T) {
	p := api.Proxy{Name: "that-proxy-name"}
	z := api.Zone{Name: "that-zone-name"}

	pRef := NewProxyRef(p, z)
	gotProxy, gotErr := pRef.Get(nil)

	assert.Nil(t, gotErr)
	assert.True(t, p.Equals(gotProxy))
}

func TestNewProxyNameProxyRef(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	zRef := NewMockZoneRef(ctrl)
	svc := NewMockAll(ctrl)
	proxySvc := NewMockProxy(ctrl)

	z := api.Zone{ZoneKey: "that-zone-key"}
	p := api.Proxy{ProxyKey: "that-proxy-key", Name: "that-proxy-name"}

	zRef.EXPECT().Get(svc).Return(z, nil)
	svc.EXPECT().Proxy().Return(proxySvc)
	proxySvc.EXPECT().Index(
		ProxyFilter{ZoneKey: z.ZoneKey, Name: p.Name},
	).Return(api.Proxies{p}, nil)

	pRef := NewProxyNameProxyRef(p.Name, zRef)
	gotProxy, gotErr := pRef.Get(svc)

	assert.Nil(t, gotErr)
	assert.True(t, p.Equals(gotProxy))

	// prove lookup is memoized
	gotProxy, gotErr = pRef.Get(svc)

	assert.Nil(t, gotErr)
	assert.True(t, p.Equals(gotProxy))
}

func TestNewProxyNameProxyRefEmptyName(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	zRef := NewMockZoneRef(ctrl)
	svc := NewMockAll(ctrl)

	pRef := NewProxyNameProxyRef("", zRef)
	gotProxy, gotErr := pRef.Get(svc)

	assert.ErrorContains(t, gotErr, "proxyName must be non-empty")
	assert.True(t, api.Proxy{}.Equals(gotProxy))
}

func TestNewProxyNameProxyRefZoneGetFails(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	zRef := NewMockZoneRef(ctrl)
	svc := NewMockAll(ctrl)

	err := errors.New("boom")

	zRef.EXPECT().Get(svc).Return(api.Zone{}, err)

	pRef := NewProxyNameProxyRef("whatever", zRef)
	gotProxy, gotErr := pRef.Get(svc)

	assert.ErrorContains(t, gotErr, "boom")
	assert.True(t, api.Proxy{}.Equals(gotProxy))
}

func TestNewProxyNameProxyRefIndexFails(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	zRef := NewMockZoneRef(ctrl)
	svc := NewMockAll(ctrl)
	proxySvc := NewMockProxy(ctrl)

	z := api.Zone{ZoneKey: "that-zone-key"}
	p := api.Proxy{ProxyKey: "that-proxy-key", Name: "that-proxy-name"}

	err := errors.New("boom")

	zRef.EXPECT().Get(svc).Return(z, nil)
	svc.EXPECT().Proxy().Return(proxySvc)
	proxySvc.EXPECT().Index(
		ProxyFilter{ZoneKey: z.ZoneKey, Name: p.Name},
	).Return(nil, err)

	pRef := NewProxyNameProxyRef(p.Name, zRef)
	gotProxy, gotErr := pRef.Get(svc)

	assert.ErrorContains(t, gotErr, "boom")
	assert.True(t, api.Proxy{}.Equals(gotProxy))
}

func TestNewProxyNameProxyRefEmptyResult(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	zRef := NewMockZoneRef(ctrl)
	svc := NewMockAll(ctrl)
	proxySvc := NewMockProxy(ctrl)

	z := api.Zone{ZoneKey: "that-zone-key"}
	p := api.Proxy{ProxyKey: "that-proxy-key", Name: "that-proxy-name"}

	zRef.EXPECT().Get(svc).Return(z, nil)
	svc.EXPECT().Proxy().Return(proxySvc)
	proxySvc.EXPECT().Index(ProxyFilter{ZoneKey: z.ZoneKey, Name: p.Name}).Return(nil, nil)

	pRef := NewProxyNameProxyRef(p.Name, zRef)
	gotProxy, gotErr := pRef.Get(svc)

	assert.ErrorContains(t, gotErr, `no Proxy found for name "that-proxy-name"`)
	assert.True(t, api.Proxy{}.Equals(gotProxy))
}

func TestNewZoneRef(t *testing.T) {
	z := api.Zone{Name: "that-zone-name"}

	zRef := NewZoneRef(z)
	gotZone, gotErr := zRef.Get(nil)

	assert.Nil(t, gotErr)
	assert.True(t, z.Equals(gotZone))
}

func TestNewZoneNameZoneRef(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	svc := NewMockAll(ctrl)
	zoneSvc := NewMockZone(ctrl)

	z := api.Zone{ZoneKey: "that-zone-key", Name: "that-zone-name"}

	svc.EXPECT().Zone().Return(zoneSvc)
	zoneSvc.EXPECT().Index(ZoneFilter{Name: z.Name}).Return(api.Zones{z}, nil)

	zRef := NewZoneNameZoneRef(z.Name)
	gotZone, gotErr := zRef.Get(svc)

	assert.Nil(t, gotErr)
	assert.True(t, z.Equals(gotZone))

	// prove zone is memoized
	gotZone, gotErr = zRef.Get(svc)

	assert.Nil(t, gotErr)
	assert.True(t, z.Equals(gotZone))
}

func TestNewZoneNameZoneRefEmptyName(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	svc := NewMockAll(ctrl)

	zRef := NewZoneNameZoneRef("")
	gotZone, gotErr := zRef.Get(svc)

	assert.ErrorContains(t, gotErr, "zoneName must be non-empty")
	assert.True(t, api.Zone{}.Equals(gotZone))
}

func TestNewZoneNameZoneRefIndexErr(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	svc := NewMockAll(ctrl)
	zoneSvc := NewMockZone(ctrl)

	z := api.Zone{ZoneKey: "that-zone-key", Name: "that-zone-name"}
	err := errors.New("boom")

	svc.EXPECT().Zone().Return(zoneSvc)
	zoneSvc.EXPECT().Index(ZoneFilter{Name: z.Name}).Return(nil, err)

	zRef := NewZoneNameZoneRef(z.Name)
	gotZone, gotErr := zRef.Get(svc)

	assert.ErrorContains(t, gotErr, "boom")
	assert.True(t, api.Zone{}.Equals(gotZone))
}

func TestNewZoneNameZoneRefEmptyResult(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	svc := NewMockAll(ctrl)
	zoneSvc := NewMockZone(ctrl)

	z := api.Zone{ZoneKey: "that-zone-key", Name: "that-zone-name"}

	svc.EXPECT().Zone().Return(zoneSvc)
	zoneSvc.EXPECT().Index(ZoneFilter{Name: z.Name}).Return(nil, nil)

	zRef := NewZoneNameZoneRef(z.Name)
	gotZone, gotErr := zRef.Get(svc)

	assert.ErrorContains(t, gotErr, `no Zone found for name "that-zone-name"`)
	assert.True(t, api.Zone{}.Equals(gotZone))
}
