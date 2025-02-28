package types

import "gorm.io/datatypes"

type EmployeeRequest struct {
	EmployeeId     uint           `json:"employee_id"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	PhoneNumber    string         `json:"phone_number"`
	DeptName       string         `json:"dept_name"`
	Remarks        string         `json:"remarks"`
	DefaultStatus  bool           `json:"default_status"`
	IsAdmin        bool           `json:"is_admin"`
	PreferenceFood datatypes.JSON `json:"preference_food"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email"`
	Link  string `json:"link"`
}

type PasswordRequest struct {
	Password string `json:"password"`
}
