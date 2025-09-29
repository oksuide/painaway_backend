package response

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//TODO: Переписать фронт под один ответов

type ErrorResponse struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string, logger *zap.Logger) {
	var errMsg string
	switch statusCode {
	case 400:
		errMsg = "BadRequestException"
		logger.Warn(message, zap.Int("status", statusCode), zap.String("error_type", errMsg), zap.String("path", c.Request.URL.Path))
	case 401:
		errMsg = "UnauthorizedException"
		logger.Warn(message, zap.Int("status", statusCode), zap.String("error_type", errMsg), zap.String("path", c.Request.URL.Path))
	case 402:
		errMsg = "ForbiddenException"
		logger.Warn(message, zap.Int("status", statusCode), zap.String("error_type", errMsg), zap.String("path", c.Request.URL.Path))
	default:
		errMsg = "InternalServerError"
		logger.Error(message, zap.Int("status", statusCode), zap.String("error_type", errMsg), zap.String("path", c.Request.URL.Path))
	}

	c.AbortWithStatusJSON(statusCode, ErrorResponse{
		Status:  statusCode,
		Error:   errMsg,
		Message: message,
	})
}

//TODO: ...
// type SuccessResponse struct {
// 	Status  int         `json:"status"`
// 	Data    interface{} `json:"data"`
// 	Message string      `json:"message"`
// }

// func NewSuccessResponse(c *gin.Context, data interface{}, message string) {
// 	if message == "" {
// 		message = "OK"
// 	}

// 	c.JSON(http.StatusOK, SuccessResponse{
// 		Status:  http.StatusOK,
// 		Data:    data,
// 		Message: message,
// 	})
// }
