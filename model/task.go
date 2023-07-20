package model

import (
	"database/sql"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Items       []string  `json:"items"`
}

type TaskDto struct {
	DB *sql.DB
}

func (taskDto TaskDto) GetAllTasks() ([]Task, error) {
	// Prepare the SQL query to fetch all tasks and their associated items.
	query := `
		SELECT t.id, t.title, t.description, t.completed, t.created_at, t.updated_at, ti.item
		FROM task t
		LEFT JOIN task_item ti ON t.id = ti.task_id
		ORDER BY t.id, ti.id
	`

	// Execute the query and retrieve the rows.
	rows, err := taskDto.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a map to store the tasks based on their ID.
	taskMap := make(map[int]*Task)

	// Create a map to store the task items grouped by task ID.
	taskItemsMap := make(map[int][]string)

	// Loop through the rows and scan each task and its item into the maps.
	for rows.Next() {
		var taskItem sql.NullString
		task := &Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
			&taskItem,
		)
		if err != nil {
			return nil, err
		}

		// Add the task item to the corresponding task items map entry.
		if taskItem.Valid {
			taskItemsMap[task.ID] = append(taskItemsMap[task.ID], taskItem.String)
		}

		// Check if the task ID already exists in the task map.
		// If not, add the task to the map.
		if _, ok := taskMap[task.ID]; !ok {
			taskMap[task.ID] = task
		}
	}

	// Combine the task items with the corresponding tasks.
	for taskID, taskItems := range taskItemsMap {
		if task, ok := taskMap[taskID]; ok {
			task.Items = taskItems
		}
	}

	// Convert the map of tasks into a slice and return it.
	tasks := make([]Task, 0, len(taskMap))
	for _, task := range taskMap {
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func (taskDto TaskDto) Insert(task *Task) error {

	// Prepare the SQL statement for the task table.
	stmt, err := taskDto.DB.Prepare(`
			INSERT INTO task (title, description, completed, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement and get the inserted task ID.
	var taskID int
	err = stmt.QueryRow(task.Title, task.Description, task.Completed, task.CreatedAt, task.UpdatedAt).Scan(&taskID)
	if err != nil {
		return err
	}

	// Set the ID field of the task struct with the generated task ID.
	task.ID = taskID

	return nil

}

func (taskDto TaskDto) Get(id int64) (*Task, error) {
	return nil, nil
}

func (taskDto TaskDto) UpdateTask(id int, task *Task) error {
	// Prepare the SQL statement to update the task by ID.
	stmt, err := taskDto.DB.Prepare(`
		UPDATE task
		SET title = $1, description = $2, completed = $3, updated_at = $4
		WHERE id = $5
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement to update the task.
	_, err = stmt.Exec(task.Title, task.Description, task.Completed, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

func (taskDto TaskDto) DeleteTask(id int) error {
	// Prepare the SQL statement to delete the task by ID.
	stmt, err := taskDto.DB.Prepare(`
		DELETE FROM task WHERE id = $1
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the task.
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (taskDto TaskDto) InsertTaskItem(taskID int, item string) error {
	// Prepare the SQL statement for the task_item table.
	stmt, err := taskDto.DB.Prepare(`
		INSERT INTO task_item (task_id, item)
		VALUES ($1, $2)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement to insert the task item.
	_, err = stmt.Exec(taskID, item)
	if err != nil {
		return err
	}

	return nil
}

func (taskDto TaskDto) GetTask(id int) (*Task, error) {
	// Prepare the SQL query to fetch the task by ID.
	query := `
		SELECT t.id, t.title, t.description, t.completed, t.created_at, t.updated_at, ti.item
		FROM task t
		LEFT JOIN task_item ti ON t.id = ti.task_id
		WHERE t.id = $1
	`

	// Execute the query and retrieve the row.
	row := taskDto.DB.QueryRow(query, id)

	// Scan the row into a new task object.
	task := &Task{}
	var taskItem sql.NullString
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
		&taskItem,
	)
	if err != nil {
		return nil, err
	}

	// If there are task items, append them to the TaskItems field of the task.
	if taskItem.Valid {
		task.Items = append(task.Items, taskItem.String)
	}

	return task, nil
}
