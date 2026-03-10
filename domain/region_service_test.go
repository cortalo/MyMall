package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type mockRegionRepo struct {
	findById func(id int64) (*Region, error)
}

func (m *mockRegionRepo) FindById(id int64) (*Region, error) {
	return m.findById(id)
}

func TestUserService_FindById(t *testing.T) {
	repo := &mockRegionRepo{
		findById: func(id int64) (*Region, error) {
			return &Region{
				Id:   id,
				Name: "China",
			}, nil
		},
	}
	userService := NewUserService(repo)
	result, err := userService.FindById(1)
	require.NoError(t, err)
	require.Equal(t, int64(1), result.Id)
	require.Equal(t, "China", result.Name)
}
