package controllers

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/domain"
	"meal-management/pkg/middleware"
	"meal-management/pkg/models"
	"net/http"
)

var PreferenceService domain.IPreferenceService

func SetPreferenceService(pService domain.IPreferenceService) {
	PreferenceService = pService
}

func CreatePreference(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	_, isAdmin, err := middleware.ParseJWT(authorizationHeader)
	if err != nil {
		if err.Error() == "token expired" {
			return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
		}
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	if !isAdmin {
		return e.JSON(http.StatusForbidden, map[string]string{"res": "Unauthorized"})
	}
	reqPref := &models.Preference{}
	if err := e.Bind(reqPref); err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}
	if err := PreferenceService.CreatePreferenceService(reqPref); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, "New Preference is created successfully")
}

func GetPreference(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	preference, err := PreferenceService.GetPreferenceService()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, "Cannot Fetch Preference")
	}
	return e.JSON(http.StatusOK, preference)
}
