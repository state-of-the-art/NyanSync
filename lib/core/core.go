package core

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/state-of-the-art/NyanSync/lib/api"
	"github.com/state-of-the-art/NyanSync/lib/config"
	"github.com/state-of-the-art/NyanSync/lib/gui"
	"github.com/state-of-the-art/NyanSync/lib/state"
)

func Init(configuration *config.Config) {
	// Init cfg variable
	cfg = configuration

	// Init state
	state.Init(cfg.StateFilePath)
	state.SourcesUpdate(cfg.Sources)
}

func RunHTTPServer() {
	router := gin.Default()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	api.InitV1(router)
	gui.Init(cfg.GuiPath, router)

	log.Fatal(router.Run(cfg.Endpoint.Address))
}

func ConfigUpdateSource(id string, src *config.Source) {
	cfg.SourceSet(id, src)
}

// Core configuration
var cfg = &config.Config{}
