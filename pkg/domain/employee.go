package domain

import (
	"meal-management/pkg/models"
	"meal-management/pkg/types"
)

type IEmployeeRepo interface {
	CreateEmployee(employee *models.Employee) error
}

type IEmployeeService interface {
	GetEmployee() ([]types.EmployeeRequest, error)
}
