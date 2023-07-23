package main

import (
	"chat/app/http/controllers"
	_ "chat/app/interfaces"
	"chat/app/middlewares"
	"chat/app/repository"
	"chat/app/ws"
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
	utils.ConnectAmqp(loadedConfig.AmqpURL, loadedConfig.AmqpQueue, loadedConfig.AmqpRouting, loadedConfig.AmqpExchange)

	ctx := context.TODO()
	db := utils.DB
	UserRepository := repository.NewUserRepositoryImpl(ctx, db)
	UserController := controllers.NewUserController(UserRepository, ctx)

	AuthRepository := repository.NewAuthRepository(ctx, db, UserRepository)
	AuthController := controllers.NewAuthController(AuthRepository, ctx)

	MessageRepository := repository.NewMessageRepositoryImpl(ctx, db)
	MessageController := controllers.NewMessageController(ctx, MessageRepository)

	r := gin.Default()
	// Enable CORS for requests from localhost
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"*"}
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

		user := apiv1.Group("/users").Use(middlewares.AuthMiddleware(UserRepository))
		{
			user.GET("/", UserController.GetAllUser)
		}

		message := apiv1.Group("/chat")
		{
			message.GET("/", MessageController.ReadMessage)
		}
	}

	server := ws.NewWebsocketServer(UserRepository)
	go server.Run()

	r.GET("/message", func(c *gin.Context) {
		ws.ServeWebsocket(server, c.Writer, c.Request, MessageRepository)
	})

	err = r.Run(":3000")
	if err != nil {
		fmt.Printf("cannot run server: %w", err)
		return
	}
}
