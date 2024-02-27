package events

import (
	"database/sql"
	"fmt"
	"go-chi/db"
	"go-chi/models"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

var DB = db.Db
var title string
var id uint

func GetEvents(c *fiber.Ctx) error {
	var events []models.Event

	// Preload to load all the events associated to the authors
	// result := db.Db.Preload("Tasks").Select("Title").Find(&events)

	// if result.Error != nil {
	// 	return c.Status(http.StatusNotFound).SendString(result.Error.Error())
	// }

	// var title string
	// var id int

	rows, err := db.Db.Query("SELECT title, id  FROM events")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&title, &id)

		if err != nil {
			log.Fatal(err)
		}

		events = append(events, models.Event{ID: id, Title: title})
	}

	return c.Status(http.StatusOK).JSON(events)
}

func GetEvent(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Preload to load all the events associated to the authors
	err := db.Db.QueryRow("SELECT id, title FROM events WHERE id = $1", idParam).Scan(&id, &title)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).SendString(err.Error())
		}
		return c.Status(http.StatusNotFound).SendString(err.Error())
	}

	var event = models.Event{ID: id, Title: title}

	return c.Status(http.StatusOK).JSON(event)
}

func PostEvent(c *fiber.Ctx) error {
	var event models.Event

	if err := c.BodyParser(&event); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid request body: %s", err.Error())})
	}

	query := `INSERT INTO events (title) VALUES ($1) RETURNING id, title`

	err := db.Db.QueryRow(query, event.Title).Scan(&id, &title)

	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	insertedEvent := models.Event{ID: id, Title: title}

	return c.Status(http.StatusOK).JSON(fiber.Map{"messaage": "Event created successfully!", "event": insertedEvent})
}

func DeleteEvent(c *fiber.Ctx) error {
	idParam := c.Params("id")

	query := `DELETE FROM events WHERE id = $1 RETURNING id, title`

	// Delete event from db
	err := db.Db.QueryRow(query, idParam).Scan(&id, &title)

	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	deletedEvent := models.Event{ID: id, Title: title}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Event deleted successfully", "event": deletedEvent})
}

func UpdateEvent(c *fiber.Ctx) error {
	idParam := c.Params("id")

	var event models.Event

	// Bind JSON body to the User Struct for partial update
	if err := c.BodyParser(&event); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid request body: %s", err.Error())})
	}

	query := `UPDATE events SET title = $1 WHERE id = $2 RETURNING id, title`

	// Update the event in the database
	db.Db.QueryRow(query, event.Title, idParam).Scan(&id, &title)

	updatedEvent := models.Event{ID: id, Title: title}

	// Return updated event to the client
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Event updated successfully!", "event": updatedEvent})
}
