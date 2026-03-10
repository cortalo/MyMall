package domain

import (
	"MyMall/domain/shared"
	"time"
)

// Region status
const (
	RegionStatusValid     int8 = 0
	RegionStatusSuspended int8 = 1
	RegionStatusAbandoned int8 = 2
)

//// Region tree node special IDs
//const (
//	RegionTopID   int64 = 0
//	RegionRootPID int64 = -1
//)

type Region struct {
	Id           int64
	CreatorId    int64
	CreatorName  string
	ModifierId   int64
	ModifierName string
	GmtCreate    time.Time
	GmtModified  time.Time
	Pid          int64
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

func (o *Region) abandon() error {
	return nil
}

func (o *Region) createSubRegion(region *Region, operator shared.Operator) (*Region, error) {
	if !(RegionStatusValid == o.Status || RegionStatusSuspended == o.Status) {
		return nil, ErrRegionAbandoned
	}
	region.Status = o.Status
	region.Level = o.Level + 1
	region.Pid = o.Id
	region.CreatorId = operator.Id
	region.CreatorName = operator.Username
	return region, nil
}
