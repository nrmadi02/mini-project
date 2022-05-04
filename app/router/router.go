package router

import (
	"github.com/labstack/echo/v4"
	http4 "github.com/nrmadi02/mini-project/Favorite/delivery/http"
	repository6 "github.com/nrmadi02/mini-project/Favorite/repository"
	usecase5 "github.com/nrmadi02/mini-project/Favorite/usecase"
	http3 "github.com/nrmadi02/mini-project/enterprise/delivery/http"
	repository4 "github.com/nrmadi02/mini-project/enterprise/repository"
	usecase3 "github.com/nrmadi02/mini-project/enterprise/usecase"
	repository5 "github.com/nrmadi02/mini-project/rating/repository"
	usecase4 "github.com/nrmadi02/mini-project/rating/usecase"
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
	enterpriseRepository := repository4.NewEnterpriseRepository(db)
	ratingRepository := repository5.NewRatingRepository(db)
	favoriteRepository := repository6.NewFavoriteRepository(db)

	authUsecase := usecase.NewAuthUsecase(userRepository, roleRepository, favoriteRepository, enterpriseRepository)
	userUsecase := usecase.NewUserUsecase(userRepository, authUsecase)
	tagUsecase := usecase2.NewTagUsecase(tagRepository)
	enterpriseUsecase := usecase3.NewEnterpriseUsecase(enterpriseRepository, tagRepository, userRepository)
	ratingUsecase := usecase4.NewRatingUsecase(userRepository, enterpriseRepository, ratingRepository)
	favoriteUsecase := usecase5.NewFavoriteUsecase(enterpriseRepository, favoriteRepository)

	authController := http.NewAuthController(authUsecase)
	userController := http.NewUserController(authUsecase, ratingUsecase)
	adminController := http.NewAdminController(authUsecase, userUsecase)
	tagController := http2.NewTagController(authUsecase, tagUsecase)
	enterpriseController := http3.NewEnterpriseController(authUsecase, enterpriseUsecase, ratingUsecase)
	favoriteController := http4.NewFavoriteController(favoriteUsecase, authUsecase, ratingUsecase)

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

	//enterprise endpoints
	c.POST("/api/v1/enterprise", enterpriseController.CreateNewEnterprise, authMiddleware)
	c.PUT("/api/v1/enterprise/:id/status", enterpriseController.UpdateStatusEnterprise, authMiddleware)
	c.GET("/api/v1/enterprises/:status", enterpriseController.GetEnterpriseByStatus, authMiddleware)
	c.PUT("/api/v1/enterprise/:id", enterpriseController.UpdateEnterpriseByID, authMiddleware)
	c.GET("/api/v1/enterprises", enterpriseController.GetAllEnterprises, authMiddleware)
	c.DELETE("/api/v1/enterprise/:id", enterpriseController.DeleteEnterpriseByID, authMiddleware)
	c.GET("/api/v1/enterprise/:id", enterpriseController.GetDetailEnterpriseByID, authMiddleware)
	c.GET("/api/v1/enterprise/:id/distance", enterpriseController.GetDistance, authMiddleware)
	c.POST("/api/v1/enterprise/:id/rating", enterpriseController.AddNewRanting, authMiddleware)
	c.GET("/api/v1/enterprise/:id/rating/user/:userid", enterpriseController.CekRatingUser, authMiddleware)
	c.DELETE("/api/v1/enterprise/:id/rating/user/:userid", enterpriseController.DeleteRatingUser, authMiddleware)
	c.PUT("/api/v1/enterprise/:id/rating/user/:userid", enterpriseController.UpdateRating, authMiddleware)

	//favorite endpoint
	c.POST("/api/v1/favorite", favoriteController.AddFavoriteEnterprise, authMiddleware)
	c.DELETE("/api/v1/favorite", favoriteController.RemoveFavoriteEnterprise, authMiddleware)
	c.GET("/api/v1/favorite", favoriteController.GetDetailFavoriteEnterprise, authMiddleware)
}
