package astra

import (
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/mitchellh/go-homedir"
)

type Config struct {
	Location      string
	Port          string
	InProduction  bool
	TemplateCache map[string]*template.Template
}

func DefaultConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Printf("unable to find the home directory: %v", err)
		os.Exit(1)
	}
	cfgLoc := filepath.Join(homeDir, ".astra")
	if err = os.MkdirAll(cfgLoc, os.ModeDir); err != nil {
		log.Printf("unable to create the config folder at the location %s: %v", cfgLoc, err)
		os.Exit(1)
	}
	cfg := &Config{
		Location: cfgLoc,
		Port:     port,
	}
	cfg.InProduction = false

	return cfg
}
