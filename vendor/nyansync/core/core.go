package core

import (
    "nyansync/state"
    "nyansync/config"
)

func Init(configuration *config.Config) {
    // Init cfg variable
    cfg = configuration

    // Init state
    state.Init(cfg.StateFilePath)
}

// Core configuration
var cfg = &config.Config{}
