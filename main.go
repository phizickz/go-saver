package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"math/rand"
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
type Address struct {
	Street string
	City   string
}
type Student struct {
	// FirstName string  `bson:"first_name,omitempty"`
	// LastName  string  `bson:"last_name,omitempty"`
	// Address   Address `bson:"inline"`
	FirstName string 
	LastName  string  
	Address   Address 
	Age       int
}
func fetchSingleDocument(coll mongo.Collection, documentID int) {
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

func insertSingleStudent(coll mongo.Collection, data Student) {
	_, err := coll.InsertOne(
		context.TODO(),
		data)
	if err != nil {
		panic(err)
	}
}

func generateDummyStudent() Student {
	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source) // initialize local pseudorandom generator

	addressPrefix := [...]string{"Skiolds","Stor","Alfabet","Trombone"}
	addressSuffix := [...]string{"vegen","veien","gata","gaten"}
	cities := [...]string{"Lillestrøm","Oslo","Drammen"}

	firstNames := [...]string{"Johan","Geir","Tomas","Filip"}
	lastNames := [...]string{"Pedersen","Hansson","Ultron"}
	ages := [...]int{20,25,35,30}

	tempAddress := Address{Street: addressPrefix[random.Intn(len(addressPrefix))] + addressSuffix[random.Intn(len(addressSuffix))], City: cities[random.Intn(len(cities))]}
	return Student{FirstName: firstNames[random.Intn(len(firstNames))], LastName: lastNames[random.Intn(len(lastNames))], Address: tempAddress, Age: ages[random.Intn(len(ages))]}
}

func main() {
	http.HandleFunc("/health", healthHandler)

	var mdbClient *mongo.Client = mongoDBClient()
	mdbName := os.Getenv("MONGODB.DATABASE")
	if mdbName == "" {
		log.Fatal("MongoDB database name set. Verify .env file or environment variables.")
	}
	mdbCollection := os.Getenv("MONGODB.COLLECTION")
	if mdbName == "" {
		log.Fatal("MongoDB collection name set. Verify .env file or environment variables.")
	}
	coll := mdbClient.Database(mdbName).Collection(mdbCollection)

	var students [5]Student

	for i := 0; i >= len(students); i++ {
		students[i] = generateDummyStudent()
	}

	insertSingleStudent(*coll, students[0])

	// fetchSingleDocument(*coll, 123)

	fmt.Printf("Starting server at port 8080.\n")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}
}