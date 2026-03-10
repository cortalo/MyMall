package domain

type RegionService interface {
	FindById(id int64) (*Region, error)
}

type RegionRepo interface {
	FindById(id int64) (*Region, error)
}

type regionService struct {
	repo RegionRepo
}

func NewUserService(repo RegionRepo) RegionService {
	return &regionService{repo: repo}
}

func (s *regionService) FindById(id int64) (*Region, error) {
	return s.repo.FindById(id)
}
