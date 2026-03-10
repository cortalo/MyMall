package infra_test

import (
	"MyMall/infra"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const dsn = "root:root@tcp(localhost:3306)/oomall_test?charset=utf8mb4&parseTime=True"

func TestRegionRepoSQL_FindById(t *testing.T) {

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	defer func() {
		if err := zap.L().Sync(); err != nil {
			t.Logf("zap sync error: %v", err)
		}
	}()

	db, err := sqlx.Connect("mysql", dsn)
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("db close error: %v", err)
		}
	}()

	// Clean
	_, err = db.Exec("DELETE FROM region_region")
	require.NoError(t, err)

	// Insert
	regionPo := infra.RegionPo{
		Pid:          -1,
		Level:        1,
		AreaCode:     "100000",
		ZipCode:      "100000",
		CityCode:     "010",
		Name:         "China",
		ShortName:    "CN",
		MergerName:   "China",
		Pinyin:       "China",
		Lng:          116.4074,
		Lat:          39.9042,
		Status:       0,
		CreatorId:    1,
		CreatorName:  "admin",
		ModifierId:   1,
		ModifierName: "admin",
		GmtCreate:    time.Now(),
		GmtModified:  time.Now(),
	}
	result, err := db.NamedExec(`
		INSERT INTO region_region 
			(pid, level, area_code, zip_code, city_code, name, short_name, merger_name, pinyin, lng, lat, status, creator_id, creator_name, modifier_id, modifier_name, gmt_create, gmt_modified)
		VALUES 
			(:pid, :level, :area_code, :zip_code, :city_code, :name, :short_name, :merger_name, :pinyin, :lng, :lat, :status, :creator_id, :creator_name, :modifier_id, :modifier_name, :gmt_create, :gmt_modified)
	`, regionPo)
	require.NoError(t, err)
	insertedId, err := result.LastInsertId()
	require.NoError(t, err)

	repo := infra.NewRegionRepoSQL(db)
	region, err := repo.FindById(insertedId)
	require.NoError(t, err)
	require.Equal(t, insertedId, region.Id)
	require.Equal(t, "China", region.Name)

}
