package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Entrenador struct {
	Nombre string
	Edad   int
	Ciudad string
}

func main() {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conectado a BD de Mongo")

	// Get a handle for your collection
	collection := client.Database("goTest").Collection("users")

	// Some dummy data to add to the Database
	ash := Entrenador{"Ash", 10, "Pueblo Paleta"}
	misty := Entrenador{"Misty", 10, "Ciudad Celeste"}
	brock := Entrenador{"Brock", 15, "Ciudad Plateada"}

	// Insert a single document
	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Insertando solo un documento: ", insertResult.InsertedID)

	// Insert multiple documents
	trainers := []interface{}{misty, brock}

	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Insertando multiples documentos: ", insertManyResult.InsertedIDs)

	// Update a document
	filter := bson.D{{"nombre", "Ash"}}

	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Se encontró %v documento y se actualizó %v documento.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Find a single document
	var result Entrenador

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Se encontró 1 documento: %+v\n", result)

	findOptions := options.Find()
	findOptions.SetLimit(2)

	var results []*Entrenador

	// Finding multiple documents returns a cursor
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the cursor
	for cur.Next(context.TODO()) {
		var elem Entrenador
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Se encontraron multiples documentos: %+v\n", results)

	// Delete all the documents in the collection
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Se eliminaron %v documentos en la coleccion de entrenadores\n", deleteResult.DeletedCount)

	// Close the connection once no longer needed
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Cerrando conexion a MongoDB.")
	}

}
