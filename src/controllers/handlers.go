package controllers

import (
	"net/http"
	"registration/models"
	"registration/utilities"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
)

type osdfPath struct {
    Name     string  `json:"name"`
}

// GET /download
// Get file
func GetFile(c *gin.Context) {
	var newPath osdfPath
	if err := c.BindJSON(&newPath); err != nil {
        return
    }
	fileName := newPath.Name
	inputFile, err := ioutil.ReadFile(fileName)
    if err != nil {
        panic(err)
    }
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, "application/data", inputFile)
}

// GET /bestcache
// Get best cache based on geolocation in osdf
func GetBestCache(c *gin.Context) {
	bestcache, err := utilities.GetBestCache()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": bestcache})
}

// GET /namespaces
// Get all namespaces in osdf
func GetNamespaces(c *gin.Context) {
	ns, err := utilities.GetNamespaces()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": ns})
}

// GET /datasets
// Get all datasets
func GetDatasets(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataset []models.Dataset
	db.Find(&dataset)
	c.JSON(http.StatusOK, gin.H{"data": dataset})
}

// POST /dataset
// Create new dataset
func CreateDataset(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input models.CreateDatasetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dataset := models.Dataset{Name: input.Name, Owner: input.Owner, Content: input.Content}
	db.Create(&dataset)
	c.JSON(http.StatusOK, gin.H{"data": dataset})
}
