package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	database string
	username string
	password string
	w        http.ResponseWriter
)

func init() {
	database = "gorestapimongo"
	username = "gorestapimongo"
	password = "joAwxM4FaZuRDpRYXgGdXhEARn4TA9wWqa1xOLMir12xZy5DlxNkSL9QYi7Hek0ILquPWMbraA1uACDblj0vtg=="
}

type todo struct {
	ID        primitive.ObjectID `json: "id" bson: "_id,omitempty"`
	Item      string             `json: "item" bson: "item"`
	Completed bool               `json: "completed"  bson: "completed"`
}

var todos = []todo{}

func getTodos(contex *gin.Context) { //this function is to convert the array todos into JSON, cause in REST API client and server understand only JSON

	client := authenticateMongoDB()

	collection := client.Database(database).Collection("samplecollection")

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo todo
		err = cursor.Decode(&todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	contex.IndentedJSON(http.StatusOK, todos)
}

func addTodo(contex *gin.Context) {

	var newTodo todo

	if err := contex.BindJSON(&newTodo); err != nil { //BIND JSON is used to add the request from our body to the passed variable which is in this case type struct.
		return //this is used if the POST request to add the data is not in the format of the struct it will give an error
	}

	client := authenticateMongoDB()
	collection := client.Database(database).Collection("samplecollection")
	_, err := collection.InsertOne(context.Background(), &newTodo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todos = append(todos, newTodo)

	contex.IndentedJSON(http.StatusCreated, newTodo)

}

func getTodoID(id string) (*todo, error) { //this function is used to search the given ID in the array TODO and return struct todo or an error

	client := authenticateMongoDB()

	collection := client.Database(database).Collection("samplecollection")

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer cursor.Close(context.Background())
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	for cursor.Next(context.Background()) {
		var todo todo
		err = cursor.Decode(&todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if todo.ID == objID {
			return &todo, nil
		}
	}

	return nil, errors.New("ID NOT FOUND")
}

func toggleTodoStatus(contex *gin.Context) {

	client := authenticateMongoDB()

	collection := client.Database(database).Collection("samplecollection")
	id := contex.Param("id") //this is used to dynamically fetch the id from the http string
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	if err != nil {
		contex.IndentedJSON(http.StatusNotFound, gin.H{"message": "ID not found"}) //with this IF block we are checking the whether ID is there or not
		return
	}
	filter := bson.M{"id": bson.M{"$eq": objID}}

	update := bson.M{"$set": bson.M{"item": "New Update Request"}}

	result, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		panic(err)
	}

	contex.IndentedJSON(http.StatusOK, result)

}

func getTodo(context *gin.Context) {

	id := context.Param("id")  //this is used to dynamically fetch the id from the http string
	todo, err := getTodoID(id) //getting the todo and err status from GETID func

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "ID not found"}) //with this IF block we are checking the whether ID is there or not
		return
	}
	context.IndentedJSON(http.StatusOK, todo)
}

func authenticateMongoDB() *mongo.Client {

	connecturi := fmt.Sprintf(
		"mongodb://%s:%s@%s.documents.azure.com:10255/?ssl=true",
		username,
		password,
		database)

	// Set the client options
	clientOptions := options.Client().ApplyURI(connecturi)

	// Set the context with a 10-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the MongoDB instance
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Cosmos DB MongoDB instance!")

	return client

}

func main() {
	router := gin.Default()                      //to create the server
	router.GET("/todos", getTodos)               //this is the method for GET request
	router.GET("/todos/:id", getTodo)            //this is to call getTodo function which is searching for dynamic ID in the array todo
	router.PATCH("/todos/:id", toggleTodoStatus) //this is to call the PATCH http request, It is changing the completed boolean.
	router.POST("/todos", addTodo)               //this is the method for POST request
	router.Run("localhost:9090")                 //to run the server on port 9090

}
