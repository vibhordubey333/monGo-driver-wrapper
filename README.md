[![GoDoc](https://godoc.org/github.com/jackyzha0/monGo-driver-wrapper?status.svg)](https://godoc.org/github.com/jackyzha0/monGo-driver-wrapper)
[![GoReportCard](https://goreportcard.com/badge/github.com/jackyzha0/monGo-driver-wrapper)](https://goreportcard.com/report/github.com/jackyzha0/monGo-driver-wrapper)
# MongoDB Go Driver Wrapper
This package reduces boilerplate when using MongoDB by handling contexts and cursor iteration for you.

## Install
Use the package by running `go get -u github.com/jackyzha0/monGo-driver-wrapper/`

## Testing / Development
Ensure that you have a `mongod` instance running locally on your machine. Run the sanity check tests by doing `go test -v`. This will create a new Database called `exampleDB` and collection named `test` which it will delete after.
