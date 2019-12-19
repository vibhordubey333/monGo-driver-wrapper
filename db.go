// Provides helper functions that make interfacing with the MongoDB Go driver library easier
package db

import (
	"context"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create Context
var Ctx context.Context

// Create MongoDB client
var Client *mongo.Client

// Timeout limit for context before cancelling
const OperationTimeOut = 5

// Wrapper for Mongo Collection
type CnctConnection struct {
	Collection *mongo.Collection
}

// Function to create new connection to Mongo Collection
func New(uri, db, cnct string) (c CnctConnection) {
	log.Printf("Attempting to connect to %q", uri)
	conn, _ := context.WithTimeout(context.Background(), OperationTimeOut*time.Second)
	clt, err := mongo.Connect(Ctx, options.Client().ApplyURI(uri))

	if err != nil {
		log.Panic(err)
	}
	log.Print("Connection Established")

	// Change Package level vars
	Ctx = conn
	Client = clt

	c.Collection = Client.Database(db).Collection(cnct)
	return c
}

// Wrapper for collection.Drop()
func (db CnctConnection) Drop() error {
	// Set context
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)
	return db.Collection.Drop(Ctx)
}

// Wrapper for collection.FindOne(). Finds first document that satisfies filter and fills res with the unmarshalled document.
func (db CnctConnection) FindOne(filter bson.D, res interface{}) error {
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	err := db.Collection.FindOne(Ctx, filter).Decode(res)
	if err != nil {
		return err
	}
	return nil
}

// Wrapper for collection.Find(). Iterates cursor and fills res with unmarshalled documents.
func (db CnctConnection) FindMany(filter bson.D, res *[]interface{}) error {
	arrtype := reflect.TypeOf(res).Elem()

	// Set context
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	cursor, err := db.Collection.Find(Ctx, filter)

	for cursor.Next(Ctx) {
		doc := reflect.New(arrtype).Interface()
		err := cursor.Decode(&doc)
		if err != nil {
			return err
		}
		*res = append(*res, doc)
	}

	// unmarshall fail
	if cursor.Err() != nil {
		return err
	}

	// Close cursor after we're done with it
	cursor.Close(Ctx)
	return nil
}

// Wrapper for collection.UpdateOne(). Returns number of documents matched and modified. Should always be either 0 or 1.
func (db CnctConnection) UpdateOne(filter, update bson.D) (error, int64, int64) {
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	updateRes, err := db.Collection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return err, 0, 0
	}
	return nil, updateRes.MatchedCount, updateRes.ModifiedCount
}

// Wrapper for collection.UpdateMany(). Returns number of documents matched and modified.
func (db CnctConnection) UpdateMany(filter, update bson.D) (error, int64, int64) {
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	updateRes, err := db.Collection.UpdateMany(Ctx, filter, update)
	if err != nil {
		return err, 0, 0
	}
	return nil, updateRes.MatchedCount, updateRes.ModifiedCount
}

// Wrapper for collection.InsertOne(), doesn't return document and accepts arbitrary structs.
// Returns inserted ID
func (db CnctConnection) InsertOne(new interface{}) (error, interface{}) {
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	insertRes, err := db.Collection.InsertOne(Ctx, new)
	if err != nil {
		return err, ""
	}
	return nil, insertRes.InsertedID
}

// Wrapper for collection.InsertMany(), takes slice of structs to insert.
// Returns list of inserted IDs
func (db CnctConnection) InsertMany(new []interface{}) (error, interface{}) {
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	insertRes, err := db.Collection.InsertMany(Ctx, new)
	if err != nil {
		return err, ""
	}
	return nil, insertRes.InsertedIDs
}

// Wrapper for collection.DeleteOne(). Deletes single document that match the bson.D filter
func (db CnctConnection) DeleteOne(filter bson.D) error {
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)
	_, err := db.Collection.DeleteOne(Ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// Wrapper for collection.DeleteMany(). Deletes all documents that match the bson.D filter
func (db CnctConnection) DeleteMany(filter bson.D) error {
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)
	_, err := db.Collection.DeleteMany(Ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
