package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Ginのデフォルトルーターを作成する
	router := gin.Default()

	// /api以下の各ルートを定義する
	api := router.Group("/api")
	{
		api.GET("/sample", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Hello, World!",
			})
		})
	}

	// Ginサーバーをポート8080で起動する
	log.Println("Ginサーバーをポート8080で起動します")
	router.Run(":8080")
}
