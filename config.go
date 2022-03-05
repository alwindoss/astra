package astra

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

type Config struct {
	Location      string
	DBName        string
	Port          string
	InProduction  bool
	TemplateCache map[string]*template.Template
}

func DefaultConfig(dbPath, dbName string) *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	cfg := &Config{
		DBName: dbName,
		Port:   port,
	}
	if dbPath == "" {
		homeDir, err := homedir.Dir()
		if err != nil {
			log.Printf("unable to find the home directory: %v", err)
			os.Exit(1)
		}
		cfgLoc := filepath.Join(homeDir, ".astra")
		if err = os.MkdirAll(cfgLoc, 0750); err != nil {
			log.Printf("unable to create the config folder at the location %s: %v", cfgLoc, err)
			os.Exit(1)
		}
		cfg.Location = cfgLoc

	} else {
		if err := os.MkdirAll(dbPath, 0750); err != nil {
			log.Printf("unable to create the config folder at the location %s: %v", dbPath, err)
			os.Exit(1)
		}
		cfg.Location = dbPath
	}
	log.Printf("Created and set the config location %s", cfg.Location)
	cfg.InProduction = false

	return cfg
}
