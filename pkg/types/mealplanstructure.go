package types

import "gorm.io/datatypes"

type GetMealPlanResponse struct {
	Date string `json:"date"`
	Menu []Menu `json:"menu"`
}

type Menu struct {
	MealType       string         `json:"meal_type"`
	Food           string         `json:"food"`
	PreferenceFood datatypes.JSON `json:"preference_food"`
}

type CreateMealPlanRequest struct {
	Date           string         `json:"date" validate:"required"`
	MealType       string         `json:"meal_type" validate:"required"`
	Food           string         `json:"food" validate:"required"`
	PreferenceFood datatypes.JSON `json:"preference_food"`
}

type GetMealPlanRequest struct {
	StartDate string `json:"start_date" validate:"required"`
	Days      int    `json:"days" validate:"required,min=1"`
}

type Meal struct {
	Food string `json:"food" validate:"required"`
}
