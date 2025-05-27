package types

import "gorm.io/datatypes"

type EmployeeRequest struct {
	EmployeeId  uint   `json:"employee_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	DeptName    string `json:"dept_name"`
	Remarks     string `json:"remarks"`
	//DefaultStatus bool   `json:"default_status"`
	DefaultStatusLunch  bool           `json:"default_status_lunch"`
	DefaultStatusSnacks bool           `json:"default_status_snacks"`
	IsAdmin             bool           `json:"is_admin"`
	PreferenceFood      datatypes.JSON `json:"preference_food"`
	//IsPermanent         bool           `json:"is_permanent"`
	IsActive    bool   `json:"is_active"`
	Designation string `json:"designation"`
	Roll        string `json:"roll"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email"`
	Link  string `json:"link"`
}

type PasswordRequest struct {
	Password string `json:"password"`
}

//type DefaultStatus struct {
//	Date   string `json:"date"`
//	Status bool   `json:"status"`
//}

type DefaultStatusNew struct {
	Date     string `json:"date"`
	Status   bool   `json:"status"`
	MealType int    `json:"meal_type"`
}
