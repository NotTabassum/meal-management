package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
)

func MealActivityRoutes(e *echo.Echo) {
	emp := e.Group("/meal_activity")

	emp.POST("", controllers.CreateMealActivity)
	emp.PATCH("/total-meal-summary", controllers.TotalMealCount)
	emp.PATCH("/total-meal-group", controllers.TotalMealADayGroup)
	//emp.PATCH("/total-meal", controllers.TotalMealADay)
	emp.PATCH("/total-penalty", controllers.TotalPenalty)
	emp.PATCH("/meal-summary", controllers.TotalMealAMonth)
	emp.PATCH("/meal-per-person", controllers.TotalMealPerPerson)
	emp.PATCH("/group-update", controllers.UpdateGroupMealActivity)
	emp.PATCH("", controllers.UpdateMealActivity)
	emp.GET("/meal-summary-graph", controllers.MealSummaryForGraph)
	emp.GET("", controllers.GetOwnMealActivity)
	emp.GET("/admin", controllers.GetMealActivity)
}
