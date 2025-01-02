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

	for _, val := range employee {
		allEmployees = append(allEmployees, types.EmployeeRequest{
			EmployeeId:    val.EmployeeId,
			Name:          val.Name,
			Email:         val.Email,
			PhoneNumber:   val.PhoneNumber,
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
			PhoneNumber:   val.PhoneNumber,
			DeptID:        val.DeptID,
			Password:      val.Password,
			Remarks:       val.Remarks,
			DefaultStatus: val.DefaultStatus,
			IsAdmin:       val.IsAdmin,
		})
	}
	return allEmployees, nil
}

//func (service *EmployeeService) SaveFile(file *multipart.FileHeader, destDir string) (string, error) {
//	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
//		return "", err
//	}
//
//	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
//	destPath := filepath.Join(destDir, filename)
//
//	src, err := file.Open()
//	if err != nil {
//		return "", err
//	}
//	defer func(src multipart.File) {
//		err := src.Close()
//		if err != nil {
//			return
//		}
//	}(src)
//
//	dst, err := os.Create(destPath)
//	if err != nil {
//		return "", err
//	}
//	defer func(dst *os.File) {
//		err := dst.Close()
//		if err != nil {
//			return
//		}
//	}(dst)
//
//	if _, err := dst.ReadFrom(src); err != nil {
//		return "", err
//	}
//
//	return destPath, nil
//}

//func (service *EmployeeService) SaveFile(file multipart.File, fileHeader *multipart.FileHeader, photoDir string) (string, error) {
//	// Ensure the directory exists
//	err := os.MkdirAll(photoDir, os.ModePerm) // os.ModePerm gives full access permissions (0755)
//	if err != nil {
//		fmt.Printf("Failed to create directory %s: %v\n", photoDir, err)
//		return "", err
//	}
//
//	// Extract the file extension
//	ext := filepath.Ext(fileHeader.Filename)
//	if ext == "" {
//		return "", fmt.Errorf("file has no extension")
//	}
//
//	// Generate a unique filename based on the current time and the original extension
//	fileName := fmt.Sprintf("%d%s", time.Now().Unix(), ext)
//	filePath := filepath.Join(photoDir, fileName) // Save it to /tmp/photos or another writable dir
//
//	// Create a new file on the system where the photo will be saved
//	outFile, err := os.Create(filePath)
//	if err != nil {
//		fmt.Println("Failed to create file %s: %v", filePath, err)
//		return "", err
//	}
//	defer outFile.Close()
//
//	// Copy the content of the uploaded file to the new file
//	_, err = io.Copy(outFile, file)
//	if err != nil {
//		fmt.Println("Failed to copy file content to %s: %v", filePath, err)
//		return "", err
//	}
//
//	// Return the file path of the saved photo
//	return filePath, nil
//}
