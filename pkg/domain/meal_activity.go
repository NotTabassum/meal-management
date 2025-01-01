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
}

type IMealActivityService interface {
	GenerateMealActivities() error
	GetMealActivityById(date string, mealType int, employeeId uint) (*models.MealActivity, error)
	GetMealActivity(startDate string, days int) ([]types.MealActivityResponse, error)
	UpdateMealActivity(mealActivity *models.MealActivity) error
}
