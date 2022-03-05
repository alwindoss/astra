package main

import (
	"flag"

	"github.com/alwindoss/astra"
	"github.com/alwindoss/astra/internal/server"
)

var (
	dbName string
	dbPath string
)

func init() {
	flag.StringVar(&dbName, "name", "astra.db", "-name=dbname.db (default is astra.db)")
	flag.StringVar(&dbPath, "loc", "", "-loc=/opt/astra/db/ (default is home directory)")
}

func main() {
	flag.Parse()
	cfg := astra.DefaultConfig(dbPath, dbName)

	server.Run(cfg)
}
