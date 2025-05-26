package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
)

func EmployeeRoutes(e *echo.Echo) {
	emp := e.Group("/employee")

	emp.GET("/profile", controllers.Profile)
	emp.POST("", controllers.CreateEmployee)
	emp.GET("/hash", controllers.MakeHash)
	emp.GET("/photo", controllers.GetPhoto)
	emp.GET("", controllers.GetEmployee)
	emp.GET("/guest-list", controllers.GetGuestList)
	//emp.PATCH("/default-status", controllers.UpdateDefaultStatus)
	emp.PATCH("/default-status", controllers.UpdateDefaultStatusNew)
	emp.PATCH("/forget-password", controllers.ForgottenPassword)
	emp.PATCH("/password-change", controllers.PasswordChange)
	emp.PATCH("", controllers.UpdateEmployee)
	emp.DELETE("", controllers.DeleteEmployee)
	emp.GET("/telegram", controllers.TelegramMessage)
}
