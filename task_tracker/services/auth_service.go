package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tbaehler/gin-keycloak/pkg/ginkeycloak"
)

type UserRole string

const (
	RolePopug   UserRole = "popug"
	RoleManager UserRole = "manager"
	RoleAdmin   UserRole = "admin"
)

type AuthService interface {
	GetUserId(c *gin.Context) uuid.UUID
	GetUserRole(c *gin.Context) (UserRole, error)
}

type authService struct {
}

func (a authService) GetUserId(c *gin.Context) uuid.UUID {
	token, _ := c.Get("token")
	ginToken := token.(ginkeycloak.KeyCloakToken)
	return uuid.MustParse(ginToken.Sub)
}

func (a authService) GetUserRole(c *gin.Context) (UserRole, error) {
	token, _ := c.Get("token")
	ginToken := token.(ginkeycloak.KeyCloakToken)
	supportedRoles := []UserRole{RolePopug, RoleManager, RoleAdmin}
	for _, role := range ginToken.RealmAccess.Roles {
		for _, sRole := range supportedRoles {
			if role == string(sRole) {
				return sRole, nil
			}
		}
	}

	return "", fmt.Errorf("no supported roles found in provided token")
}

func NewAuthService() AuthService {
	return &authService{}
}
