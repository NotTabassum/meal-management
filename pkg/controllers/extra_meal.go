package controllers

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/consts"
	"meal-management/pkg/domain"
	"meal-management/pkg/middleware"
	"meal-management/pkg/models"
	"net/http"
	"time"
)

var ExtraMealService domain.IExtraMealService

func SetExtraMealService(extraMealService domain.IExtraMealService) {
	ExtraMealService = extraMealService
}

func CreateExtraMeal(e echo.Context) error {
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

	err = ExtraMealService.GenerateExtraMeal()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}
	return e.JSON(http.StatusCreated, map[string]string{"res": "New Extra Meal Activity Created"})
}

func UpdateExtraMeal(e echo.Context) error {
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
		return e.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
	}

	reqExtraMeal := &models.ExtraMeal{}

	if err := e.Bind(reqExtraMeal); err != nil {
		return e.JSON(http.StatusUnprocessableEntity, map[string]string{"res": "invalid request"})
	}

	if reqExtraMeal.Date == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"res": "Employee ID, Date and Meal Type are required"})
	}
	date := reqExtraMeal.Date
	LunchCount := reqExtraMeal.LunchCount
	SnackCount := reqExtraMeal.SnackCount

	requestedDate, err := time.Parse(consts.DateFormat, date)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format, use YYYY-MM-DD"})
	}

	now := time.Now()
	if requestedDate.Year() == now.Year() && requestedDate.YearDay() < now.YearDay() {
		return e.JSON(http.StatusForbidden, map[string]string{"error": "You cant change previous meal activity"})
	}

	if err := ExtraMealService.UpdateExtraMeal(date, LunchCount, SnackCount); err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}
	return e.JSON(http.StatusCreated, map[string]string{"res": "Updated Extra Meal Activity."})
}

func FetchExtraMeal(e echo.Context) error {
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
		return e.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
	}
	date := e.QueryParam("date")
	extraMeal, err := ExtraMealService.FetchExtraMeal(date)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}
	return e.JSON(http.StatusOK, extraMeal)

}
