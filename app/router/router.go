package router

import (
	"github.com/labstack/echo/v4"
	repository2 "github.com/nrmadi02/mini-project/role/repository"
	http2 "github.com/nrmadi02/mini-project/tag/delivery/http"
	repository3 "github.com/nrmadi02/mini-project/tag/repository"
	usecase2 "github.com/nrmadi02/mini-project/tag/usecase"
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
	tagRepository := repository3.NewTagRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepository, roleRepository)
	userUsecase := usecase.NewUserUsecase(userRepository, authUsecase)
	tagUsecase := usecase2.NewTagUsecase(tagRepository)
	authController := http.NewAuthController(authUsecase)
	userController := http.NewUserController(authUsecase)
	adminController := http.NewAdminController(authUsecase, userUsecase)
	tagController := http2.NewTagController(authUsecase, tagUsecase)

	// Auth Endpoints (User)
	c.POST("/api/v1/register", authController.Register)
	c.POST("/api/v1/login", authController.Login)

	//user endpoints
	c.GET("/api/v1/users", adminController.GetUserList, authMiddleware)
	c.GET("/api/v1/user", userController.User, authMiddleware)

	//tag endpoints
	c.GET("/api/v1/tags", tagController.GetTagsList, authMiddleware)
	c.DELETE("/api/v1/tag/:id", tagController.DeleteTag, authMiddleware)
	c.POST("/api/v1/tag", tagController.CreateTag, authMiddleware)
}
