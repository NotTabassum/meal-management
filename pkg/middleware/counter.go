package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"sync/atomic"
)

var mealPlanCounter uint64

func RequestCounter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		count := atomic.AddUint64(&mealPlanCounter, 1)

		fmt.Printf("The /mealplan endpoint has been called %d times\n", count)

		return next(c)
	}
}
