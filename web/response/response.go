package response

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type JSONSuccessResult struct {
	Code    int         `json:"code"`
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type JSONBadRequestResult struct {
	Code    int    `json:"code"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type JSONUnauthorizedResult struct {
	Code    int    `json:"code"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func SuccessResponse(c echo.Context, code int, status bool, message string, data interface{}) error {
	return c.JSON(code, JSONSuccessResult{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func FailResponse(c echo.Context, code int, status bool, message string) error {
	if code == http.StatusUnauthorized {
		return c.JSON(http.StatusOK, JSONUnauthorizedResult{
			Code:    code,
			Message: message,
			Status:  status,
		})
	}

	if code == http.StatusBadRequest {
		return c.JSON(http.StatusOK, JSONBadRequestResult{
			Code:    code,
			Message: message,
			Status:  status,
		})
	}

	return nil
}
