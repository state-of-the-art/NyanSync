/**
 * Copyright (C) 2020, State Of The Art https://www.state-of-the-art.io/
 */

package main

import (
	"fmt"

	"github.com/state-of-the-art/NyanSync/lib/config"
	"github.com/state-of-the-art/NyanSync/lib/core"
)

func main() {
	cfg := config.Load()
	fmt.Printf("%+v\n", cfg)

	core.Init(cfg)
	core.RunHTTPServer()
}
