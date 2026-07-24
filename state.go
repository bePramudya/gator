package main

import (
	"github.com/bePramudya/gator/internal/config"
	"github.com/bePramudya/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
	// Config *config.Config
}
