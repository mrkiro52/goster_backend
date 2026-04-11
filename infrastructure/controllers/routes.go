package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	userCtrl *UserController,
) {

	router.GET("/_info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	// auth
	api.POST("/register", userCtrl.Register)
	api.POST("/login", userCtrl.Login)

}
