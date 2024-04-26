package controllers

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"registration/models"
	"registration/utilities"

	"github.com/gin-gonic/gin"
	"github.com/pelicanplatform/pelican/client"
	"github.com/pelicanplatform/pelican/config"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type osdfPath struct {
	Name string `json:"name"`
}

// GET /download
// Get file
func GetFile(c *gin.Context) {
	var newPath osdfPath
	if err := c.BindJSON(&newPath); err != nil {
		return
	}
	federationDiscovURL := os.Getenv("FEDERATION_DISCOVERY_URL")
	log.Println("federationDiscovURL: ", federationDiscovURL)
	log.Println("Initializing pelican init client")
	viper.Reset()
	config.InitConfig()
	viper.Set("Federation.DiscoveryUrl", federationDiscovURL)
	err := config.InitClient()
	if err != nil {
		log.Println("Failed to init pelican client:", err)
	}
	te := client.NewTransferEngine(c)
	defer func() {
		if err := te.Shutdown(); err != nil {
			log.Println("Failure when shutting down transfer engine:", err)
		}
	}()
	project := "" // Used for condor jobs
	localDestination := "tmp"
	remoteObjectUrl, err := url.Parse(newPath.Name)
	if err != nil {
		log.Println("Failed to parse source URL:", err)
	}
	tc, err := te.NewClient()
	if err != nil {
		log.Println("Failure when creating new client:", err)
	}
	tj, err := tc.NewTransferJob(context.Background(), remoteObjectUrl, localDestination, false, false, project)
	if err != nil {
		log.Println("Failure when creating new transfer job:", err)
	}
	err = tc.Submit(tj)
	if err != nil {
		log.Println("Failure when submitting job:", err)
	}
	transferResults, err := tc.Shutdown()
	var downloaded int64 = 0
	for _, result := range transferResults {
		downloaded += result.TransferredBytes
		if err == nil && result.Error != nil {
			err = result.Error
		}
	}

	inputFile, err := os.ReadFile(localDestination)
	if err != nil {
		panic(err)
	}
	log.Println("Downloaded results:", downloaded)
	c.Header("Content-Disposition", "attachment; filename="+newPath.Name)
	c.Data(http.StatusOK, "application/data", inputFile)
	//Remove file
	err = os.Remove(localDestination)
	if err != nil {
		log.Println(err)
	}
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
