package router

import (
	"github.com/labstack/echo/v4"
	http3 "github.com/nrmadi02/mini-project/internal/enterprise/delivery/http"
	repository4 "github.com/nrmadi02/mini-project/internal/enterprise/repository"
	usecase3 "github.com/nrmadi02/mini-project/internal/enterprise/usecase"
	http4 "github.com/nrmadi02/mini-project/internal/favorite/delivery/http"
	repository6 "github.com/nrmadi02/mini-project/internal/favorite/repository"
	usecase5 "github.com/nrmadi02/mini-project/internal/favorite/usecase"
	repository5 "github.com/nrmadi02/mini-project/internal/rating/repository"
	usecase4 "github.com/nrmadi02/mini-project/internal/rating/usecase"
	http5 "github.com/nrmadi02/mini-project/internal/review/delivery/http"
	repository7 "github.com/nrmadi02/mini-project/internal/review/repository"
	usecase6 "github.com/nrmadi02/mini-project/internal/review/usecase"
	repository2 "github.com/nrmadi02/mini-project/internal/role/repository"
	http2 "github.com/nrmadi02/mini-project/internal/tag/delivery/http"
	repository3 "github.com/nrmadi02/mini-project/internal/tag/repository"
	usecase2 "github.com/nrmadi02/mini-project/internal/tag/usecase"
	http6 "github.com/nrmadi02/mini-project/internal/user/delivery/http"
	mid "github.com/nrmadi02/mini-project/internal/user/delivery/http/middleware"
	"github.com/nrmadi02/mini-project/internal/user/repository"
	usecase7 "github.com/nrmadi02/mini-project/internal/user/usecase"
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
	reviewRepository := repository7.NewReviewRepository(db)

	authUsecase := usecase7.NewAuthUsecase(userRepository, roleRepository, favoriteRepository, enterpriseRepository)
	userUsecase := usecase7.NewUserUsecase(userRepository, authUsecase)
	tagUsecase := usecase2.NewTagUsecase(tagRepository)
	enterpriseUsecase := usecase3.NewEnterpriseUsecase(enterpriseRepository, tagRepository, userRepository)
	ratingUsecase := usecase4.NewRatingUsecase(userRepository, enterpriseRepository, ratingRepository)
	favoriteUsecase := usecase5.NewFavoriteUsecase(enterpriseRepository, favoriteRepository)
	reviewUsecase := usecase6.NewReviewUsecase(enterpriseRepository, userRepository, reviewRepository)

	authController := http6.NewAuthController(authUsecase)
	userController := http6.NewUserController(authUsecase, ratingUsecase)
	adminController := http6.NewAdminController(authUsecase, userUsecase)
	tagController := http2.NewTagController(authUsecase, tagUsecase)
	enterpriseController := http3.NewEnterpriseController(authUsecase, enterpriseUsecase, ratingUsecase)
	favoriteController := http4.NewFavoriteController(favoriteUsecase, authUsecase, ratingUsecase)
	reviewController := http5.NewReviewController(reviewUsecase, enterpriseUsecase, authUsecase)

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

	//favorite endpoints
	c.POST("/api/v1/favorite", favoriteController.AddFavoriteEnterprise, authMiddleware)
	c.DELETE("/api/v1/favorite", favoriteController.RemoveFavoriteEnterprise, authMiddleware)
	c.GET("/api/v1/favorite", favoriteController.GetDetailFavoriteEnterprise, authMiddleware)

	//review endpoints
	c.POST("/api/v1/review/enterprise/:id", reviewController.AddReviewEnterprise, authMiddleware)
	c.GET("/api/v1/review/enterprise/:id", reviewController.GetListReviewByEnterpriseID, authMiddleware)
	c.PUT("/api/v1/review/enterprise/:id", reviewController.UpdateReviewEnterprise, authMiddleware)
	c.DELETE("/api/v1/review/enterprise/:id", reviewController.DeleteReviewEnterprise, authMiddleware)
	c.GET("/api/v1/review/:id", reviewController.GetDetailReviewByID, authMiddleware)
}
