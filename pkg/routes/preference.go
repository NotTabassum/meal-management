package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
)

func Preference(e *echo.Echo) {
	mp := e.Group("/preference")
	mp.POST("", controllers.CreatePreference)
	mp.GET("", controllers.GetPreference)
	//mp.PATCH("", controllers.UpdatePreference)
	//mp.DELETE("", controllers.DeletePreference)
}
