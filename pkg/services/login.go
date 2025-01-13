package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
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

//	func VerifyPassword(hashedPassword, password string) bool {
//		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
//		return err == nil
//	}
func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (service *LoginService) Login(Auth models.Login) (string, error) {
	employee, err := service.repo.Login(Auth.Email)
	if err != nil {
		return "", err
	}
	//fmt.Println(employee.Password, Auth.Password)
	//if VerifyPassword(employee.Password, Auth.Password) == false {
	//	return "", errors.New("invalid Email or Password")
	//}
	if employee.Password != Auth.Password {
		return "", errors.New("invalid password")
	}

	token, err := domain.GenerateJWT(&employee)
	if err != nil {
		return "", err
	}
	return token, nil
}
