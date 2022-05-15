package main

import (
	"github.com/szuecs/gin-glog"
	"github.com/tbaehler/gin-keycloak/pkg/ginkeycloak"
	"os"
	"task_tracker/api"
	"task_tracker/services"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.Use(ginglog.Logger(3 * time.Second))
	router.Use(ginkeycloak.RequestLogger([]string{"uid"}, "data"))
	router.Use(gin.Recovery())

	keycloakCertPath := "realms/uber-popug/protocol/openid-connect/certs"
	keycloakConfig := ginkeycloak.KeycloakConfig{
		Url:           os.Getenv("KEYCLOAK_URL"),
		Realm:         "uber-popug",
		FullCertsPath: &keycloakCertPath,
	}

	auth := ginkeycloak.Auth(ginkeycloak.AuthCheck(), keycloakConfig)
	router.Use(auth)
	api.Init(router, services.NewAuthService(), services.NewTaskService())
	router.Run("localhost:" + os.Getenv("PORT"))
}
