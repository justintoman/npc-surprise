package main

import (
	"github.com/justintoman/npc-surprise/pkg/db"
	"github.com/justintoman/npc-surprise/pkg/router"
)

func main() {
	config := LoadConfig()
	db := db.New(config.DatabaseURL, config.ApiKey)
	r := router.New(db, config.AdminKey)
	r.Run() // listen and serve on 0.0.0.0:8080
}
