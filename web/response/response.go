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

type JSONSuccessDeleteResult struct {
	Code    int    `json:"code"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
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

func SuccessDeleteResponse(c echo.Context, code int, status bool, message string) error {
	return c.JSON(code, JSONSuccessDeleteResult{
		Code:    code,
		Status:  status,
		Message: message,
	})
}

func FailResponse(c echo.Context, code int, status bool, message string) error {
	if code == http.StatusUnauthorized {
		return c.JSON(http.StatusUnauthorized, JSONUnauthorizedResult{
			Code:    code,
			Message: message,
			Status:  status,
		})
	}

	if code == http.StatusBadRequest {
		return c.JSON(http.StatusBadRequest, JSONBadRequestResult{
			Code:    code,
			Message: message,
			Status:  status,
		})
	}

	if code == http.StatusNotFound {
		return c.JSON(http.StatusNotFound, JSONBadRequestResult{
			Code:    code,
			Message: message,
			Status:  status,
		})
	}

	return nil
}
