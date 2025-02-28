package domain

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"meal-management/pkg/models"
	"time"
)

type ILoginRepo interface {
	Login(Email string) (models.Employee, error)
	LoginPhone(Phone string) (models.Employee, error)
}

type ILoginService interface {
	Login(Auth models.Login) (string, error)
	LoginPhone(Auth models.Login) (string, error)
}

var secretKey = []byte("jwtserversidesecret")

func GenerateJWT(employee *models.Employee) (string, error) {
	claims := jwt.MapClaims{
		"employee_id": employee.EmployeeId,
		"email":       employee.Email,
		"is_admin":    employee.IsAdmin,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
		//"exp": time.Now().Add(time.Minute * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("could not sign the token: %w", err)
	}
	return tokenString, nil
}

//func GenerateResetToken(employee models.Employee) (string, error) {
//	claims := jwt.MapClaims{
//		"employee_id": employee.EmployeeId,
//		"email":       employee.Email,
//		"is_admin":    employee.IsAdmin,
//		"exp":     time.Now().Add(15 * time.Minute).Unix(), // Token expires in 15 minutes
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	tokenString, err := token.SignedString(secretKey)
//	if err != nil {
//		return "", fmt.Errorf("could not sign the token: %w", err)
//	}
//	return tokenString, nil
//}
