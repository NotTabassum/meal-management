package services

import (
	"errors"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
)

type LoginService struct {
	repo domain.ILoginRepo
}

func LoginServiceInstance(login domain.ILoginRepo) domain.ILoginService {
	return &LoginService{
		repo: login,
	}
}

func (service *LoginService) Login(Auth models.Login) (string, error) {
	employee, err := service.repo.Login(Auth.Email)
	if err != nil {
		return "", err
	}
	if employee.Password != Auth.Password {
		return "", errors.New("Invalid Email or Password")
	}
	token, err := domain.GenerateJWT(&employee)
	if err != nil {
		return "", err
	}
	return token, nil
}
