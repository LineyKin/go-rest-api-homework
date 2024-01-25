package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
	// w.Write возвращает 2 ответа: кол-во байтов и ошибку
	// тут и других местах, где вызывается метод, проверим на ошибку
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}

// обработчик для добавления задачи
func postTasks(w http.ResponseWriter, r *http.Request) {
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

	// если элемент c id = task.ID уже есть, вернём ошибку
	_, ok := tasks[task.ID]
	if ok {
		http.Error(w, "Элемент с таким id уже существует", http.StatusBadRequest)
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
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// сериализуем полученные данные
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// записываем в ответ тип контента в заголовок
	w.Header().Set("Content-Type", "application/json")

	// записываем в ответ статус 200 OK
	w.WriteHeader(http.StatusOK)

	// записываем в ответ сам контент, а именно сериализованый json объект
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	// забираем параметр из урла
	id := chi.URLParam(r, "id")

	// нас интересует лишь само существование задачи
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// удаляем задачу из мапы
	delete(tasks, id)

	// сериализуем оставшиеся данные
	_, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// записываем в ответ тип контента в заголовок
	w.Header().Set("Content-Type", "application/json")

	// записываем в ответ статус 200 OK
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getTastks`
	// http://localhost:8080/tasks
	r.Get("/tasks", getTastks)

	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `postTastks`
	// http://localhost:8080/tasks
	// пример для постмана:
	// {"id": "3", "description": "Выпить чай", "note": "таёжный сбор, 2 ч. ложки, 5 мин.", "applications":["кружка", "заварной чайник"]}
	r.Post("/tasks", postTasks)

	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом GET, для которого используется обработчик `getTask`
	// http://localhost:8080/tasks/2
	r.Get("/tasks/{id}", getTask)

	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом DELETE, для которого используется обработчик `deleteTask`
	// http://localhost:8080/tasks/3
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
