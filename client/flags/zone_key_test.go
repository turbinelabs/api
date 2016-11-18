package flags

import (
	"errors"
	"flag"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/turbinelabs/api"
	"github.com/turbinelabs/api/service"
	"github.com/turbinelabs/test/assert"
)

func TestNewZoneKeyFromFlags(t *testing.T) {
	flagset := flag.NewFlagSet("NewZoneKeyFromFlags options", flag.PanicOnError)

	ff := NewZoneKeyFromFlags(flagset)
	ffImpl := ff.(*zoneKeyFromFlags)

	flagset.Parse([]string{"-api.zone-name=red-sector-a"})

	assert.Equal(t, ffImpl.zoneName, "red-sector-a")
}

func TestZoneKeyFromFlagsMake(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	mockService := service.NewMockAll(ctrl)
	mockZone := service.NewMockZone(ctrl)
	zoneKey := api.ZoneKey("zk-1")

	mockService.EXPECT().Zone().Return(mockZone)
	mockZone.EXPECT().
		Index(service.ZoneFilter{Name: "z1"}).
		Return(api.Zones{{ZoneKey: zoneKey}}, nil)

	ff := &zoneKeyFromFlags{"z1"}

	zk, err := ff.Get(mockService)
	assert.Nil(t, err)
	assert.Equal(t, zk, zoneKey)
}

func TestZoneKeyFromFlagsZoneName(t *testing.T) {
	var ff ZoneKeyFromFlags
	ff = &zoneKeyFromFlags{"red-sector-a"}
	assert.Equal(t, ff.ZoneName(), "red-sector-a")
}

func TestZoneKeyFromFlagsMakeGetZoneError(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	mockService := service.NewMockAll(ctrl)
	mockZone := service.NewMockZone(ctrl)

	mockService.EXPECT().Zone().Return(mockZone)
	mockZone.EXPECT().Index(service.ZoneFilter{Name: "z1"}).Return(nil, errors.New("boom"))

	ff := &zoneKeyFromFlags{"z1"}

	wantErr := errors.New("error finding Zone with name z1: boom")

	zk, err := ff.Get(mockService)
	assert.DeepEqual(t, err, wantErr)
	assert.Equal(t, zk, api.ZoneKey(""))
}

func TestZoneKeyFromFlagsMakeNoZone(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	mockService := service.NewMockAll(ctrl)
	mockZone := service.NewMockZone(ctrl)

	mockService.EXPECT().Zone().Return(mockZone)
	mockZone.EXPECT().Index(service.ZoneFilter{Name: "z1"}).Return(nil, nil)

	ff := &zoneKeyFromFlags{"z1"}

	wantErr := errors.New("Zone with name z1 does not exist")

	zk, err := ff.Get(mockService)
	assert.DeepEqual(t, err, wantErr)
	assert.Equal(t, zk, api.ZoneKey(""))
}