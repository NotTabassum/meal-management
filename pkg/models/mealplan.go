package models

type MealPlan struct {
	Date     string `gorm:"primaryKey" json:"date"`
	MealType string `gorm:"primaryKey" json:"meal_type"`
	Food     string `json:"food"`
}
