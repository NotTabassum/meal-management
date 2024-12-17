package containers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"meal-management/pkg/config"
	"meal-management/pkg/connection"
	"meal-management/pkg/controllers"
	"meal-management/pkg/repositories"
	"meal-management/pkg/routes"
	"meal-management/pkg/services"
)

func Serve(e *echo.Echo) {
	config.SetConfig()
	db := connection.GetDB()
	EmployeeRepo := repositories.EmloyeeDBInstance(db)
	EmployeeService := services.EmployeeServiceInstance(EmployeeRepo)
	controllers.SetEmployeeService(EmployeeService)
	routes.EmployeeRoutes(e)
	log.Fatal(e.Start(fmt.Sprintf(":%s", config.LocalConfig.ServerPort)))
}
