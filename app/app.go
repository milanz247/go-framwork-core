// Package app wires application configuration, persistence and routes to the core.
package app

import (
	"log"

	"ticket-sync-mvc/app/repositories"
	"ticket-sync-mvc/config"
	"ticket-sync-mvc/framework/core"
	"ticket-sync-mvc/framework/database"
	"ticket-sync-mvc/framework/env"
	"ticket-sync-mvc/routes"
)

func Run() {
	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}
	if err := repositories.Open(env.Get("DATABASE_PATH", "data/ticket-sync.db")); err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()
	core.Run(":"+env.Get("PORT", "8001"), routes.Handler())
}
