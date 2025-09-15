package users

type IUserRepository interface {
	GetAllUsers() ([]User, error)
	CreateNewUser(User) (User, error)
	DeleteUser(string) error
}

type UserService struct {
	repo IUserRepository
}

func NewUserService(repo IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(user User) (User, error) {
	if err := ValidateUser(user); err != nil {
		return User{}, err
	}

	newUser, err := s.repo.CreateNewUser(user)
	if err != nil {
		return user, err
	}

	return newUser, nil
}

func (s *UserService) GetUsers() ([]User, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}
