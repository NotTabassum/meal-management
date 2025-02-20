package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
	"meal-management/pkg/types"
	"net/http"
)

var LoginService domain.ILoginService

func SetLoginService(login domain.ILoginService) {
	LoginService = login
}

func Login(e echo.Context) error {
	reqLogin := &types.CreateLoginRequest{}
	if err := e.Bind(reqLogin); err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}
	if reqLogin.Email == "" && reqLogin.PhoneNumber == "" {
		return fmt.Errorf("either email or phone number must be provided")
	}
	var login string
	var err error
	if reqLogin.Email != "" {
		auth := models.Login{
			Email:    reqLogin.Email,
			Password: reqLogin.Password,
		}

		login, err = LoginService.Login(auth)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, err.Error())
		}
	} else if reqLogin.PhoneNumber != "" {
		auth := models.Login{
			PhoneNumber: reqLogin.PhoneNumber,
			Password:    reqLogin.Password,
		}

		login, err = LoginService.LoginPhone(auth)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return e.JSON(http.StatusOK, login)
}
