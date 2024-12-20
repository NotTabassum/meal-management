package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
)

func MealPlanRoutes(e *echo.Echo) {
	mp := e.Group("/mealplan")
	//mp.Use(middleware.RequestCounter)

	mp.POST("", controllers.CreateMealPlan)
	mp.GET("", controllers.GetMealPlan)
	mp.PATCH("", controllers.UpdateMealPlan)
	mp.DELETE("", controllers.DeleteMealPlan)
}
