package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bus struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Stopages map[int]string `json:"stopages"`
}

func (cfg *apiConfig) handlerGetBus(client *mongo.Client) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request)  {
		type parameters struct {
			Source string `json:"source"`
			Destination string `json:"destination"`
		}
	
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramaters")
			return
		}
	
		collection := client.Database("Buses").Collection("bus-info")
		source := strings.ToLower(params.Source)
		destination := strings.ToLower(params.Destination)

		filter := bson.M{
			"$and": []bson.M{
				{"stopages": bson.M{"$elemMatch": bson.M{"$eq": source}}},
				{"stopages": bson.M{"$elemMatch": bson.M{"$eq": destination}}},
			},
		}

		cursor, err := collection.Find(context.Background(), filter)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't find buses")
			return
		}
		defer cursor.Close(context.Background())

		var buses []string
		for cursor.Next(context.Background()) {
			var bus Bus
			if err := cursor.Decode(&bus); err != nil {
				log.Fatal(err)
				respondWithError(w, http.StatusNotFound, "Buses couldn't be fetched")
				return 
			}
			buses = append(buses, bus.Name)
		}
		if err := cursor.Err(); err != nil {
			log.Fatal(err)
			respondWithError(w, http.StatusInternalServerError, "Error while getting buses")
			return
		}
		respondWithJson(w, http.StatusFound, buses)
	}
}