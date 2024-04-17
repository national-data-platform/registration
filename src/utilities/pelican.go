package utilities

import (
	"log"
	"os"

	"github.com/pelicanplatform/pelican/client"
	"github.com/pelicanplatform/pelican/config"
	"github.com/pelicanplatform/pelican/namespaces"
	"github.com/spf13/viper"

	"registration/models"
)

func GetNamespaces() ([]models.Namespaces, error) {
	os.Setenv("PELICAN_NAMESPACE_URL", "https://topology.opensciencegrid.org/osdf/namespaces")
	viper.Reset()
	config.SetPreferredPrefix("OSDF") //deaults to PELICAN
	prefix := config.GetPreferredPrefix()
	log.Println("Prefix: ", prefix)
	err := config.InitClient()
	if err != nil {
		log.Println("Failed to init config client: ", err)
	}
	osdfNS, err := namespaces.GetNamespaces()
	if err != nil {
		log.Println("Failed to get namespaces: ", err)
	}
	allNs := []models.Namespaces{}
	for _, ns := range osdfNS {
		iNS := models.Namespaces{Name: ns.Path}
		allNs = append(allNs, iNS)
	}
	return allNs, nil
}

func GetBestCache() ([]string, error) {
	cacheListName := "xroot"
	// var bestCaches []string
	bestCaches, err := client.GetBestCache(cacheListName)
	if err != nil {
		log.Println("Failed to get best caches:", err)
	}
	return bestCaches, nil
}