package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
)

func MealActivityRoutes(e *echo.Echo) {
	emp := e.Group("/meal_activity")
	emp.POST("", controllers.CreateMealActivity)
	emp.PATCH("", controllers.UpdateMealActivity)
	emp.GET("", controllers.GetMealActivity)
	//emp.GET("", controllers.GetMealActivity)
}
