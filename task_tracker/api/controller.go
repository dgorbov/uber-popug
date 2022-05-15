package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"task_tracker/services"
)

type controller struct {
	taskService services.TaskService
	authService services.AuthService
}

func Init(router *gin.Engine, authService services.AuthService, taskService services.TaskService) {
	controller := controller{taskService: taskService, authService: authService}

	router.POST("/tasks", controller.createTask)
	router.GET("/tasks/:id", controller.getTask)
	router.GET("/tasks/my", controller.getAllMyTasks)
	router.POST("/tasks/:id/done", controller.completeTask)
	router.POST("/tasks/shuffle", controller.shuffleTasks)
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

func (con *controller) getTaskByIdOrWriteError(c *gin.Context) (services.Task, error) {
	id, err := uuid.Parse(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return services.Task{}, err
	}

	task, err := con.taskService.GetTask(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return services.Task{}, err
	}

	return task, nil
}

func (con *controller) getTask(c *gin.Context) {
	task, err := con.getTaskByIdOrWriteError(c)
	if err == nil {
		c.JSON(http.StatusOK, task)
	}
}

func (con *controller) getAllMyTasks(c *gin.Context) {
	myId := con.authService.GetUserId(c)
	myTasks := con.taskService.GetAllUserTasks(myId)
	c.JSON(http.StatusOK, myTasks)
}

func (con *controller) completeTask(c *gin.Context) {
	task, err := con.getTaskByIdOrWriteError(c)
	if err != nil {
		return
	}

	userId := con.authService.GetUserId(c)
	task, err = con.taskService.CompleteTask(task.Id, userId)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, task)
}

func (con *controller) shuffleTasks(c *gin.Context) {
	role, _ := con.authService.GetUserRole(c)
	c.String(http.StatusNotImplemented, string(role))
}
