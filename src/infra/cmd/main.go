package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/robertokbr/events-worker/src/domain/dtos"
	"github.com/robertokbr/events-worker/src/infra/controllers"
	"github.com/robertokbr/events-worker/src/infra/database"
	"github.com/robertokbr/events-worker/src/infra/database/repositories"
	"github.com/robertokbr/events-worker/src/infra/middlewares"
	"github.com/robertokbr/events-worker/src/infra/providers"
	"github.com/robertokbr/events-worker/src/usecases"
)

func init() {
	godotenv.Load()
}

func main() {
	mux := http.NewServeMux()
	connection := database.NewMysqlDB()
	repository := repositories.NewMySqlRepository(connection)
	mailProvider := &providers.SesMailProvider{}
	jobs := make(chan *dtos.EventDTO, 100)
	eventsWorkerManager := usecases.NewEventsWorkerManager(jobs, repository, mailProvider)
	eventsController := controllers.NewEventsController(jobs)

	go eventsWorkerManager.Start(7)

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			middlewares.RequestAuth(w, r, eventsController.Add)
			break
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Not found"}`))
		}
	})

	fmt.Printf("Server running on port %s\n", ":8080")
	http.ListenAndServe(":8080", mux)
}
