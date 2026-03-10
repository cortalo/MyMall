package infra

import (
	"MyMall/domain"

	"github.com/jmoiron/sqlx"
)

type RegionRepoSQL struct {
	db *sqlx.DB
}

func NewRegionRepoSQL(db *sqlx.DB) *RegionRepoSQL {
	return &RegionRepoSQL{db: db}
}

func (r *RegionRepoSQL) FindById(id int64) (*domain.Region, error) {
	sqlstr := `SELECT * FROM region_region WHERE id = ?`
	regionPo := RegionPo{}
	if err := r.db.Get(&regionPo, sqlstr, id); err != nil {
		return nil, err
	}
	return toDomain(&regionPo), nil
}
