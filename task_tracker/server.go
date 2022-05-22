package main

import (
	"fmt"
	"github.com/szuecs/gin-glog"
	"github.com/tbaehler/gin-keycloak/pkg/ginkeycloak"
	"os"
	"strings"
	"sync"
	"task_tracker/api"
	"task_tracker/messaging"
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

	userConsumer := messaging.NewUserInfoConsumer()
	api.Init(router, services.NewAuthService(), services.NewTaskService(), services.NewUserService(&userConsumer))

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		newUserTopicName := os.Getenv("KAFKA_TOPIC_NEW_USER")
		brokersUrl := os.Getenv("KAFKA_BROKERS_URL")
		err := userConsumer.ReadFromTopic(strings.Split(brokersUrl, ","), newUserTopicName)
		fmt.Println(err)
		wg.Done()
	}()
	go func() {
		err := router.Run("localhost:" + os.Getenv("PORT"))
		fmt.Println(err)
		wg.Done()
	}()

	fmt.Println("Started...")
	wg.Wait()
}
