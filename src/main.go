package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/mux"
	"encoding/json"
)

type Book struct {
	Id string `json:"Id"`
	Author string `json:"Author"`
	Title string `json:"Title"`
	PublicationYear string `json:"PublicationYear"`
}


type Library struct {
	Books []Book `json:"Books"`
}

var BOOKS []Book

func getBook(writer http.ResponseWriter, request *http.Request) {
	routeVariables := mux.Vars(request)
	bookId := routeVariables["id"]
	found := false
	for _, book := range BOOKS {
		if book.Id == bookId {
			json.NewEncoder(writer).Encode(book)
			found = true
		}
	}
	if !found {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func createNewBook(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return 
	}
	var book Book
	json.Unmarshal(requestBody, &book)
	book.Id = fmt.Sprint(len(BOOKS) + 1)
	BOOKS = append(BOOKS, book)
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(book)
}

func updateBook(writer http.ResponseWriter, request *http.Request) {
	routeVariables := mux.Vars(request)
	bookId := routeVariables["id"]
	requestBody, _ := ioutil.ReadAll(request.Body)
	for _, book := range BOOKS {
		if book.Id == bookId {
			json.Unmarshal(requestBody, &book)
			book.Id = bookId
			json.NewEncoder(writer).Encode(book)
		}
	}
	writer.WriteHeader(http.StatusNotFound)
}

func deleteBook(writer http.ResponseWriter, request *http.Request) {
	routeVariables := mux.Vars(request)
	bookId := routeVariables["id"]
	for index, book := range BOOKS {
		if book.Id == bookId {
			BOOKS = append(BOOKS[:index], BOOKS[index+1:]...)
			writer.WriteHeader(http.StatusNoContent)
			return
		}
	}
	writer.WriteHeader(http.StatusNotFound)

}

func retrieveAllBooks(writer http.ResponseWriter, request *http.Request) {
	json.NewEncoder(writer).Encode(BOOKS)
}


func getRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/books", retrieveAllBooks).Methods("GET")
	router.HandleFunc("/book",  createNewBook).Methods("POST")
	router.HandleFunc("/book/{id}", getBook).Methods("GET")
	router.HandleFunc("/book/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/book/{id}", deleteBook).Methods("DELETE")
	return router
}


func HandleRequest(port string){
	router := getRouter()
	log.Fatal(http.ListenAndServe(port, router))
}

func defineBooks(){
	BOOKS = []Book{
		Book{Id: "1", Author: "JM", Title: "Aprendiendo a programar en GO", PublicationYear: "2021"},
		Book{Id: "2", Author: "Ale", Title: "Una casa en la pradera", PublicationYear: "1930"},
	}
}

func main() {
	port := "8000"
	defineBooks()
	fmt.Println("Starting server on port", port)
	HandleRequest(port)
}
