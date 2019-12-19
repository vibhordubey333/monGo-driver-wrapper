// Package db provides helper functions that make interfacing with the MongoDB Go driver library easier
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

// Ctx holds the current context
var Ctx context.Context

// Client holds the reference to the underlying MongoDB client
var Client *mongo.Client

// OperationTimeOut is the time limit for context before cancelling
const OperationTimeOut = 5

// CnctConnection is the wrapper for Mongo Collection
type CnctConnection struct {
	Collection *mongo.Collection
}

// New creates a new connection to Mongo Collection
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

// Drop drops the current CnctConnection (collection)
func (db CnctConnection) Drop() error {
	// Set context
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)
	return db.Collection.Drop(Ctx)
}

// FindOne finds first document that satisfies filter and fills res with the unmarshalled document.
func (db CnctConnection) FindOne(filter bson.D, res interface{}) error {
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	err := db.Collection.FindOne(Ctx, filter).Decode(res)
	if err != nil {
		return err
	}
	return nil
}

// FindMany iterates cursor of all docs matching filter and fills res with unmarshalled documents.
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

// UpdateOne updates single document matching filter and applies update to it. Returns number of documents matched and modified. Should always be either 0 or 1.
func (db CnctConnection) UpdateOne(filter, update bson.D) (int64, int64, error) {
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	updateRes, err := db.Collection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return 0, 0, err
	}
	return updateRes.MatchedCount, updateRes.ModifiedCount, nil
}

// UpdateMany updates all documents matching the filter by applying the update query on it. Returns number of documents matched and modified.
func (db CnctConnection) UpdateMany(filter, update bson.D) (int64, int64, error) {
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	updateRes, err := db.Collection.UpdateMany(Ctx, filter, update)
	if err != nil {
		return 0, 0, err
	}
	return updateRes.MatchedCount, updateRes.ModifiedCount, nil
}

// InsertOne inserts a single struct as a document into the database and returns its ID.
// Returns inserted ID
func (db CnctConnection) InsertOne(new interface{}) (interface{}, error) {
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	insertRes, err := db.Collection.InsertOne(Ctx, new)
	if err != nil {
		return "", err
	}
	return insertRes.InsertedID, nil
}

// InsertMany takes a slice of structs, inserts them into the database, and returns list of inserted IDs
func (db CnctConnection) InsertMany(new []interface{}) (interface{}, error) {
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)

	insertRes, err := db.Collection.InsertMany(Ctx, new)
	if err != nil {
		return "", err
	}
	return insertRes.InsertedIDs, nil
}

// DeleteOne deletes single document that match the bson.D filter
func (db CnctConnection) DeleteOne(filter bson.D) error {
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)
	_, err := db.Collection.DeleteOne(Ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMany deletes all documents that match the bson.D filter
func (db CnctConnection) DeleteMany(filter bson.D) error {
	Ctx, _ = context.WithTimeout(context.Background(), OperationTimeOut*time.Second)
	_, err := db.Collection.DeleteMany(Ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
