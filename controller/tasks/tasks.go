package tasks

import (
	"database/sql"
	"fmt"
	"go-chi/db"
	"go-chi/models"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var (
	id       uint
	title    string
	isDone   bool
	eventId  uint
	taskType string
)

func GetTasks(c *fiber.Ctx) error {
	var taskList []models.Task

	eventId := c.Query("eventId")
	taskType := c.Query("taskType")

	// Construct the base SQL query
	sqlQuery := `SELECT id, title, is_done, event_id, task_type FROM tasks`

	queryParams := map[string]interface{}{}

	if eventId != "" {
		queryParams["event_id"] = eventId
	}
	if taskType != "" {
		queryParams["task_type"] = taskType
	}

	// Check if any query parameters are provided
	if len(queryParams) > 0 {
		// Construct the WHERE clause based on the provided query parameters
		whereClause := " WHERE"
		idx := 1
		for key, value := range queryParams {
			if idx > 1 {
				whereClause += " AND"
			}

			// Check if field type is string or number
			if key == "event_id" {
				whereClause += fmt.Sprintf(" %s = %s", key, value)
			} else if key == "task_type" {
				whereClause += fmt.Sprintf(" %s = '%s'", key, value)
			}
			idx++
		}
		sqlQuery += whereClause
	}

	rows, err := db.Db.Query(sqlQuery)

	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	defer rows.Close()

	// Iterate over the rows and populate tasks slice
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.IsDone, &task.EventId, &task.TaskType)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		taskList = append(taskList, task)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to iterate over tasks"})
	}

	// Return tasklist as JSON response
	return c.Status(http.StatusOK).JSON(taskList)
}

func GetTask(c *fiber.Ctx) error {
	// get id parameters
	idParam := c.Params("id")

	query := `SELECT id, title, is_done, event_id, task_type FROM tasks WHERE id = $1`

	// Find task based on ID
	err := db.Db.QueryRow(query, idParam).Scan(&id, &title, &isDone, &eventId, &taskType)

	// Check if task exists
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).SendString(err.Error())
		}
		return c.Status(http.StatusNotFound).SendString(err.Error())
	}

	var task = models.Task{ID: id, Title: title, IsDone: isDone, EventId: eventId, TaskType: taskType}
	return c.Status(http.StatusOK).JSON(task)
}

func PostTask(c *fiber.Ctx) error {
	var taskData models.Task

	// Bind the Json body to Task struct
	// BodyParser is more lenient, doesn't check for validation error like missing required fields or invalid data types
	// BindJSON is stricter and will return 400 Bad request status and abort request.
	if err := c.BodyParser(&taskData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid request body: %s", err.Error())})

	}

	// Set default value for taskType if not provided in the request body
	if taskData.TaskType == "" {
		taskData.TaskType = "backlog"
	}

	query := `INSERT INTO tasks (title, is_done, event_id, task_type) VALUES ($1, $2, $3, $4) RETURNING id, title, is_done, event_id, task_type`

	err := db.Db.QueryRow(query, taskData.Title, taskData.IsDone, taskData.EventId, taskData.TaskType).Scan(&id, &title, &isDone, &eventId, &taskType)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Cannot find event"})

	}

	insertedTask := models.Task{ID: id, Title: title, IsDone: isDone, EventId: eventId, TaskType: taskType}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Task created successfully", "task": insertedTask})
}

func DeleteTask(c *fiber.Ctx) error {
	idParam := c.Params("id")

	query := `DELETE FROM tasks WHERE id = $1 RETURNING id, title, is_done, event_id, task_type`
	// Find task
	err := db.Db.QueryRow(query, idParam).Scan(&id, &title, &isDone, &eventId, &taskType)

	// Check if task exists
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	deletedTask := models.Task{ID: id, Title: title, IsDone: isDone, EventId: eventId, TaskType: taskType}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Task deleted successfully", "task": deletedTask})
}

func UpdateTask(c *fiber.Ctx) error {
	pkId := c.Params("id")

	// use map instead of Task model so that it will only map the fields defined in req.body
	var task map[string]interface{}

	if err := c.BodyParser(&task); err != nil {
		c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// query := "UPDATE tasks SET title = $1, is_done = $2, event_id = $3, task_type = $4 WHERE id = $5 RETURNING id, title, is_done, event_id, task_type"
	sqlQuery := "UPDATE tasks SET"
	args := []interface{}{}
	idx := 1

	// Loop through task fields and construct SET clauses
	for key, value := range task {
		if idx > 1 {
			sqlQuery += ","
		}
		sqlQuery += " " + key + " = $" + strconv.Itoa(idx)
		args = append(args, value)
		idx++
	}

	sqlQuery += " WHERE id = $" + strconv.Itoa(idx)
	sqlQuery += " RETURNING id, title, is_done, event_id, task_type"
	args = append(args, pkId)

	fmt.Println("sql query: ", sqlQuery)
	fmt.Println("args: ", args)

	err := db.Db.QueryRow(sqlQuery, args...).Scan(&id, &title, &isDone, &eventId, &taskType)

	// Find task based on ID
	// err := db.Db.QueryRow(query, task.Title, task.IsDone, task.EventId, task.TaskType, pkId).Scan(&id, &title, &isDone, &eventId, &taskType)

	// Check if task exists
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	updatedTask := models.Task{ID: id, Title: title, IsDone: isDone, EventId: eventId, TaskType: taskType}

	// Return updated task to the client
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Task updated successfully!", "task": updatedTask})
}
