package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
	"meal-management/pkg/types"
	"net/http"
	"strconv"
)

var EmployeeService domain.IEmployeeService

func SetEmployeeService(empService domain.IEmployeeService) {
	EmployeeService = empService
}

func CreateEmployee(e echo.Context) error {
	reqEmployee := &types.EmployeeRequest{}
	if err := e.Bind(reqEmployee); err != nil {
		fmt.Println(err)
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}

	//if err := reqEmployee.Validate(); err != nil {
	//	return e.JSON(http.StatusBadRequest, err.Error())
	//}

	employee := &models.Employee{
		Name:          reqEmployee.Name,
		Email:         reqEmployee.Email,
		DeptID:        reqEmployee.DeptID,
		Password:      reqEmployee.Password,
		Remarks:       reqEmployee.Remarks,
		DefaultStatus: reqEmployee.DefaultStatus,
	}
	if err := EmployeeService.CreateEmployee(employee); err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, "Employee created successfully")
}
func GetEmployee(e echo.Context) error {
	tempEmployeeID := e.QueryParam("employee_id")
	EmployeeID, err := strconv.ParseUint(tempEmployeeID, 0, 0)
	if err != nil && tempEmployeeID != "" {
		return e.JSON(http.StatusBadRequest, "Invalid Employee ID")
	}
	Employee, err := EmployeeService.GetEmployee(uint(EmployeeID))

	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, Employee)
}
func UpdateEmployee(e echo.Context) error {
	reqEmployee := &types.EmployeeRequest{}

	if err := e.Bind(reqEmployee); err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Input")
	}

	//if err := reqEmployee.Validate(); err != nil {
	//	return e.JSON(http.StatusBadRequest, err.Error())
	//}

	if reqEmployee.EmployeeId == 0 {
		return e.JSON(http.StatusBadRequest, "EmployeeID is required and must be greater than zero")
	}
	EmployeeID := reqEmployee.EmployeeId
	//if err != nil {
	//	return e.JSON(http.StatusBadRequest, err.Error())
	//}

	existingEmployee, err := EmployeeService.GetEmployee(uint(EmployeeID))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	updatedEmployee := &models.Employee{
		EmployeeId:    uint(EmployeeID),
		Name:          ifNotEmpty(reqEmployee.Name, existingEmployee[0].Name),
		Email:         ifNotEmpty(reqEmployee.Email, existingEmployee[0].Email),
		Password:      ifNotEmpty(reqEmployee.Password, existingEmployee[0].Password),
		DeptID:        ifNotZero(reqEmployee.DeptID, existingEmployee[0].DeptID),
		Remarks:       ifNotEmpty(reqEmployee.Remarks, existingEmployee[0].Remarks),
		DefaultStatus: reqEmployee.DefaultStatus,
	}

	if err := EmployeeService.UpdateEmployee(updatedEmployee); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, "Employee was updated successfully")
}

func ifNotEmpty(new, existing string) string {
	if new != "" {
		return new
	}
	return existing
}

func ifNotZero(new, existing int) int {
	if new != 0 {
		return new
	}
	return existing
}

func DeleteEmployee(e echo.Context) error {
	tempEmployeeID := e.QueryParam("employee_id")
	EmployeeID, err := strconv.ParseUint(tempEmployeeID, 0, 0)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}

	_, err = EmployeeService.GetEmployee(uint(EmployeeID))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	if err := EmployeeService.DeleteEmployee(uint(EmployeeID)); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, "Employee was deleted successfully")
}
