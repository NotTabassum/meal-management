package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
)

func MealPlanRoutes(e *echo.Echo) {
	mp := e.Group("/mealplan")
	mp.POST("/create", controllers.CreateMealPlan)
	mp.GET("/get/:start/:days", controllers.GetMealPlan)
	//mp.GET("/get/date/mealType", controllers.GetMealPlanByPrimaryKey)
	mp.PATCH("/update/:date/:meal_type", controllers.UpdateMealPlan)
	mp.DELETE("/delete/:date/:meal_type", controllers.DeleteMealPlan)
}
