package services

import (
	"errors"
	"fmt"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
)

type DeptService struct {
	repo  domain.IDeptRepo
	repo2 domain.IEmployeeRepo
}

func DeptServiceInstance(deptRepo domain.IDeptRepo, employeeRepo domain.IEmployeeRepo) domain.IDeptService {
	return &DeptService{
		repo:  deptRepo,
		repo2: employeeRepo,
	}
}

func (service *DeptService) CreateDepartment(dept *models.Department) error {
	if err := service.repo.CreateDepartment(dept); err != nil {
		return errors.New("department was not created")
	}
	return nil
}
func (service *DeptService) UpdateDepartment(dept *models.Department) error {
	if err := service.repo.UpdateDepartment(dept); err != nil {
		return errors.New("department was not updated")
	}
	return nil
}

func (service *DeptService) DeleteDepartment(deptID int) error {
	employee, err := service.repo2.GetEmployeeByDepartment(deptID)
	if err != nil {
		return err
	}
	if employee != nil {
		return fmt.Errorf("department still has %d employees", len(employee))
	}
	if err := service.repo.DeleteDepartment(deptID); err != nil {
		return errors.New("department does not exist")
	}
	return nil
}

func (service *DeptService) GetAllDepartments() ([]models.Department, error) {
	var departments []models.Department
	departments, err := service.repo.GetAllDepartments()
	if err != nil {
		return []models.Department{}, err
	}
	return departments, nil
}
