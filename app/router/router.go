package router

import (
	"github.com/labstack/echo/v4"
	repository2 "github.com/nrmadi02/mini-project/role/repository"
	"github.com/nrmadi02/mini-project/user/delivery/http"
	mid "github.com/nrmadi02/mini-project/user/delivery/http/middleware"
	"github.com/nrmadi02/mini-project/user/repository"
	"github.com/nrmadi02/mini-project/user/usecase"
	"gorm.io/gorm"
)

func SetupRouter(c *echo.Echo, db *gorm.DB) {
	authMiddleware := mid.NewGoMiddleware().AuthMiddleware()

	userRepository := repository.NewUserRepository(db)
	roleRepository := repository2.NewRoleRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepository, roleRepository)
	userUsecase := usecase.NewUserUsecase(userRepository, authUsecase)
	authController := http.NewAuthController(authUsecase)
	userController := http.NewUserController(authUsecase)
	adminController := http.NewAdminController(authUsecase, userUsecase)

	// ADMIN ENDPOINTS:
	c.GET("/api/v1/admin/users", adminController.GetUserList, authMiddleware)

	// Auth Endpoints (User)
	c.POST("/api/v1/register", authController.Register)
	c.POST("/api/v1/login", authController.Login)

	//user endpoints
	c.GET("/api/v1/user", userController.User, authMiddleware)

}
