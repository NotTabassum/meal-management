package domain

import (
	"meal-management/pkg/models"
	"meal-management/pkg/types"
)

type IMealPlanRepo interface {
	CreateMealPlan(mealPlan *models.MealPlan) error
	GetMealPlanByPrimaryKey(Date string, MealType string) (*models.MealPlan, error)
	GetMealPlan(startDate, endDate string) []models.MealPlan
	UpdateMealPlan(mealPlan *models.MealPlan) error
	DeleteMealPlan(date string, mealType string) error
}

type IMealPlanService interface {
	CreateMealPlan(mealPlan *models.MealPlan) error
	GetMealPlanByPrimaryKey(Date string, MealType string) (models.MealPlan, error)
	GetMealPlan(startDate string, days int) ([]types.GetMealPlanResponse, error)
	UpdateMealPlan(mealPlan *models.MealPlan) error
	DeleteMealPlan(date string, mealType string) error
}
