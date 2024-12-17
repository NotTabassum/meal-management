package types

import validation "github.com/Go-ozzo/ozzo-validation"

type EmployeeRequest struct {
	EmployeeId    uint   `json:"EmployeeID"`
	Name          string `json:"Name"`
	Email         string `json:"Email"`
	DeptID        int    `json:"DeptID"`
	Password      string `json:"Password"`
	Remarks       string `json:"Remarks"`
	DefaultStatus bool   `json:"DefaultStatus"`
}

func (employee *EmployeeRequest) Validate() error {
	return validation.ValidateStruct(&employee,
		validation.Field(&employee.Name,
			validation.Required.Error("Employee name is required"),
			validation.Length(1, 30)),
		validation.Field(&employee.Email,
			validation.Required.Error("Email is required"),
			validation.Length(1, 30)),
		validation.Field(&employee.Password,
			validation.Required.Error("Password is required"),
			validation.Length(1, 30)),
	)
}
