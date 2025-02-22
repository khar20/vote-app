package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Cors Config
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:8000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}

	router.Use(cors.New(config))

	// Rutas
	router.GET("/vote", GetVotesHandler())
	router.GET("/public-key", GetPublicKey())
	router.POST("/vote", PostVoteHandler())
	router.GET("/vote/check", CheckChainHandler())
	router.GET("/vote/tally", TallyVotesHandler())

	return router
}
