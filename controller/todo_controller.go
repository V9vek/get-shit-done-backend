package controller

import (
	"fmt"
	"get-shit-done/model"
	"get-shit-done/service"
	"get-shit-done/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TodoController struct {
	todoService *service.TodoService
	jwtService  *service.JWTAuth
}

func NewTodoController(todoService *service.TodoService, jwtService *service.JWTAuth) *TodoController {
	return &TodoController{todoService: todoService, jwtService: jwtService}
}

// func getAccessTokenFromCookies(w http.ResponseWriter, r *http.Request) string {
// 	accessTokenCookie, err := r.Cookie("access_token")
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("access token not found: %v", err), http.StatusUnauthorized)
// 		return ""
// 	}

// 	return accessTokenCookie.Value
// }

func (c *TodoController) AddTodo(writer http.ResponseWriter, requests *http.Request) {
	var todo model.Todo
	if err := utils.ReadFromRequestBody(requests, &todo); err != nil {
		http.Error(writer, fmt.Sprintf("Invalid field for a todo: %v", err), http.StatusBadRequest)
		return
	}

	// setting the userid from the access token
	userIdStr, err := c.jwtService.GetSubjectFromAccessToken(getAccessTokenFromHeaders(writer, requests))
	if err != nil {
		http.Error(writer, fmt.Sprintf("token's subject has invalid format: %v", err), http.StatusUnauthorized)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(writer, fmt.Sprintf("can not parse subject from token: %v", err), http.StatusBadRequest)
		return
	}

	todo.UserId = userId

	err = c.todoService.AddTodo(requests.Context(), todo)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	webResponse := model.WebResponse{
		Code:   http.StatusOK,
		Status: "todo added successfully",
		Data:   todo,
	}
	utils.WriteResponseBody(writer, webResponse)
}

func (c *TodoController) UpdateTodo(writer http.ResponseWriter, requests *http.Request) {
	todoIdParam := chi.URLParam(requests, "todoId")
	todoId, err := strconv.Atoi(todoIdParam)
	if err != nil {
		http.Error(writer, fmt.Sprintf("can not parse todoId params: %v", err), http.StatusBadRequest)
		return
	}

	var todo model.Todo
	if err := utils.ReadFromRequestBody(requests, &todo); err != nil {
		http.Error(writer, fmt.Sprintf("Invalid field for a todo: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", todo)

	// setting the userid from the access token
	// userIdStr, err := c.jwtService.GetSubjectFromAccessToken(getAccessTokenFromHeaders(writer, requests))
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("token's subject has invalid format: %v", err), http.StatusUnauthorized)
	// 	return
	// }

	// userId, err := strconv.Atoi(userIdStr)
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("can not parse subject from token: %v", err), http.StatusBadRequest)
	// 	return
	// }

	// fetch all details of Todo
	oldTodo, err := c.todoService.FindTodoByTodoId(todoId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	if todo.Title != "" {
		oldTodo.Title = todo.Title
	}
	if todo.Description != "" {
		oldTodo.Description = todo.Description
	}
	if oldTodo.IsCompleted != todo.IsCompleted {
		oldTodo.IsCompleted = todo.IsCompleted
	}

	err = c.todoService.UpdateTodo(requests.Context(), todoId, *oldTodo)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	webResponse := model.WebResponse{
		Code:   http.StatusOK,
		Status: "updated the todo successfully",
		Data:   nil,
	}
	utils.WriteResponseBody(writer, webResponse)
}

func (c *TodoController) FindTodoByUserId(writer http.ResponseWriter, requests *http.Request) {
	userIdStr, err := c.jwtService.GetSubjectFromAccessToken(getAccessTokenFromHeaders(writer, requests))
	if err != nil {
		http.Error(writer, fmt.Sprintf("token's subject has invalid format: %v", err), http.StatusUnauthorized)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(writer, fmt.Sprintf("can not parse subject from token: %v", err), http.StatusBadRequest)
		return
	}

	todos, err := c.todoService.FindTodosByUserId(requests.Context(), userId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	webResponse := model.WebResponse{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   todos,
	}
	utils.WriteResponseBody(writer, webResponse)
}
