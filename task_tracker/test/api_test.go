package test

import (
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"strings"
	"task_tracker/api"
	"task_tracker/services"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateTask(t *testing.T) {
	router := gin.New()
	api.Init(router, services.NewTaskService())

	req, _ := http.NewRequest(
		"POST", "/tasks/",
		strings.NewReader(`{"description": "test task for TestCreateTask", "assigned": "8d278e16-5da5-4105-a0d4-6b7a8fa4e163"}`))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("status: %d", w.Code)
	t.Logf("response: %s", w.Body.String())
}

func TestGetTask(t *testing.T) {
	ts := services.NewTaskService()
	task := ts.CreateTask("test task for TestGetTask", uuid.New())

	router := gin.New()
	api.Init(router, ts)

	req, _ := http.NewRequest("GET", "/tasks/"+task.Id.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("status: %d", w.Code)
	t.Logf("response: %s", w.Body.String())
}
