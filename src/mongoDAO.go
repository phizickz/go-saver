package main

import "go.mongodb.org/mongo-driver/mongo"

serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
clientOptions := options.Client().
    ApplyURI("mongodb+srv://<username>:<password>@<host-url>/?retryWrites=true&w=majority").
    SetServerAPIOptions(serverAPIOptions)
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
client, err := mongo.Connect(ctx, clientOptions)
if err != nil {
    log.Fatal(err)
}