package main

import (
	"net/http"
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

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
