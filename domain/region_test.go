package domain

import (
	"MyMall/domain/shared"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegionCreateSubRegion(t *testing.T) {
	parent := &Region{
		Id:     1,
		Pid:    -1,
		Level:  0,
		Status: RegionStatusValid,
	}
	subRegion := &Region{
		Name: "subRegion",
	}
	operator := shared.Operator{
		Id:       777,
		Username: "admin",
	}
	result, err := parent.createSubRegion(subRegion, operator)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, result.Status, parent.Status)
	require.Equal(t, result.Level, parent.Level+1)
	require.Equal(t, result.Pid, parent.Id)
	require.Equal(t, result.Name, subRegion.Name)
	require.Equal(t, result.CreatorId, operator.Id)
	require.Equal(t, result.CreatorName, operator.Username)
}

func TestRegionCreateSubRegionWithStatusAbandoned(t *testing.T) {
	parent := &Region{
		Id:     1,
		Pid:    -1,
		Level:  0,
		Status: RegionStatusAbandoned,
	}
	subRegion := &Region{
		Name: "subRegion",
	}
	operator := shared.Operator{
		Id:       777,
		Username: "admin",
	}
	result, err := parent.createSubRegion(subRegion, operator)
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, ErrRegionAbandoned, err)
}
