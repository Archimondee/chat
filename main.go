package main

import (
	"chat/app/http/controllers"
	"chat/app/repository"
	"chat/config"
	"chat/utils"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("cannot load loadedConfig: %w", err)
		return
	}

	utils.ConnectDatabase(loadedConfig.DBDriver, loadedConfig.DBSource, "")

	ctx := context.TODO()
	db := utils.DB
	UserRepository := repository.NewUserRepositoryImpl(ctx, db)

	AuthRepository := repository.NewAuthRepository(ctx, db, UserRepository)
	AuthController := controllers.NewAuthController(AuthRepository, ctx)

	r := gin.Default()
	// Enable CORS for requests from localhost
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.GET("/", func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusOK, utils.ResponseData("success", "Server is up", nil))
	})

	apiv1 := r.Group("/api/v1")
	{
		auth := apiv1.Group("/auth")
		{
			auth.POST("/signin", AuthController.SigninUser)
			auth.POST("/signup", AuthController.SignupUser)
		}
	}

	err = r.Run(":3000")
	if err != nil {
		fmt.Printf("cannot run server: %w", err)
		return
	}
}
