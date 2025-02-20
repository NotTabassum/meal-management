package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"meal-management/pkg/consts"
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

//func StartCronJobMealActivity() {
//	c := cron.New()
//	_, err := c.AddFunc("*/1 * * * *", func() { //"*/1 * * * *" for every minute
//		log.Println("Running GenerateMealActivities at:", time.Now())
//		if err := MealActivityService.GenerateMealActivities(); err != nil {
//			log.Printf("Error generating meal activities: %v", err)
//		}
//	})
//
//	if err != nil {
//		log.Fatalf("Failed to schedule cron job: %v", err)
//	}
//	c.Start()
//	log.Println("Cron job started")
//	select {}
//}

//func StartCronJobMealActivity() {
//	log.Println("✅ Initializing CronJob for MealActivity")
//
//	if MealActivityService == nil {
//		log.Fatal("🚨 MealActivityService is nil! Check initialization.")
//	}
//
//	c := cron.New()
//
//	_, err := c.AddFunc("*/1 * * * *", func() {
//		log.Println("🕛 Running GenerateMealActivities at:", time.Now())
//
//		if MealActivityService == nil {
//			log.Println("🚨 MealActivityService is nil! Skipping job.")
//			return
//		}
//
//		if err := MealActivityService.GenerateMealActivities(); err != nil {
//			log.Printf("❌ Error generating meal activities: %v", err)
//		}
//	})
//
//	if err != nil {
//		log.Fatalf("🚨 Failed to schedule GenerateMealActivities: %v", err)
//	}
//
//	c.Start()
//	log.Println("✅ Cron jobs started successfully.")
//
//	//select {} // Keep running
//}

func CreateMealActivity(e echo.Context) error {
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

	err = MealActivityService.GenerateMealActivities()
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

	requestedDate, err := time.Parse(consts.DateFormat, date)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format, use YYYY-MM-DD"})
	}

	now := time.Now()
	if requestedDate.Year() == now.Year() && requestedDate.YearDay() < now.YearDay() {
		return e.JSON(http.StatusForbidden, map[string]string{"error": "You cant change previous meal activity"})
	} else if requestedDate.Year() == now.Year() && requestedDate.YearDay() == now.YearDay() {
		if mealType == 1 {
			cutoff := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())

			if now.After(cutoff) {
				return e.JSON(http.StatusForbidden, map[string]string{"error": "Lunch update is not allowed after 10 AM"})
			}
		} else if mealType == 2 {
			cutoff := time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, now.Location())

			if now.After(cutoff) {
				return e.JSON(http.StatusForbidden, map[string]string{"error": "Snacks update is not allowed after 2 PM"})
			}
		}
	}

	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	ID, isAdmin, err := middleware.ParseJWT(authorizationHeader)
	if err != nil {
		if err.Error() == "token expired" {
			return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
		}
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	existingActivity, err := MealActivityService.GetMealActivityById(string(date), int(mealType), uint(employeeId))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}

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
	ID, _, err := middleware.ParseJWT(authorizationHeader)
	if err != nil {
		if err.Error() == "token expired" {
			return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
		}
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
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

//
//func TotalMealADay(e echo.Context) error {
//	reqMealActivity := &types.MealActivityRequest{}
//
//	if err := e.Bind(reqMealActivity); err != nil {
//		return e.JSON(http.StatusUnprocessableEntity, map[string]string{"res": "invalid request"})
//	}
//	date := reqMealActivity.Date
//	mealType := reqMealActivity.MealType
//	mealCount, err := MealActivityService.TotalMealADay(date, mealType)
//	if err != nil {
//		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
//	}
//	return e.JSON(http.StatusCreated, mealCount)
//}

func TotalMealADayGroup(e echo.Context) error {
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

	reqMealActivity := &types.TotalMealGroupRequest{}
	if err := e.Bind(reqMealActivity); err != nil {
		return e.JSON(http.StatusUnprocessableEntity, map[string]string{"res": "invalid request"})
	}
	date := reqMealActivity.Date
	mealType := reqMealActivity.MealType
	days := reqMealActivity.Days
	groupedMealCount, err := MealActivityService.TotalMealADayGroup(date, mealType, days)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}
	return e.JSON(http.StatusCreated, groupedMealCount)
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
		ID, isAdmin, err := middleware.ParseJWT(authorizationHeader)
		if err != nil {
			if err.Error() == "token expired" {
				return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
			}
			return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
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
			PenaltyScore: &val.PenaltyScore,
		}

		if err := MealActivityService.UpdateMealActivity(updatedActivity); err != nil {
			return e.JSON(http.StatusInternalServerError, err)
		}
	}
	return e.JSON(http.StatusCreated, "Meal Activity is updated successfully")
}

func TotalMealAMonth(e echo.Context) error {
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
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Not Authorized"})
	}
	reqMealActivity := &types.MealSummaryReq{}

	if err := e.Bind(reqMealActivity); err != nil {
		return e.JSON(http.StatusUnprocessableEntity, map[string]string{"res": "invalid request"})
	}
	stdate := reqMealActivity.StartDate
	days := reqMealActivity.Days
	mealSummary, err := MealActivityService.TotalMealAMonth(stdate, days)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}
	return e.JSON(http.StatusCreated, mealSummary)
}

func TotalMealPerPerson(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	ID, _, err := middleware.ParseJWT(authorizationHeader)
	if err != nil {
		if err.Error() == "token expired" {
			return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
		}
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	id, err := strconv.ParseUint(ID, 10, 32)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}
	reqMealActivity := &types.MealSummaryReq{}

	if err := e.Bind(reqMealActivity); err != nil {
		return e.JSON(http.StatusUnprocessableEntity, map[string]string{"res": "invalid request"})
	}
	stdate := reqMealActivity.StartDate
	days := reqMealActivity.Days
	mealCount, err := MealActivityService.TotalMealPerPerson(stdate, days, uint(id))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}
	return e.JSON(http.StatusCreated, mealCount)
}

func TotalMealCount(e echo.Context) error {
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
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Not Authorized"})
	}

	reqMealActivity := &types.MealSummaryReq{}
	if err := e.Bind(reqMealActivity); err != nil {
		return e.JSON(http.StatusUnprocessableEntity, map[string]string{"res": "invalid request"})
	}
	stdate := reqMealActivity.StartDate
	days := reqMealActivity.Days
	mealCount, err := MealActivityService.TotalMealCount(stdate, days)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}
	return e.JSON(http.StatusCreated, mealCount)
}

func MealSummaryForGraph(e echo.Context) error {
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
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Not Authorized"})
	}

	monthStr := e.QueryParam("month")
	month, err := strconv.Atoi(monthStr)
	mealSummary, err := MealActivityService.MealSummaryForGraph(month)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"res": "Internal server error"})
	}
	return e.JSON(http.StatusCreated, mealSummary)
}

func TodayLunch(e echo.Context) error {
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
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Not Authorized"})
	}

	TodayLunchSummary := MealActivityService.LunchToday()
	return e.JSON(http.StatusOK, TodayLunchSummary)
}
