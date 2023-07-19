package main

import (
	"chat/config"
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

	config.ConnectDatabase(loadedConfig.DBDriver, loadedConfig.DBSource, "")
	r := gin.Default()
	// Enable CORS for requests from localhost
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.GET("/", func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusOK, config.ResponseData("success", "Server is up", nil))
	})

	err = r.Run(":3000")
	if err != nil {
		fmt.Printf("cannot run server: %w", err)
		return
	}
}
