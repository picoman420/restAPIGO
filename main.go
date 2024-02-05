package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
)

var (
	database string
	username string
	password string
)

func init() {
	database = "gorestapimongodb"
	username = "gorestapimongodb"
	password = "ec4ZJJOn1Ig8D9W9rhXaI6LAchfEzzEAorzkbK3JMp842q7D01Ml44aHnPiyalBNA3pPeRVJiygCACDbimWopg=="
}

type todo struct {
	ID        string `json: "id"`
	Item      string `json: "item"`
	Completed bool   `json: "completed"`
}

var todos = []todo{
	{ID: "1", Item: "Clean Room", Completed: false},
	{ID: "2", Item: "Read Book", Completed: false},
	{ID: "3", Item: "Record Video", Completed: false},
}

func getTodos(context *gin.Context) { //this function is to convert the array todos into JSON, cause in REST API client and server understand only JSON

	context.IndentedJSON(http.StatusOK, todos)
}

func addTodo(context *gin.Context) {

	var newTodo todo
	if err := context.BindJSON(&newTodo); err != nil { //BIND JSON is used to add the request from our body to the passed variable which is in this case type struct.
		return //this is used if the POST request to add the data is not in the format of the struct it will give an error
	}

	todos = append(todos, newTodo)

	context.IndentedJSON(http.StatusCreated, newTodo)

}

func getTodoID(id string) (*todo, error) { //this function is used to search the given ID in the array TODO and return struct todo or an error

	for i, t := range todos { //this is used to iterate the todos array and find the given ID, if ID is not here then user defined error is thrown using errors package
		if t.ID == id {
			return &todos[i], nil
		}
	}
	return nil, errors.New("ID NOT FOUND")
}

func toggleTodoStatus(context *gin.Context) {
	id := context.Param("id")  //this is used to dynamically fetch the id from the http string
	todo, err := getTodoID(id) //getting the todo and err status from GETID func

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "ID not found"}) //with this IF block we are checking the whether ID is there or not
		return
	}

	todo.Completed = !todo.Completed

	context.IndentedJSON(http.StatusOK, todo)

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

func authenticateMongoDB() {
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
}

func main() {
	authenticateMongoDB()
	router := gin.Default()                      //to create the server
	router.GET("/todos", getTodos)               //this is the method for GET request
	router.GET("/todos/:id", getTodo)            //this is to call getTodo function which is searching for dynamic ID in the array todo
	router.PATCH("/todos/:id", toggleTodoStatus) //this is to call the PATCH http request, It is changing the completed boolean.
	router.POST("/todos", addTodo)               //this is the method for POST request
	router.Run("localhost:9090")                 //to run the server on port 9090

}
