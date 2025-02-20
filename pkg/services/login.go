package services

import (
	"errors"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
	"meal-management/pkg/security"
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
	if ok := security.CheckPasswordHash(Auth.Password, employee.Password); ok == false {
		return "", errors.New("invalid Email or Password")
	}

	token, err := domain.GenerateJWT(&employee)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (service *LoginService) LoginPhone(Auth models.Login) (string, error) {
	employee, err := service.repo.LoginPhone(Auth.PhoneNumber)
	if err != nil {
		return "", err
	}
	if ok := security.CheckPasswordHash(Auth.Password, employee.Password); ok == false {
		return "", errors.New("invalid Email or Password")
	}

	token, err := domain.GenerateJWT(&employee)
	if err != nil {
		return "", err
	}
	return token, nil
}
