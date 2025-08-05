package main

import (
	"net/http"
	"os"
	"planning-poker/application/planningsvc"
	"planning-poker/delivery/websocket"
	"planning-poker/infra/in_memory"
)

func main() {
	planningRepo := in_memory.NewPlanningRepository()
	planningSvc := planningsvc.NewPlanningService(planningRepo)
	wsHandler := websocket.NewWebsocketHandler(planningSvc)

	http.Handle("/ws", wsHandler)
	http.Handle("/", http.FileServer(http.Dir("./frontend/")))
	http.HandleFunc("/session/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
