package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"vk-inter-test-go/internal/db/repo"
	"vk-inter-test-go/internal/io/ioutils"
	"vk-inter-test-go/internal/io/models"
	"vk-inter-test-go/internal/utils"
)

// CreateUser creates a new user.
//
// @Summary Creates a new user
// @Description Creates a new user with the data provided in the request body.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param body body repo.User true "User data"
// @Success 200 {object} models.TokenResponse "Generated token"
// @Failure 400 {object} models.ErrorResponse "Invalid data format or unsupported characters in the username"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/create/user [post]
func (c *Controller) CreateUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		ioutils.HandleInvalidMethodResponse(w, req.Method)
		return
	}

	var user repo.User
	err := ioutils.DecodeRequestBody(req, &user)
	if err != nil || !ioutils.UserJsonValidate(user) {
		ioutils.HandleInvalidJson(w)
		return
	}

	var answer interface{}

	var cleanInput bool
	user.Login, cleanInput = utils.Sanitize(user.Login)
	if !cleanInput {
		answer = models.ErrorResponse{
			Error: "имя содержит непотдерживаемые символы",
		}
	} else {
		token, err := c.Bl.CreateUser(user)
		if err != nil {
			c.logger.Info("err", zap.Error(err))
			answer = models.ErrorResponse{
				Error: "user is exist: '" + user.Login + "'",
			}
		} else {
			answer = models.TokenResponse{Bearer: token}
		}
	}

	ioutils.RespJson(w, answer)
}

// AuthUser authenticates a user.
//
// @Summary Authenticates a user
// @Description Authenticates a user with the provided data in the request body.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param body body repo.User true "User data"
// @Success 200 {object} models.TokenResponse "Generated token"
// @Failure 400 {object} models.ErrorResponse "Invalid data format or unsupported characters in the username"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/login [get]
func (c *Controller) AuthUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		ioutils.HandleInvalidMethodResponse(w, req.Method)
		return
	}
	var user repo.User
	err := ioutils.DecodeRequestBody(req, &user)
	if err != nil || !ioutils.UserJsonValidate(user) {
		ioutils.HandleInvalidJson(w)
		return
	}

	var answer interface{}

	var cleanInput bool
	user.Login, cleanInput = utils.Sanitize(user.Login)
	if !cleanInput {
		answer = models.ErrorResponse{
			Error: "имя содержит непотдерживаемые символы",
		}
	} else {
		token, err := c.Bl.AuthUser(user)
		if err != nil {
			c.logger.Info("err", zap.Error(err))
			answer = models.ErrorResponse{
				Error: "wrong pass",
			}
		} else {
			answer = models.TokenResponse{Bearer: token}
		}
	}

	ioutils.RespJson(w, answer)

}
