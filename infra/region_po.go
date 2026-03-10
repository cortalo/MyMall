package infra

import "time"

//go:generate go run ../gen/mapper_gen.go -src=RegionPo -dst=domain.Region -pkg=MyMall/domain -out=region_mapper.go
type RegionPo struct {
	Id           int64     `db:"id"`
	Pid          int64     `db:"pid"`
	Level        int8      `db:"level"`
	AreaCode     string    `db:"area_code"`
	ZipCode      string    `db:"zip_code"`
	CityCode     string    `db:"city_code"`
	Name         string    `db:"name"`
	ShortName    string    `db:"short_name"`
	MergerName   string    `db:"merger_name"`
	Pinyin       string    `db:"pinyin"`
	Lng          float64   `db:"lng"`
	Lat          float64   `db:"lat"`
	Status       int8      `db:"status"`
	CreatorId    int64     `db:"creator_id"`
	CreatorName  string    `db:"creator_name"`
	ModifierId   int64     `db:"modifier_id"`
	ModifierName string    `db:"modifier_name"`
	GmtCreate    time.Time `db:"gmt_create"`
	GmtModified  time.Time `db:"gmt_modified"`
}
