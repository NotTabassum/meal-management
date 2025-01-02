package controllers

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/domain"
	"meal-management/pkg/middleware"
	"meal-management/pkg/models"
	"meal-management/pkg/types"
	"net/http"
	"strconv"
)

var MealPlanService domain.IMealPlanService

func SetMealPlanService(mpService domain.IMealPlanService) {
	MealPlanService = mpService
}

func CreateMealPlan(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	_, isAdmin, _ := middleware.ParseJWT(authorizationHeader)
	if !isAdmin {
		return e.JSON(http.StatusForbidden, map[string]string{"res": "Unauthorized"})
	}

	reqMeal := &types.CreateMealPlanRequest{}
	if err := e.Bind(reqMeal); err != nil {
		//fmt.Println(err)
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}

	//if err := reqEmployee.Validate(); err != nil {
	//	return e.JSON(http.StatusBadRequest, err.Error())
	//}

	meal := &models.MealPlan{
		Date:     reqMeal.Date,
		MealType: reqMeal.MealType,
		Food:     reqMeal.Food,
	}
	if err := MealPlanService.CreateMealPlan(meal); err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, "MealPlan created successfully")
}

func GetMealPlanByPrimaryKey(e echo.Context) error {
	date := e.QueryParam("date")
	mealType := e.QueryParam("meal_type")

	mealPlan, err := MealPlanService.GetMealPlanByPrimaryKey(date, mealType)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, mealPlan)
}

func GetMealPlan(e echo.Context) error {
	StDate := e.QueryParam("start")
	tempDays := e.QueryParam("days")

	days, err := strconv.Atoi(tempDays)
	//fmt.Println(days)

	if err != nil || days < 1 {
		return e.JSON(http.StatusBadRequest, "Enter a valid number of days (must be 1 or more)")
	}

	mealPlan, err := MealPlanService.GetMealPlan(StDate, days)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}
	return e.JSON(http.StatusOK, mealPlan)

}
func UpdateMealPlan(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	_, isAdmin, _ := middleware.ParseJWT(authorizationHeader)
	if !isAdmin {
		return e.JSON(http.StatusForbidden, map[string]string{"res": "Unauthorized"})
	}
	reqMealPlan := &types.CreateMealPlanRequest{}

	if err := e.Bind(reqMealPlan); err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Input")
	}

	//if err := reqEmployee.Validate(); err != nil {
	//	return e.JSON(http.StatusBadRequest, err.Error())
	//}
	if reqMealPlan.Date == "" || reqMealPlan.MealType == "" || reqMealPlan.Food == "" {
		return e.JSON(http.StatusBadRequest, "All the fields are required")
	}
	tempDate := reqMealPlan.Date
	tempMealType := reqMealPlan.MealType

	meal, err := MealPlanService.GetMealPlanByPrimaryKey(tempDate, tempMealType)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	updatedMeal := &models.MealPlan{
		Date:     ifNotEmpty(reqMealPlan.Date, meal.Date),
		MealType: ifNotEmpty(reqMealPlan.MealType, meal.MealType),
		Food:     ifNotEmpty(reqMealPlan.Food, meal.Food),
	}

	if err := MealPlanService.UpdateMealPlan(updatedMeal); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, "Meal is updated successfully")
}

func DeleteMealPlan(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	_, isAdmin, _ := middleware.ParseJWT(authorizationHeader)
	if !isAdmin {
		return e.JSON(http.StatusForbidden, map[string]string{"res": "Unauthorized"})
	}
	reqMealPlan := &types.CreateMealPlanRequest{}

	if err := e.Bind(reqMealPlan); err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Input")
	}

	if reqMealPlan.Date == "" || reqMealPlan.MealType == "" {
		return e.JSON(http.StatusBadRequest, "Both date and meal type are required")
	}
	tempDate := reqMealPlan.Date
	tempMealType := reqMealPlan.MealType

	_, err := MealPlanService.GetMealPlanByPrimaryKey(tempDate, tempMealType)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	if err := MealPlanService.DeleteMealPlan(tempDate, tempMealType); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, "Meal is deleted successfully")
}
