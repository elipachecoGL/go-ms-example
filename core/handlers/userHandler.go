package handlers

import (
	"fmt"
	"net/http"
	"user/app/utils/response"
	"user/core/entities"

	"github.com/gofiber/fiber/v2"
)

//NewUserHandler factory method for creating a method handler to users
func NewUserHandler(repository entities.UserRepository) MethodHandlers {
	return userHandler{store: repository}
}

type userHandler struct {
	store entities.UserRepository
}

type BasicUser struct {
	Name  string
	Email string
}

func (handler userHandler) RegisterMethods(app *fiber.App) {
	app.Get("/api/v1/users", handler.getUsers)
	app.Get("/api/v1/users/email", handler.getUser)
	app.Post("/api/v1/users", handler.newUser)
	app.Put("/api/v1/users/:id", handler.updateUser)
	app.Delete("/api/v1/users/:id", handler.deleteUser)
}

func (handler userHandler) getUsers(context *fiber.Ctx) error {
	users, err := handler.store.Users()

	if err != nil {
		return handler.send(nil, &response.ResponseError{StatusCode: 500, Message: err.Error()}, context)
	}

	if len(users) == 0 {
		return handler.send([]entities.User{}, nil, context)
	}

	return handler.send(users, nil, context)
}

func (handler userHandler) getUser(context *fiber.Ctx) error {
	userEmail := context.Query("address")

	if userEmail == "" {
		return handler.send(nil, &response.ResponseError{StatusCode: http.StatusBadRequest, Message: "address is not present on url as a query param"}, context)
	}

	user, err := handler.store.User(userEmail)

	if err != nil {
		return handler.send(nil, &response.ResponseError{StatusCode: 404, Message: err.Error()}, context)
	}

	return handler.send(user, nil, context)
}

func (handler userHandler) newUser(context *fiber.Ctx) error {
	user := new(entities.User)

	if err := context.BodyParser(user); err != nil {
		context.SendStatus(503)
		return nil
	}

	return context.JSON(fiber.Map{
		"result": fmt.Sprintf("Welcome %s!", user.Nickname),
	})
}

func (handler userHandler) updateUser(context *fiber.Ctx) error {
	userID := context.Params("id")

	if userID == "" {
		context.SendStatus(503)
		return nil
	}

	return context.JSON(fiber.Map{
		"result": fmt.Sprintf("User %s updated!", userID),
	})
}

func (handler userHandler) deleteUser(context *fiber.Ctx) error {
	userID := context.Params("id")

	if userID == "" {
		context.SendStatus(503)
		return nil
	}

	return context.JSON(fiber.Map{
		"result": fmt.Sprintf("User %s deleted!", userID),
	})
}

func (handler userHandler) send(value interface{}, err *response.ResponseError, context *fiber.Ctx) error {
	if err != nil {
		return response.MakeJSON(response.Fail, nil, err, context)
	}

	return response.MakeJSON(response.Success, &value, nil, context)
}