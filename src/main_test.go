package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setUp() *mux.Router {
	defineBooks()
	router := getRouter()
	return router
}

func TestRetrieveBook(t *testing.T) {
	router := setUp()
	assert := assert.New(t)
	var book Book
	request, err := http.NewRequest("GET", "/book/1", nil)
	assert.Nil(err)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)
	assert.Equal(200, responseRecorder.Code, "The books exists, should get 200")
	json.Unmarshal([]byte(responseRecorder.Body.String()), &book)
	assert.Equal(book.Id, "1")
}

func TestRetrieveAllBooks(t *testing.T){
	router := setUp()
	assert := assert.New(t)
	request, err := http.NewRequest("GET", "/books", nil)
	assert.Nil(err)
	responseRecorder := httptest.NewRecorder()
	var books []map[string]string
	router.ServeHTTP(responseRecorder, request)
	assert.Equal(200, responseRecorder.Code)
	json.Unmarshal([]byte(responseRecorder.Body.String()), &books)
	assert.Equal(2, len(books))
}


func TestCreateBook(t *testing.T) {
	router := setUp()
	assert := assert.New(t)
	data := map[string]string{"Author": "Joaco", "Title": "Twilight", "PublicationYear": "2005"}
    json_data, err := json.Marshal(data)
	var book Book
	request, err := http.NewRequest("POST", "/book", bytes.NewBuffer(json_data))
	assert.Nil(err)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)
	assert.Equal(201, responseRecorder.Code, "An object should be created")
	json.Unmarshal([]byte(responseRecorder.Body.String()), &book)
	assert.Equal("3", book.Id)
	assert.Equal("Joaco", book.Author)
	assert.Equal("Twilight", book.Title)
	assert.Equal("2005", book.PublicationYear)
}

 func TestUpdateBook(t *testing.T){
	router := setUp()
	assert := assert.New(t)
	var book Book
	data := map[string]string{"Author": "Ale", "Title": "Una casa en la pradera", "PublicationYear": "2012"}
	json_data, err := json.Marshal(data)
	request, err := http.NewRequest("PUT", "/book/2", bytes.NewBuffer(json_data))
	assert.Nil(err)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)
	json.Unmarshal([]byte(responseRecorder.Body.String()), &book)
	assert.Equal(200, responseRecorder.Code)
	assert.Equal("2", book.Id)
	assert.Equal("Ale", book.Author)
	assert.Equal("Una casa en la pradera", book.Title)
	assert.Equal("2012", book.PublicationYear)
}

func TestDeleteBook(t *testing.T) {
	router := setUp()
	assert := assert.New(t)
	request, err := http.NewRequest("DELETE", "/book/2", nil)
	assert.Nil(err)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)
	assert.Equal(204, responseRecorder.Code)
}
