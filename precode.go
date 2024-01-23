package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Обработчик для получения всех задач
func getTastks(w http.ResponseWriter, r *http.Request) {

	// сериализуем данные из мапы tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// записываем в ответ тип контента в заголовок
	w.Header().Set("Content-Type", "application/json")

	// записываем в ответ статус 200 OK
	w.WriteHeader(http.StatusOK)

	// записываем в ответ сам контент, а именно сериализованый json объект
	w.Write(resp)
}

// обработчик для добавления задачи
// пример для постмана:
// {"id": "3", "description": "Выпить чай", "note": "таёжный сбор, 2 ч. ложки, 5 мин.", "applications":["кружка", "заварной чайник"]}
func postTastks(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// обработчик для выгрузки конкретной задачи с заданым id
func getTask(w http.ResponseWriter, r *http.Request) {

	// забираем параметр из урла
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusNoContent)
		return
	}

	// сериализуем полученные данные
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// записываем в ответ тип контента в заголовок
	w.Header().Set("Content-Type", "application/json")

	// записываем в ответ статус 200 OK
	w.WriteHeader(http.StatusOK)

	// записываем в ответ сам контент, а именно сериализованый json объект
	w.Write(resp)
}

func main() {
	r := chi.NewRouter()

	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getTastks`
	// http://localhost:8080/tasks
	r.Get("/tasks", getTastks)

	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `postTastks`
	// http://localhost:8080/tasks
	r.Post("/tasks", postTastks)

	// регистрируем в роутере эндпоинт `/task/{id}` с методом GET, для которого используется обработчик `getTask`
	// http://localhost:8080/task/2
	r.Get("/task/{id}", getTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
