package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"meal-management/pkg/domain"
	"meal-management/pkg/middleware"
	"meal-management/pkg/models"
	"net/http"
	"strconv"
)

var DeptService domain.IDeptService

func SetDeptService(deptService domain.IDeptService) {
	DeptService = deptService
}

func CreateDepartment(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	_, isAdmin, _ := middleware.ParseJWT(authorizationHeader)
	if !isAdmin {
		return e.JSON(http.StatusForbidden, map[string]string{"res": "Unauthorized"})
	}
	reqDept := &models.Department{}
	if err := e.Bind(reqDept); err != nil {
		fmt.Println(err)
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}
	dept := &models.Department{
		DeptID:   reqDept.DeptID,
		DeptName: reqDept.DeptName,
		Weekend:  reqDept.Weekend,
	}

	if err := DeptService.CreateDepartment(dept); err != nil {
		fmt.Println(err)
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, "New Department is created successfully")
}

func DeleteDepartment(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	_, isAdmin, _ := middleware.ParseJWT(authorizationHeader)
	if !isAdmin {
		return e.JSON(http.StatusForbidden, map[string]string{"res": "Unauthorized"})
	}

	tempDept := e.QueryParam("dept_id")
	deptID, err := strconv.ParseInt(tempDept, 0, 0)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}

	if err := DeptService.DeleteDepartment(int(deptID)); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, "Department is deleted successfully")

}

func GetAllDept(e echo.Context) error {
	departments, err := DeptService.GetAllDepartments()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch departments"})
	}
	return e.JSON(http.StatusOK, departments)
}

func UpdateDepartment(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	_, isAdmin, _ := middleware.ParseJWT(authorizationHeader)
	if !isAdmin {
		return e.JSON(http.StatusForbidden, map[string]string{"res": "Unauthorized"})
	}

	reqDept := &models.Department{}
	if err := e.Bind(reqDept); err != nil {
		fmt.Println(err)
	}
	if reqDept.DeptID == 0 {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}
	dept := &models.Department{
		DeptID:   reqDept.DeptID,
		DeptName: reqDept.DeptName,
		Weekend:  reqDept.Weekend,
	}
	if err := DeptService.UpdateDepartment(dept); err != nil {
		fmt.Println(err)
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, "department is updated successfully")
}
