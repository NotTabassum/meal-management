package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
)

func HolidayRoutes(e *echo.Echo) {
	holiday := e.Group("/holiday")
	holiday.POST("", controllers.CreateHoliday)
	holiday.GET("", controllers.GetHoliday)
	holiday.DELETE("", controllers.DeleteHoliday)
}
