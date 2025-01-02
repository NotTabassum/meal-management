package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
)

func EmployeeRoutes(e *echo.Echo) {
	emp := e.Group("/employee")

	emp.GET("/profile", controllers.Profile)
	emp.POST("", controllers.CreateEmployee)
	emp.GET("", controllers.GetEmployee)
	emp.PATCH("", controllers.UpdateEmployee)
	emp.DELETE("", controllers.DeleteEmployee)
}
