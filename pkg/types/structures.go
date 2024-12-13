package types

type EmployeeRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	DeptID string `json:"department"`
}
