package types

type EmployeeRequest struct {
	EmployeeId    uint   `json:"employee_id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	DeptID        int    `json:"dept_id"`
	Remarks       string `json:"remarks"`
	DefaultStatus bool   `json:"default_status"`
	IsAdmin       bool   `json:"is_admin"`
}

//func (employee *EmployeeRequest) Validate() error {
//	return validation.ValidateStruct(&employee,
//		validation.Field(&employee.Name,
//			validation.Required.Error("Employee name is required"),
//			validation.Length(1, 30)),
//		validation.Field(&employee.Email,
//			validation.Required.Error("Email is required"),
//			validation.Length(1, 30)),
//		validation.Field(&employee.Password,
//			validation.Required.Error("Password is required"),
//			validation.Length(1, 30)),
//	)
//}
