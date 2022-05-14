package main

import (
	"os"
	"task_tracker/api"
	"task_tracker/services"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	api.Init(router, services.NewTaskService())
	router.Run("localhost:" + os.Getenv("PORT"))
}
