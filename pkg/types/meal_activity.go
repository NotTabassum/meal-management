package types

import "gorm.io/datatypes"

type MealActivityRequest struct {
	Date         string   `json:"date" validate:"required"`
	EmployeeId   uint     `json:"employee_id" validate:"required"`
	MealType     int      `json:"meal_type" validate:"required"`
	Status       *bool    `json:"status"`
	GuestCount   *int     `json:"guest_count"`
	Penalty      *bool    `json:"penalty"`
	IsOffDay     bool     `json:"is_off_day"`
	PenaltyScore *float64 `gorm:"type:decimal(10,2);" json:"penalty_score"`
}

type PenaltyRequest struct {
	Date       string `json:"date"`
	EmployeeId uint   `json:"employee_id"`
	Days       int    `json:"days"`
}
type TotalMealGroupRequest struct {
	Date     string `json:"date" validate:"required"`
	Days     int    `json:"days"`
	MealType int    `json:"meal_type" validate:"required"`
}

type TotalMealGroupResponse struct {
	Date         string `json:"date"`
	RegularCount int    `json:"count"`
	SpecialCount int    `json:"special_count"`
}

type MealActivityResponse struct {
	EmployeeId      uint              `json:"employee_id"`
	EmployeeName    string            `json:"employee_name"`
	EmployeeDetails []EmployeeDetails `json:"employee_details"`
}

type EmployeeDetails struct {
	Date    string        `json:"date"`
	Holiday bool          `json:"holiday"`
	Meal    []MealDetails `json:"meal"`
}

type MealDetails struct {
	MealType   int             `json:"meal_type"`
	MealStatus []StatusDetails `json:"meal_status"`
}

type StatusDetails struct {
	Status       bool     `json:"status"`
	GuestCount   int      `json:"guest_count"`
	Penalty      bool     `json:"penalty"`
	PenaltyScore *float64 `gorm:"type:decimal(10,2);" json:"penaltyScore"`
}

type MealSummaryReq struct {
	StartDate string `json:"start_date"`
	Days      int    `json:"days"`
}

type MealSummaryResponse struct {
	EmployeeId uint   `json:"employee_id"`
	Name       string `json:"name"`
	Lunch      int    `json:"lunch"`
	Snacks     int    `json:"snacks"`
}

type Employee struct {
	EmployeeId     uint           `json:"employee_id"`
	Name           string         `json:"employee_name" gorm:"column:employee_name"`
	PreferenceFood datatypes.JSON `json:"preference_food" gorm:"column:preference_food"`
}

type TotalMealCounts struct {
	TotalLunch  int `json:"total_lunch"`
	TotalSnacks int `json:"total_snacks"`
}

type MealSummaryForGraph struct {
	Month string `json:"month"`
	Year  string `json:"year"`
	Lunch int    `json:"lunch"`
	Snack int    `json:"snack"`
}

type MonthData struct {
	Month           string  `json:"month"`
	Year            string  `json:"year"`
	TotalLunch      int     `json:"total_lunch"`
	TotalGuestLunch int     `json:"total_guest_lunch"`
	TotalSnack      int     `json:"total_snack"`
	TotalGuestSnack int     `json:"total_guest_snack"`
	LunchPenalty    float64 `json:"lunch_penalty"`
	SnackPenalty    float64 `json:"snack_penalty"`
}
