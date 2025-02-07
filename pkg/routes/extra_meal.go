package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
)

func ExtraMealRoutes(e *echo.Echo) {
	emp := e.Group("/extra_meal")

	emp.POST("", controllers.CreateExtraMeal)
	emp.PATCH("", controllers.UpdateExtraMeal)
	emp.GET("", controllers.FetchExtraMeal)
}
