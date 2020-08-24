package handler

import (
	"context"
	"curse-count/routers"
	"curse-count/websocket"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"time"
)

func serverWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		_, _ = fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func Handler(ctx context.Context) (err error) {
	router := mux.NewRouter()
	pool := websocket.NewPool()

	go pool.Start()

	router.HandleFunc("/", routers.Home).Methods("GET")
	router.HandleFunc("/new", routers.Person).Methods("POST")
	router.HandleFunc("/count/{id}", routers.Count).Methods("PUT")
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWs(pool, w, r)
	})

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	handler := cors.AllowAll().Handler(router)
	server := &http.Server{Addr: ":" + PORT, Handler: handler}

	go func() {
		if err = server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Server started")
	log.Println("Server listening on port " + PORT)
	<-ctx.Done()
	log.Println("Server stopping...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server graceful shutdown failed: %s\n", err)
	}

	log.Println("Server shut down gracefully")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
