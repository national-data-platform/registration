package main

import (
	"registration/controllers"
	"registration/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db := models.SetupDB()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.POST("/download", controllers.GetFile)
	r.GET("/bestcache", controllers.GetBestCache)
	r.GET("/namespaces", controllers.GetNamespaces)
	r.GET("/datasets", controllers.GetDatasets)
	r.POST("/datasets", controllers.CreateDataset)
	r.Run()
}
