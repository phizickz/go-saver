package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	path := "/health"

	if r.URL.Path != path  {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}


	fmt.Fprintf(w, "Healthy!")
}

func fetchDocument(client mongo.Client, database string, collection string, documentID int) {
	coll := client.Database(database).Collection(collection)

	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"id", documentID}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the ID %d\n", documentID)
		return
	}
	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
	// return jsonData
}

func mongoDBClient() *mongo.Client {
	mdbConnStr := os.Getenv("MONGODB.CONNECTIONSTRING")
	if mdbConnStr == "" {
		log.Fatal("MongoDB connection string not set. Verify .env file or environment variables.")
	}
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
    clientOptions := options.Client().
        ApplyURI(mdbConnStr).
        SetServerAPIOptions(serverAPIOptions)
    ctx, cancel:= context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        panic(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

func main() {
	http.HandleFunc("/health", healthHandler)

	var mdbClient *mongo.Client = mongoDBClient()

	mdbName := os.Getenv("MONGODB.DATABASE")
	if mdbName == "" {
		log.Fatal("MongoDB database name set. Verify .env file or environment variables.")
	}

	mdbCollection := "test"
	fetchDocument(*mdbClient, mdbName, mdbCollection, 123)
	// fmt.Printf("%s\n", )

	fmt.Printf("Starting server at port 8080.\n")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}
}