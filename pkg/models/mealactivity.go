package models

type MealActivity struct {
	Date         string `gorm:"primaryKey"`
	EmployeeId   uint   `gorm:"primaryKey"`
	MealType     int    `gorm:"primaryKey"`
	EmployeeName string
	Status       *bool
	GuestCount   *int
	Penalty      *bool
}
