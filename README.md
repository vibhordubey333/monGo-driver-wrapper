[![GoDoc](https://godoc.org/github.com/jackyzha0/monGo-driver-wrapper?status.svg)](https://godoc.org/github.com/jackyzha0/monGo-driver-wrapper)
[![GoReportCard](https://goreportcard.com/badge/github.com/jackyzha0/monGo-driver-wrapper)](https://goreportcard.com/report/github.com/jackyzha0/monGo-driver-wrapper)
# MongoDB Go Driver Wrapper
This package reduces boilerplate when using MongoDB by handling contexts and cursor iteration for you.


## Install
Use the package by running `go get -u github.com/jackyzha0/monGo-driver-wrapper/` and adding the import to the top of your file

```go
db "github.com/jackyzha0/monGo-driver-wrapper"
```

## Usage
#### Instantiate Connection
Begin by creation a connection to a MongoDB collection by specifying the URL, database name, and collection name.

```go
// Create new connection to the `users` collection in the `exampleDB` database on the local instance
var Collection = db.New("mongodb://localhost:27017", "exampleDB", "users")
```

#### Insert a document
```go
// First, create a struct to represent your data
type Doc struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

// Create a new document
john := Doc{"john", "smith"}

// Insert said document into the collection, and get the insertID assigned by MongoDB
err, insertID = Collection.InsertOne(john)
```
#### Insert multiple documents
Assuming you've already designed a struct to represent your data and made a connection,
```go
// Define some structs
john := Doc{"john", "smith"}
betty := Doc{"betty", "hansen"}

// Make a slice out of them
sl := []interface{}{john, betty}
err, insertIDs := Collection.InsertMany(sl)
```

#### Find a document
```go
// Define a filter to search for. Below is a query to search for all documents with the name `bob`
filter := bson.D{{"name", "bob"}}

// Define the struct that you wish to unmarshall the response into. If the call is successful, res will take on the value of the result.
var res Doc
err := Collection.FindOne(filter, &res)
```
#### Find multiple documents
```go
filter := bson.D{{"surname", "joe"}}

// Make sure that when searching for multiple documents, you use a slice of interfaces
var res []interface{}
err := Collection.FindMany(filter, &res)

// Then, create a new slice of the desired type of struct
var val []Doc

// This part unmarshalls the slice of interfaces into the desired slice of structs.
for _, el := range got {
  var d Doc
  bsonBytes, _ := bson.Marshal(el)
  bson.Unmarshal(bsonBytes, &d)
  val = append(val, d)
}
```

#### Update a document
```go
// Define a filter to select which types of documents should be updated
update_filter := bson.D{{"name", "rebecca"}}

// Define how to update said documents
update := bson.D{{"$set", bson.D{{"surname", "o'connor"}}}}

// Returns number of documents matched and modified.
err, match, mod := Collection.UpdateOne(update_filter, update)
```
#### Update multiple documents
```go
// Increment the points of anyone on team red by one
update_filter := bson.D{{"team", "red"}}
update := bson.D{{"$inc", bson.D{{"points", 1}}}}
err, match, mod := Collection.UpdateMany(update_filter, update)
```

#### Delete a document
```go
// Delete a single document with `name: albert`
filter := bson.D{{"name", "albert"}}
err = TestCollection.DeleteOne(filter)
```
#### Delete multiple documents
```go
// Delete all documents with `name: albert`
filter := bson.D{{"name", "albert"}}
err = TestCollection.DeleteOne(filter)
```

## Testing / Development
Ensure that you have a `mongod` instance running locally on your machine. Run the sanity check tests by doing

```go test -v```

This will create a new Database called `exampleDB` and collection named `test` which it will delete after.
