package main

import (
	"fmt"
	"go-chi/db"
	"go-chi/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	r := fiber.New()

	r.Use(
		cors.New(),
	)

	db.ConnectDb()

	// Group routes
	api := r.Group("/api")
	api.Route("/user", routes.UserRoute)
	api.Route("/events", routes.EventRoute)
	api.Route("/tasks", routes.TaskRoute)

	port := "3001"

	fmt.Println("Listening of port " + port)
	log.Fatal(r.Listen(":" + port))
}
