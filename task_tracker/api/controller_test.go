package api

import (
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"strings"
	"task_tracker/services"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type authServiceStub struct {
	userId   uuid.UUID
	userRole services.UserRole
}

func (a *authServiceStub) GetUserId(_ *gin.Context) uuid.UUID {
	return a.userId
}

func (a *authServiceStub) GetUserRole(_ *gin.Context) (services.UserRole, error) {
	return a.userRole, nil
}

type userServiceStub struct {
	stubUser services.UserInfo
}

func (u *userServiceStub) GetUser(userId uuid.UUID) (services.UserInfo, error) {
	return u.stubUser, nil
}

func (u *userServiceStub) GetRandomUser() (services.UserInfo, error) {
	return u.stubUser, nil
}

func Test_CreateTask_ProvideValidPayload_TaskCreated(t *testing.T) {
	router := gin.New()
	Init(router, &authServiceStub{}, services.NewTaskService(), &userServiceStub{})

	req, _ := http.NewRequest(
		"POST", "/tasks/",
		strings.NewReader(`{"description": "test_resources task for TestCreateTask", "assigned": "8d278e16-5da5-4105-a0d4-6b7a8fa4e163"}`))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("status: %d", w.Code)
	t.Logf("response: %s", w.Body.String())
}

func Test_GetTask_ProvideExistingTaskId_ReturnTask(t *testing.T) {
	ts := services.NewTaskService()
	task := ts.CreateTask("test_resources task for TestGetTask", uuid.New())

	router := gin.New()
	Init(router, &authServiceStub{}, ts, &userServiceStub{})

	req, _ := http.NewRequest("GET", "/tasks/"+task.Id.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("status: %d", w.Code)
	t.Logf("response: %s", w.Body.String())
}

func Test_GetAllMyTasks_CreateSeveralTaskAndAssigneeToUserRequester_ReturnAllAssignedTasks(t *testing.T) {
	userA := uuid.New()
	userB := uuid.New()

	ts := services.NewTaskService()
	ts.CreateTask("test_resources task1 for userA", userA)
	ts.CreateTask("test_resources task2 for userA", userA)
	ts.CreateTask("test_resources task3 for userB", userB)

	router := gin.New()
	Init(router, &authServiceStub{userId: userA}, ts, &userServiceStub{})

	req, _ := http.NewRequest("GET", "/tasks/my", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	t.Logf("response: %s", w.Body.String())
}

func Test_CompleteTask_ProvideValidTaskIdAndUserId_TaskCompleted(t *testing.T) {
	user := services.UserInfo{Id: uuid.New(), Name: "task owner"}
	ts := services.NewTaskService()
	task := ts.CreateTask("test_resources task for user to complete", user.Id)

	router := gin.New()
	Init(router, &authServiceStub{userId: user.Id}, ts, &userServiceStub{stubUser: user})

	req, _ := http.NewRequest("POST", "/tasks/"+task.Id.String()+"/done", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	t.Logf("response: %s", w.Body.String())
}
