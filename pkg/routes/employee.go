package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
)

func EmployeeRoutes(e *echo.Echo) {
	emp := e.Group("/employee")

	emp.POST("/create", controllers.CreateEmployee)
	emp.GET("/list", controllers.GetEmployee)
	emp.PATCH("/update/:EmployeeID", controllers.UpdateEmployee)
	emp.DELETE("/delete/:EmployeeID", controllers.DeleteEmployee)
}
