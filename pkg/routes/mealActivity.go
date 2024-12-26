package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
)

func MealActivityRoutes(e *echo.Echo) {
	emp := e.Group("/meal_activity")

	emp.POST("", controllers.CreateEmployee)
	emp.GET("", controllers.GetEmployee)
	emp.PATCH("", controllers.UpdateEmployee)
	emp.DELETE("", controllers.DeleteEmployee)
}
