package main

import (
	"context"
	"log"
	"net/http"

	"github.com/PratikforCoding/search-bus.git/internal/database"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileServerHits int
}

func main() {

	client, err := db.ConnectToMongo()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	router := chi.NewRouter()
	root := "."
	apicfg := apiConfig {
		fileServerHits: 0,
	}
	
	fsHandler := apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(root))))
	corsMux := middlewareCors(router)

	router.Handle("/app/", fsHandler)
	router.Handle("/app/*", fsHandler)
	router.Get("/app/getbus", apicfg.handlerGetBus(client))

	server := &http.Server{
		Addr : ":8080",
		Handler : corsMux,
	}

	

	log.Println("Server running at port: 8080...")
	log.Fatal(server.ListenAndServe())
}