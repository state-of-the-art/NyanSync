package core

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/state-of-the-art/NyanSync/lib/api"
	"github.com/state-of-the-art/NyanSync/lib/gui"
	"github.com/state-of-the-art/NyanSync/lib/config"
	"github.com/state-of-the-art/NyanSync/lib/state"
)

func Init(configuration *config.Config) {
	// Init cfg variable
	cfg = configuration

	// Init state
	state.Init(cfg.StateFilePath)
}

func RunHTTPServer() {
	router := gin.Default()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	gui.Init(cfg.GuiPath, router)
	api.InitV1(router)

	log.Fatal(router.Run(cfg.Endpoint.Address))
}

// Core configuration
var cfg = &config.Config{}
