package controllers

import (
	"meal-management/pkg/domain"
)

var AuthService domain.IAuthService

func SetAuthService(auth domain.IAuthService) {
	AuthService = auth
}
