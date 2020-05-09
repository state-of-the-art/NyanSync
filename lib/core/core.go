package core

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/state-of-the-art/NyanSync/lib/api"
	"github.com/state-of-the-art/NyanSync/lib/config"
	"github.com/state-of-the-art/NyanSync/lib/gui"
	"github.com/state-of-the-art/NyanSync/lib/state"
)

func Init() {
	config.Load()

	// Init state
	state.Init()
}

func RunHTTPServer() {
	router := gin.Default()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	api.InitV1(router)
	gui.Init(config.Cfg().GuiPath, router)

	log.Fatal(router.Run(config.Cfg().Endpoint.Address))
}
