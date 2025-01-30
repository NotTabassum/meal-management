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
	FindMealADay(date string, mealType int) ([]models.MealActivity, error)
	FindPenaltyAMonth(startDate string, endDate string, employeeID uint) ([]models.MealActivity, error)
	TotalEmployees() ([]types.Employee, error)
}

type IMealActivityService interface {
	GenerateMealActivities() error
	GetMealActivityById(date string, mealType int, employeeId uint) (*models.MealActivity, error)
	GetMealActivity(startDate string, days int) ([]types.MealActivityResponse, error)
	UpdateMealActivity(mealActivity *models.MealActivity) error
	GetOwnMealActivity(ID uint, startDate string, days int) ([]types.MealActivityResponse, error)
	TotalMealADay(date string, mealType int) (int, error)
	TotalPenaltyAMonth(date string, employeeID uint, days int) (int, error)
	TotalMealAMonth(date string, days int) ([]types.MealSummaryResponse, error)
	TotalMealPerPerson(date string, days int, employeeID uint) (int, error)
}
