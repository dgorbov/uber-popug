package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"task_tracker/services"
)

type controller struct {
	taskService services.TaskService
}

func Init(router *gin.Engine, taskService services.TaskService) {
	controller := controller{taskService: taskService}

	router.POST("/tasks/", controller.createTask)
	router.GET("/tasks/:id", controller.getTask)
}

func (con *controller) createTask(c *gin.Context) {
	type CreateTaskRequest struct {
		Description string    `json:"description"`
		Assigned    uuid.UUID `json:"assigned"`
	}

	var ctr CreateTaskRequest
	err := c.ShouldBindJSON(&ctr)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	task := con.taskService.CreateTask(ctr.Description, ctr.Assigned)
	c.JSON(http.StatusOK, task)
}

func (con *controller) getTask(c *gin.Context) {
	id, err := uuid.Parse(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	task, err := con.taskService.GetTask(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}
