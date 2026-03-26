package domain

import (
	"MyMall/domain/shared"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestRegion(status int8) *Region {
	return &Region{
		ID:       1,
		ParentID: -1,
		Level:    0,
		Status:   status,
	}
}

func newTestOperator() shared.Operator {
	return shared.Operator{
		ID:       777,
		Username: "admin",
	}
}

func TestRegionCreateSubRegion(t *testing.T) {
	parent := newTestRegion(RegionStatusValid)
	subRegion := &Region{
		Name: "subRegion",
	}
	operator := newTestOperator()
	result, err := parent.createSubRegion(subRegion, operator)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, parent.Status, result.Status)
	require.Equal(t, parent.Level+1, result.Level)
	require.Equal(t, parent.ID, result.ParentID)
	require.Equal(t, subRegion.Name, result.Name)
	require.Equal(t, operator.ID, result.CreatorID)
	require.Equal(t, operator.Username, result.CreatorName)
}

func TestRegionCreateSubRegionWithStatusAbandoned(t *testing.T) {
	parent := newTestRegion(RegionStatusAbandoned)
	subRegion := &Region{
		Name: "subRegion",
	}
	operator := newTestOperator()
	result, err := parent.createSubRegion(subRegion, operator)
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrRegionAbandoned)
}

func TestRegionChangeStatus(t *testing.T) {
	region := newTestRegion(RegionStatusValid)
	operator := newTestOperator()
	result, err := region.changeStatus(RegionStatusSuspended, operator)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, RegionStatusSuspended, result.Status)
	require.Equal(t, operator.ID, result.ModifierID)
	require.Equal(t, operator.Username, result.ModifierName)
}

func TestRegionChangeStatusNotAllowedStatus(t *testing.T) {
	region := newTestRegion(RegionStatusAbandoned)
	operator := newTestOperator()
	result, err := region.changeStatus(RegionStatusSuspended, operator)
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrRegionNotAllowedStatus)
}
