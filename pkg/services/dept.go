package services

import (
	"errors"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
)

type DeptService struct {
	repo domain.IDeptRepo
}

func DeptServiceInstance(deptRepo domain.IDeptRepo) domain.IDeptRepo {
	return &DeptService{
		repo: deptRepo,
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
	if err := service.repo.DeleteDepartment(deptID); err != nil {
		return errors.New("department was not created")
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
