package models

type MealPlan struct {
	Date     string `gorm:"primaryKey"`
	MealType int    `gorm:"primaryKey"`
	Food     string
}
