package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"log"
	"meal-management/pkg/domain"
	"meal-management/pkg/middleware"
	"meal-management/pkg/models"
	"meal-management/pkg/types"
	"net/http"
	"strconv"
	"time"
)

var MealActivityService domain.IMealActivityService

func SetMealActivityService(mealActivityService domain.IMealActivityService) {
	MealActivityService = mealActivityService
}

func StartCronJob() {
	c := cron.New()
	_, err := c.AddFunc("@hourly", func() { //"*/1 * * * *" for every minute
		log.Println("Running GenerateMealActivities at:", time.Now())
		if err := MealActivityService.GenerateMealActivities(); err != nil {
			log.Printf("Error generating meal activities: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Failed to schedule cron job: %v", err)
	}
	c.Start()
	log.Println("Cron job started")
}

func CreateMealActivity(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	_, isAdmin, _ := middleware.ParseJWT(authorizationHeader)
	if !isAdmin {
		return e.JSON(http.StatusForbidden, map[string]string{"res": "Unauthorized"})
	}

	err := MealActivityService.GenerateMealActivities()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}
	return e.JSON(http.StatusCreated, map[string]string{"res": "New Meal Activity Created"})
}

func GetMealActivity(e echo.Context) error {
	stDate := e.QueryParam("start")
	tempDays := e.QueryParam("days")
	days, err := strconv.Atoi(tempDays)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Invalid input"})
	}

	mealActivity, err := MealActivityService.GetMealActivity(stDate, days)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}
	return e.JSON(http.StatusOK, mealActivity)
}

func UpdateMealActivity(e echo.Context) error {
	reqMealActivity := &types.MealActivityRequest{}

	if err := e.Bind(reqMealActivity); err != nil {
		return e.JSON(http.StatusUnprocessableEntity, map[string]string{"res": "invalid request"})
	}

	if reqMealActivity.Date == "" || reqMealActivity.MealType == 0 || reqMealActivity.EmployeeId == 0 {
		return e.JSON(http.StatusBadRequest, map[string]string{"res": "Employee ID, Date and Meal Type are required"})
	}

	date := reqMealActivity.Date
	mealType := reqMealActivity.MealType
	employeeId := reqMealActivity.EmployeeId

	existingActivity, err := MealActivityService.GetMealActivityById(string(date), int(mealType), uint(employeeId))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}

	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	ID, isAdmin, _ := middleware.ParseJWT(authorizationHeader)
	if !isAdmin {
		reqMealActivity.Penalty = existingActivity.Penalty
	}
	NewID, err := strconv.ParseUint(ID, 10, 32)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}
	if !isAdmin && uint(NewID) != employeeId {
		return e.JSON(http.StatusBadRequest, "You cannot change others activity")
	}

	updatedActivity := &models.MealActivity{
		EmployeeId:   uint(employeeId),
		Date:         date,
		MealType:     mealType,
		EmployeeName: existingActivity.EmployeeName,
		Status:       reqMealActivity.Status,
		GuestCount:   reqMealActivity.GuestCount,
		Penalty:      reqMealActivity.Penalty,
		IsOffDay:     &reqMealActivity.IsOffDay,
	}

	if err := MealActivityService.UpdateMealActivity(updatedActivity); err != nil {
		return e.JSON(http.StatusInternalServerError, err)
	}
	return e.JSON(http.StatusCreated, "Meal Activity is updated successfully")
}

func GetOwnMealActivity(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	ID, _, _ := middleware.ParseJWT(authorizationHeader)
	id, err := strconv.ParseUint(ID, 10, 32)

	stDate := e.QueryParam("start")
	tempDays := e.QueryParam("days")
	days, err := strconv.Atoi(tempDays)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Invalid input"})
	}

	mealActivity, err := MealActivityService.GetOwnMealActivity(uint(id), stDate, days)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}
	return e.JSON(http.StatusOK, mealActivity)

}

func TotalMealADay(e echo.Context) error {
	reqMealActivity := &types.MealActivityRequest{}

	if err := e.Bind(reqMealActivity); err != nil {
		return e.JSON(http.StatusUnprocessableEntity, map[string]string{"res": "invalid request"})
	}
	date := reqMealActivity.Date
	mealType := reqMealActivity.MealType
	mealCount, err := MealActivityService.TotalMealADay(date, mealType)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}
	return e.JSON(http.StatusCreated, mealCount)
}

func TotalPenalty(e echo.Context) error {
	reqMealActivity := &types.PenaltyRequest{}

	if err := e.Bind(reqMealActivity); err != nil {
		return e.JSON(http.StatusUnprocessableEntity, map[string]string{"res": "invalid request"})
	}
	date := reqMealActivity.Date
	employeeID := reqMealActivity.EmployeeId
	days := reqMealActivity.Days
	mealCount, err := MealActivityService.TotalPenaltyAMonth(date, employeeID, days)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}
	return e.JSON(http.StatusCreated, mealCount)
}

func UpdateGroupMealActivity(e echo.Context) error {
	var groupMeal []types.MealActivityRequest
	if err := e.Bind(&groupMeal); err != nil {
		fmt.Println(err)
		return e.JSON(http.StatusUnprocessableEntity, map[string]string{"res": "invalid request"})
	}
	for _, val := range groupMeal {
		if val.Date == "" || val.MealType == 0 || val.EmployeeId == 0 {
			return e.JSON(http.StatusBadRequest, map[string]string{"res": "Employee ID, Date and Meal Type are required"})
		}

		date := val.Date
		mealType := val.MealType
		employeeId := val.EmployeeId

		existingActivity, err := MealActivityService.GetMealActivityById(string(date), int(mealType), uint(employeeId))
		if err != nil {
			return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
		}

		authorizationHeader := e.Request().Header.Get("Authorization")
		if authorizationHeader == "" {
			return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
		}
		ID, isAdmin, _ := middleware.ParseJWT(authorizationHeader)
		if !isAdmin {
			val.Penalty = existingActivity.Penalty
		}
		NewID, err := strconv.ParseUint(ID, 10, 32)
		if err != nil {
			return e.JSON(http.StatusBadRequest, err.Error())
		}
		if !isAdmin && uint(NewID) != employeeId {
			return e.JSON(http.StatusBadRequest, "You cannot change others activity")
		}

		updatedActivity := &models.MealActivity{
			EmployeeId:   uint(employeeId),
			Date:         date,
			MealType:     mealType,
			EmployeeName: existingActivity.EmployeeName,
			Status:       val.Status,
			GuestCount:   val.GuestCount,
			Penalty:      val.Penalty,
			IsOffDay:     &val.IsOffDay,
		}

		if err := MealActivityService.UpdateMealActivity(updatedActivity); err != nil {
			return e.JSON(http.StatusInternalServerError, err)
		}
	}
	return e.JSON(http.StatusCreated, "Meal Activity is updated successfully")
}
