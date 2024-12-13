package containers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"meal-management/pkg/config"
	"meal-management/pkg/routes"
)

func Serve(e *echo.Echo) {
	config.SetConfig()
	routes.EmployeeRoutes(e)
	log.Fatal(e.Start(fmt.Sprintf("%s", config.LocalConfig.DBPort)))
}
