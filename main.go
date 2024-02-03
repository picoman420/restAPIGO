package main

type todo struct{
	id 		string 	`json: "id"`
	item 	string		`json: "item"`
	completed 	bool	`json: "completed"`
}

var todos = []todo{
	{id:"1", item:"Clean Room", completed:false}
	{id:"2", item:"Read Book", completed:false}
	{id:"3", item:"Record Video", completed:false}
}