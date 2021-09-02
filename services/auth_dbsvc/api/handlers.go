// package api

// import (
// 	"encoding/json"
// 	"net/http"

// 	"deblasis.net/space-traffic-control/services/auth_dbsvc/api/model"
// 	"deblasis.net/space-traffic-control/services/auth_dbsvc/service"
// 	"github.com/labstack/echo/v4"
// 	"github.com/sirupsen/logrus"
// )

// var internalErrorMessage = NewErrorMessage("internal server error")

// type Server struct {
// 	userManager service.UserManager
// 	logger      *logrus.Logger
// }

// func NewServer(userManager service.UserManager, logger *logrus.Logger) *Server {
// 	return &Server{
// 		userManager: userManager,
// 		logger:      logger,
// 	}
// }

// func (s *Server) CreateUser(ctx echo.Context) error {
// 	context := ctx.Request().Context()
// 	var createUserRequest = new(model.CreateUserRequest)

// 	body := ctx.Request().Body
// 	if body == nil {
// 		s.logger.Warn("Body is nil")
// 		return ctx.JSON(http.StatusBadRequest, NewErrorMessage("Body is nil"))
// 	}

// 	decoder := json.NewDecoder(body)
// 	err := decoder.Decode(&createUserRequest)
// 	if err != nil {
// 		s.logger.WithContext(context).Error(err)
// 		return ctx.JSON(http.StatusBadRequest, NewErrorMessage("invalid input: cannot decode"))
// 	}

// 	user := model.User{
// 		ID:       0,
// 		Username: createUserRequest.Username,
// 		Password: createUserRequest.Password,
// 		RoleID:   createUserRequest.RoleID,
// 	}
// 	err = s.userManager.CreateUser(context, &user)
// 	if err != nil {
// 		s.logger.WithContext(context).Error(err)
// 		return ctx.JSON(http.StatusInternalServerError, internalErrorMessage)
// 	}

// 	//accountResponse := model.AccountResponse(account)
// 	return ctx.JSON(http.StatusOK, user)
// }
