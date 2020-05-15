/**
 * Copyright (C) 2020, State Of The Art https://www.state-of-the-art.io/
 */

package main

import (
	"github.com/state-of-the-art/NyanSync/lib/core"
	_ "github.com/state-of-the-art/NyanSync/lib/processors/uri_file"
)

func main() {
	core.Init()
	core.RunHTTPServer()
}
