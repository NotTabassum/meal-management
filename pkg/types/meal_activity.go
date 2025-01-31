package types

type MealActivityRequest struct {
	Date       string `json:"date" validate:"required"`
	EmployeeId uint   `json:"employee_id" validate:"required"`
	MealType   int    `json:"meal_type" validate:"required"`
	Status     *bool  `json:"status"`
	GuestCount *int   `json:"guest_count"`
	Penalty    *bool  `json:"penalty"`
	IsOffDay   bool   `json:"is_off_day"`
}

type PenaltyRequest struct {
	Date       string `json:"date"`
	EmployeeId uint   `json:"employee_id"`
	Days       int    `json:"days"`
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
	Status     bool `json:"status"`
	GuestCount int  `json:"guest_count"`
	Penalty    bool `json:"penalty"`
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
	EmployeeId uint   `json:"employee_id"`
	Name       string `json:"name"`
}
