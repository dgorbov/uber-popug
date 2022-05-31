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
	userService services.UserService
}

type TaskDto struct {
	AssignedName string `json:"assigned_name"`
	services.Task
}

func Init(router *gin.Engine, authService services.AuthService, taskService services.TaskService, userService services.UserService) {
	controller := controller{taskService: taskService, authService: authService, userService: userService}

	router.POST("/tasks", controller.createTask)
	router.GET("/tasks/:id", controller.getTask)
	router.GET("/tasks/", controller.getAllTasks)
	router.GET("/tasks/my", controller.getAllMyTasks)
	router.POST("/tasks/:id/done", controller.completeTask)
	router.POST("/tasks/shuffle", controller.shuffleTasks)
}

func createTaskDto(task services.Task, user services.UserInfo) TaskDto {
	return TaskDto{user.Name, task}
}

func (con *controller) createTaskDtoArray(tasks []services.Task) ([]TaskDto, error) {
	tasksDto := make([]TaskDto, len(tasks))
	for idx, task := range tasks {
		user, err := con.userService.GetUser(task.Assigned)
		if err != nil {
			return nil, err
		}
		tasksDto[idx] = createTaskDto(task, user)
	}

	return tasksDto, nil
}

func (con *controller) getMyUserInfo(c *gin.Context) (services.UserInfo, error) {
	myId := con.authService.GetUserId(c)
	return con.userService.GetUser(myId)
}

func (con *controller) createTask(c *gin.Context) {
	type CreateTaskRequest struct {
		Description string `json:"description"`
	}

	var ctr CreateTaskRequest
	err := c.ShouldBindJSON(&ctr)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	user, err := con.userService.GetRandomUser()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	task := con.taskService.CreateTask(ctr.Description, user.Id)
	c.JSON(http.StatusOK, createTaskDto(task, user))
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
		user, err := con.userService.GetUser(task.Assigned)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, createTaskDto(task, user))
	}
}

func (con *controller) getAllTasks(c *gin.Context) {
	role, err := con.authService.GetUserRole(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if role != services.RoleAdmin && role != services.RoleManager {
		c.String(http.StatusForbidden, err.Error())
		return
	}

	allTasks := con.taskService.GetAllTasks()
	allTasksDto, err := con.createTaskDtoArray(allTasks)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allTasksDto)
}

func (con *controller) getAllMyTasks(c *gin.Context) {
	myUser, err := con.getMyUserInfo(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	myTasks := con.taskService.GetAllUserTasks(myUser.Id)
	myTasksDto, _ := con.createTaskDtoArray(myTasks)
	c.JSON(http.StatusOK, myTasksDto)
}

func (con *controller) completeTask(c *gin.Context) {
	task, err := con.getTaskByIdOrWriteError(c)
	if err != nil {
		return
	}

	myUser, err := con.getMyUserInfo(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	task, err = con.taskService.CompleteTask(task.Id, myUser.Id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, createTaskDto(task, myUser))
}

func (con *controller) shuffleTasks(c *gin.Context) {
	role, _ := con.authService.GetUserRole(c)
	c.String(http.StatusNotImplemented, string(role))
}
