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
	allEmployees := []types.EmployeeRequest{}
	employee := service.repo.GetEmployee(EmployeeID)
	//if len(employee) == 0 {
	//	//fmt.Println(EmployeeID)
	//	return nil, errors.New("employee not found")
	//}
	for _, val := range employee {
		allEmployees = append(allEmployees, types.EmployeeRequest{
			EmployeeId:    val.EmployeeId,
			Name:          val.Name,
			Email:         val.Email,
			DeptID:        val.DeptID,
			Remarks:       val.Remarks,
			DefaultStatus: val.DefaultStatus,
			IsAdmin:       val.IsAdmin,
		})
	}
	return allEmployees, nil
}
func (service *EmployeeService) CreateEmployee(employee *models.Employee) error {
	if err := service.repo.CreateEmployee(employee); err != nil {
		return errors.New("employee was not created")
	}
	return nil
}
func (service *EmployeeService) UpdateEmployee(employee *models.Employee) error {
	if err := service.repo.UpdateEmployee(employee); err != nil {
		return errors.New("employee update was unsuccessful")
	}
	return nil
}
func (service *EmployeeService) DeleteEmployee(EmployeeId uint) error {
	if err := service.repo.DeleteEmployee(EmployeeId); err != nil {
		fmt.Println(err)
		return errors.New("employee was not deleted")
	}
	return nil
}

func (service *EmployeeService) GetEmployeeWithPassword(EmployeeID uint) ([]models.Employee, error) {
	allEmployees := []models.Employee{}
	employee := service.repo.GetEmployee(EmployeeID)
	//if len(employee) == 0 {
	//	//fmt.Println(EmployeeID)
	//	return nil, errors.New("employee not found")
	//}
	for _, val := range employee {
		allEmployees = append(allEmployees, models.Employee{
			EmployeeId:    val.EmployeeId,
			Name:          val.Name,
			Email:         val.Email,
			DeptID:        val.DeptID,
			Password:      val.Password,
			Remarks:       val.Remarks,
			DefaultStatus: val.DefaultStatus,
			IsAdmin:       val.IsAdmin,
		})
	}
	return allEmployees, nil
}
