package crud_via_gorm

import (
	"crud/internal/grpc/pb"
	"testing"
)

func TestCrudViaGorm(t *testing.T) {
	databaseConnection := Connect()

	t.Run("Create task and send to database", func(t *testing.T) {

		newtask := pb.Task{Name: "A task just for testing", Completed: true}

		rowsAffected, err := databaseConnection.CreateTask(&newtask)

		assertError(t, err, nil)

		assertRowsAffected(t, rowsAffected, 1)
	})

	t.Run("Retrieve all tasks from database", func(t *testing.T) {

		_, err := databaseConnection.ListAll()

		assertError(t, err, nil)
	})

	t.Run("Retrieve task by id", func(t *testing.T) {
		task, err := databaseConnection.GetTaskById(6)

		assertError(t, err, nil)

		if task == nil {

			t.Errorf("Got nil, but wanted a task")
		}
	})

	t.Run("Retrieve tasks by completion", func(t *testing.T) {

		completion := true
		tasks, err := databaseConnection.GetTaskByCompletion(completion)

		areAllTasksAccordingTo := VerifyCompletion(completion, tasks)

		assertError(t, err, nil)

		if areAllTasksAccordingTo != true {
			t.Errorf("Got %v Want %v", areAllTasksAccordingTo, completion)
		}

	})

	t.Run("Delete task by id", func(t *testing.T) {

		rowsAffected, err := databaseConnection.DeleteTask(4)

		assertError(t, err, nil)
		assertRowsAffected(t, rowsAffected, 1)
	})

	t.Run("Update task by id", func(t *testing.T) {
		newtask := pb.Task{Name: "Task for update purpose only", Completed: false}
		rowsAffected, err := databaseConnection.UpdateTask(8, &newtask)

		assertError(t, err, nil)
		assertRowsAffected(t, rowsAffected, 1)
	})

}

func VerifyCompletion(completion bool, tasks []*pb.Task) bool {
	for _, task := range tasks {

		if task.Completed != completion {
			return false
		}
	}

	return true
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("Got %v Want %v", got, want)
	}
}

func assertRowsAffected(t testing.TB, got, want int64) {
	t.Helper()

	if got != want {

		t.Errorf("Got: %d Want: %d", got, want)
	}
}
