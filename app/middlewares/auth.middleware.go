package middlewares

import (
	"chat/app/interfaces"
	"chat/app/models/entity"
	"chat/config"
	"chat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"net/http"
	"strings"
)

func AuthMiddleware(userRepository interfaces.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var access_token string
		authorizationHeader := ctx.GetHeader("Authorization")

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		}

		if access_token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResponseData("error", "you are not logged in", nil))
			return
		}

		config, _ := config.LoadConfig(".")
		sub, err := utils.ValidateToken(access_token, config.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResponseData("error", "you are not logged in", err))
			return
		}
		fmt.Println("data", sub)
		var user *entity.User
		//data := fmt.Sprint(sub)
		//jsonData := []byte(data)
		//fmt.Println("data", jsonData)
		dataSub, err := json.Marshal(sub)
		if err != nil {
			fmt.Println("Error : ", err)
		}

		if err := json.Unmarshal(dataSub, &user); err != nil {
			fmt.Println("Error:", err)
			return
		}

		user, err = userRepository.FindUserById(user.Id)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResponseData("error", "you are not logged in", err))
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()

	}
}
