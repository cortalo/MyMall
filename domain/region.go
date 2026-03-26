package domain

import (
	"MyMall/domain/shared"
	"slices"
	"time"
)

const (
	RegionStatusValid int8 = iota
	RegionStatusSuspended
	RegionStatusAbandoned
)

var regionAllowedTransitions = map[int8][]int8{
	RegionStatusValid:     {RegionStatusSuspended, RegionStatusAbandoned},
	RegionStatusSuspended: {RegionStatusValid, RegionStatusAbandoned},
}

type Region struct {
	ID           int64
	CreatorID    int64
	CreatorName  string
	ModifierID   int64
	ModifierName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ParentID     int64
	Level        int8
	AreaCode     string
	ZipCode      string
	CityCode     string
	Name         string
	ShortName    string
	MergerName   string
	Pinyin       string
	Lng          float64
	Lat          float64
	Status       int8
}

func (r *Region) abandon() error {
	return nil
}

func (r *Region) createSubRegion(region *Region, operator shared.Operator) (*Region, error) {
	if !(r.Status == RegionStatusValid || r.Status == RegionStatusSuspended) {
		return nil, ErrRegionAbandoned
	}
	region.Status = r.Status
	region.Level = r.Level + 1
	region.ParentID = r.ID
	region.CreatorID = operator.ID
	region.CreatorName = operator.Username
	return region, nil
}

func (r *Region) allowTransitionTo(status int8) bool {
	return slices.Contains(regionAllowedTransitions[r.Status], status)
}

func (r *Region) changeStatus(status int8, operator shared.Operator) (*Region, error) {
	if !r.allowTransitionTo(status) {
		return nil, ErrRegionNotAllowedStatus
	}
	r.Status = status
	r.ModifierID = operator.ID
	r.ModifierName = operator.Username
	return r, nil
}
