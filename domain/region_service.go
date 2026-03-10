package domain

type UserService interface {
	FindById(id int64) (*Region, error)
}

type RegionRepo interface {
	FindById(id int64) (*Region, error)
}

type userService struct {
	repo RegionRepo
}

func NewUserService(repo RegionRepo) UserService {
	return &userService{repo: repo}
}

func (s *userService) FindById(id int64) (*Region, error) {
	return s.repo.FindById(id)
}
