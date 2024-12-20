package controllers

import (
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
	auth := models.Login{
		Email:    reqLogin.Email,
		Password: reqLogin.Password,
	}

	login, err := LoginService.Login(auth)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, login)
}
