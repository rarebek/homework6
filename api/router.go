package api

import (
	mid "EXAM3/api-gateway/api/middleware"

	"github.com/casbin/casbin/util"
	casbinN "github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "EXAM3/api-gateway/api/docs"
	v1 "EXAM3/api-gateway/api/handlers/v1"
	"EXAM3/api-gateway/api/handlers/v1/tokens"
	"EXAM3/api-gateway/config"
	"EXAM3/api-gateway/pkg/logger"
	"EXAM3/api-gateway/services"
	"EXAM3/api-gateway/storage/repo"
)

type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Reds           repo.RedisStorageI
}

// @title Welcome to User-Product service
// @version 1.0
// @description This code is written by Nodirbek in third mont exam GOLANG
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func New(option Option) *gin.Engine {
	casbinEnforcer, err := casbinN.NewEnforcer(option.Conf.CasbinConfigPath, option.Conf.AuthCSVPath)
	if err != nil {
		option.Logger.Error("cannot create a new enforcer", logger.Error(err))
	}
	_ = casbinEnforcer.LoadPolicy()

	casbinEnforcer.GetRoleManager().AddMatchingFunc("keyMatch", util.KeyMatch)
	casbinEnforcer.GetRoleManager().AddMatchingFunc("keyMatch3", util.KeyMatch3)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	jwtHandle := tokens.JWTHandler{
		SignInKey: option.Conf.SigningKey,
		Log:       option.Logger,
	}

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		Reds:           option.Reds,
		Casbin:         casbinEnforcer,
		JWTHandler:     jwtHandle,
	})

	api := router.Group("/v1")
	api.Use(mid.NewAuth(casbinEnforcer, option.Conf))

	api.POST("/user/register", handlerV1.RegisterUser)
	api.POST("/user/verify/:email/:code", handlerV1.Verify)
	api.POST("/user/login/:email/:password", handlerV1.Login)
	api.POST("/user/create", handlerV1.CreateUser)
	api.DELETE("/user/delete/:id", handlerV1.DeleteUser)
	api.GET("/user/getall/:page/:limit", handlerV1.GetAllUsers)
	api.PUT("/user/update/:id", handlerV1.UpdateUser)

	api.POST("/product/create", handlerV1.CreateProduct)
	api.PUT("/product/update/:id", handlerV1.UpdateProduct)
	api.GET("/product/get/:id", handlerV1.GetProductById)
	api.DELETE("/product/delete/:id", handlerV1.DeleteProduct)
	// api.GET("/products/get/:id", handlerV1.GetPurchasedProductsByUserId)
	api.GET("/product/:page/:limit", handlerV1.ListProducts)
	// api.POST("/product/buy", handlerV1.BuyProduct)

	api.GET("/admin/:username/:password", handlerV1.GenerateAccessTokenForAdmin)

	url := ginSwagger.URL("swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
