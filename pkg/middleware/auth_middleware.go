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
			EmployeeID, parseErr := parseJWT(authorizationHeader)
			if parseErr != nil {
				return c.JSON(http.StatusBadRequest, parseErr)
			}
			c.Request().Header.Set(consts.UserIdHeader, EmployeeID)
			return next(c)
		}
	}
}

func GetEmployeeIDHandler(c echo.Context) error {
	authorizationHeader := c.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	EmployeeID, err := parseJWT(authorizationHeader)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"EmployeeID": EmployeeID})
}

func parseJWT(jwtToken string) (string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("jwtserversidesecret"), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", err
	}

	EmployeeID, found := claims[consts.UserIdHeader]
	if !found {
		return "", err
	}
	//return EmployeeID.(string), nil
	switch v := EmployeeID.(type) {
	case string:
		return v, nil
	case float64:
		return fmt.Sprintf("%.0f", v), nil // Convert float64 to string without decimal
	default:
		return "", fmt.Errorf("unexpected type for EmployeeID: %T", v)
	}
}
