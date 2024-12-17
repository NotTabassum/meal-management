package services

import (
	"errors"
	"fmt"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
	"meal-management/pkg/types"
)

type EmployeeService struct {
	repo domain.IEmployeeRepo
}

func EmployeeServiceInstance(employeeRepo domain.IEmployeeRepo) domain.IEmployeeService {
	return &EmployeeService{
		repo: employeeRepo,
	}
}

func (service *EmployeeService) GetEmployee(EmployeeID uint) ([]types.EmployeeRequest, error) {
	var allEmployees []types.EmployeeRequest
	employee := service.repo.GetEmployee(EmployeeID)
	if len(employee) == 0 {
		//fmt.Println(EmployeeID)
		return nil, errors.New("Employee not found")
	}
	for _, val := range employee {
		allEmployees = append(allEmployees, types.EmployeeRequest{
			EmployeeId:    val.EmployeeId,
			Name:          val.Name,
			Email:         val.Email,
			DeptID:        val.DeptID,
			Password:      val.Password,
			Remarks:       val.Remarks,
			DefaultStatus: val.DefaultStatus,
		})
	}
	return allEmployees, nil
}
func (service *EmployeeService) CreateEmployee(employee *models.Employee) error {
	if err := service.repo.CreateEmployee(employee); err != nil {
		return errors.New("Employee was not created")
	}
	return nil
}
func (service *EmployeeService) UpdateEmployee(employee *models.Employee) error {
	if err := service.repo.UpdateEmployee(employee); err != nil {
		return errors.New("Employee update was unsuccessful")
	}
	return nil
}
func (service *EmployeeService) DeleteEmployee(EmployeeId uint) error {
	if err := service.repo.DeleteEmployee(EmployeeId); err != nil {
		fmt.Println(err)
		return errors.New("Employee was not deleted")
	}
	return nil
}
