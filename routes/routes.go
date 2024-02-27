package routes

import (
	"go-chi/controller/events"
	"go-chi/controller/tasks"
	"go-chi/controller/users"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(router fiber.Router) {
	router.Get("/", users.AllUsers)
	// router.Get("/:id", users.GetUser)
	// router.Post("/", users.PostUser)
	// router.Patch("/:id", users.UpdateUser)
}

func EventRoute(router fiber.Router) {
	router.Get("/", events.GetEvents)
	router.Get("/:id", events.GetEvent)
	router.Post("/", events.PostEvent)
	router.Delete("/:id", events.DeleteEvent)
	router.Patch("/:id", events.UpdateEvent)
}

func TaskRoute(router fiber.Router) {
	router.Get("/", tasks.GetTasks)
	router.Get("/:id", tasks.GetTask)
	router.Post("/", tasks.PostTask)
	router.Delete("/:id", tasks.DeleteTask)
	router.Patch("/:id", tasks.UpdateTask)
}
