package main

import (
	"github.com/alwindoss/astra"
	"github.com/alwindoss/astra/internal/server"
)

func main() {
	cfg := astra.DefaultConfig()

	server.Run(cfg)
}
