package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"meal-management/pkg/domain"
	"meal-management/pkg/middleware"
	"meal-management/pkg/models"
	"net/http"
)

var HolidayService domain.IHolidayService

func SetHolidayService(holidayService domain.IHolidayService) {
	HolidayService = holidayService
}

func CreateHoliday(e echo.Context) error {
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
	var reqHoliday []models.Holiday
	if err := e.Bind(&reqHoliday); err != nil {
		fmt.Println(err)
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}

	holidates, err := HolidayService.CreateHoliday(reqHoliday)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	if len(holidates) > 0 {
		go func() {
			err := MealActivityService.UpdateMealStatusForHolidays(holidates)
			if err != nil {
				fmt.Println("Error in updating meal status:", err)
			}
		}()
	}
	return e.JSON(http.StatusCreated, "New Holiday is created successfully")
}
