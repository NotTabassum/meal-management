package domain

import (
	"meal-management/pkg/models"
	"meal-management/pkg/types"
)

type IMealActivityRepo interface {
	GetWeekend(deptID int) (models.Department, error)
	GetEmployeeByEmployeeID(EmployeeID uint) (models.Employee, error)
	FindAllEmployees() ([]models.Employee, error)
	FindMealActivity(date string, employeeId uint, mealType int) (*models.MealActivity, error)
	GetMealActivity(startDate, endDate string) ([]models.MealActivity, error)
	CreateMealActivity(activity *models.MealActivity) error
	UpdateMealActivity(mealActivity *models.MealActivity) error
	GetOwnMealActivity(ID uint, startDate, endDate string) ([]models.MealActivity, error)
	//FindMealADay(date string, mealType int) ([]models.MealActivity, error)
	FindPenaltyAMonth(startDate string, endDate string, employeeID uint) ([]models.MealActivity, error)
	TotalEmployees() ([]types.Employee, error)
	GetEmployeeMealCounts(startDate, endDate string) ([]types.MealSummaryResponse, error)
	GetTotalMealCounts(startDate, endDate string) (types.TotalMealCounts, error)
	GetTotalExtraMealCounts(startDate, endDate string) (int64, error)
	TotalMealADayGroup(startDate, endDate string, mealType int) ([]types.TotalMealGroupResponse, error)
	LunchToday(date string) ([]types.Employee, error)
	SnackToday(date string) ([]types.Employee, error)
	MealSummaryForGraph(startDate, endDate string) ([]models.MealActivity, error)
	ExtraMealSummaryForGraph(startDate, endDate string) ([]models.ExtraMeal, error)
	MealSummaryForMonthData(startDate string, endDate string, id uint) ([]models.MealActivity, error)
	UpdateMealStatusOff(date string) error
	CheckHoliday(date string) (bool, error)
}

type IMealActivityService interface {
	GenerateMealActivities() error
	GetMealActivityById(date string, mealType int, employeeId uint) (*models.MealActivity, error)
	GetMealActivity(startDate string, days int) ([]types.MealActivityResponse, error)
	UpdateMealActivity(mealActivity *models.MealActivity) error
	GetOwnMealActivity(ID uint, startDate string, days int) ([]types.MealActivityResponse, error)
	//TotalMealADay(date string, mealType int) (int, error)
	TotalPenaltyAMonth(date string, employeeID uint, days int) (float64, error)
	TotalMealAMonth(date string, days int) ([]types.MealSummaryResponse, error)
	TotalMealPerPerson(date string, days int, employeeID uint) (int, error)
	TotalMealCount(date string, days int) (types.TotalMealCounts, error)
	TotalMealADayGroup(date string, mealType int, days int) ([]types.TotalMealGroupResponse, error)
	LunchSummaryForEmail() error
	SnackSummaryForEmail() error
	MealSummaryForGraph(month int) ([]types.MealSummaryForGraph, error)
	MonthData(monthCount int, id uint) ([]types.MonthData, error)
	LunchToday() (string, error)
	SnackToday() (string, error)
	UpdateMealStatusForHolidays(holidayDates []string) error
	//RegularCount(date string) (int, int, error)
}
