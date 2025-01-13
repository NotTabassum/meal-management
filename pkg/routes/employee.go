package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
)

func EmployeeRoutes(e *echo.Echo) {
	emp := e.Group("/employee")

	emp.GET("/profile", controllers.Profile)
	emp.POST("", controllers.CreateEmployee)
	//emp.GET("/hash", controllers.MakeHash)
	emp.GET("/photo", controllers.GetPhoto)
	emp.GET("", controllers.GetEmployee)
	emp.PATCH("/default-status", controllers.UpdateDefaultStatus)
	emp.PATCH("/forget-password", controllers.ForgottenPassword)
	emp.PATCH("", controllers.UpdateEmployee)
	emp.DELETE("", controllers.DeleteEmployee)
}
