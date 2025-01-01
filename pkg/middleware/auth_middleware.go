package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"meal-management/pkg/consts"
	"net/http"
)

func Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorizationHeader := c.Request().Header.Get("Authorization")
			if authorizationHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
			}
			EmployeeID, isAdmin, parseErr := ParseJWT(authorizationHeader)
			if parseErr != nil {
				return c.JSON(http.StatusBadRequest, parseErr)
			}
			c.Request().Header.Set(consts.UserIdHeader, EmployeeID)
			c.Request().Header.Set("isAdmin", fmt.Sprintf("%v", isAdmin))
			return next(c)
		}
	}
}

func GetEmployeeIDHandler(c echo.Context) error {
	authorizationHeader := c.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	EmployeeID, _, err := ParseJWT(authorizationHeader)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"EmployeeID": EmployeeID})
}

func ParseJWT(jwtToken string) (string, bool, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("jwtserversidesecret"), nil
	})
	if err != nil {
		return "", false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", false, err
	}

	EmployeeID, found := claims["employee_id"]

	if !found {
		return "", false, err
	}
	isAdmin, ok := claims["is_admin"].(bool)
	if !ok {
		return "", false, fmt.Errorf("isAdmin is not a boolean")
	}

	//return EmployeeID, isAdmin, nil
	//m := make(map[string]interface{})

	//if m["isAdmin"] == nil {
	//	isAdmin = false
	//}
	switch v := EmployeeID.(type) {
	case string:
		return v, isAdmin, nil
	case float64:
		return fmt.Sprintf("%.0f", v), isAdmin, nil // Convert float64 to string without decimal
	default:
		return "", isAdmin, fmt.Errorf("unexpected type for EmployeeID: %T", v)
	}
}