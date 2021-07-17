package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/register", register)
	myRouter.HandleFunc("/login", login)
	myRouter.HandleFunc("/article", profile)
	myRouter.HandleFunc("/profile", editProfile).Methods("PUT")
	myRouter.HandleFunc("/items", returnItems)
	myRouter.HandleFunc("buy", buy).Methods("POST")
	myRouter.HandleFunc("balance", balance)
	myRouter.HandleFunc("userItems", getUserItems)
	
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func register(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the register!")
}
func login(w http.ResponseWriter, r *http.Request) {}
func profile(w http.ResponseWriter, r *http.Request){}
func editProfile(w http.ResponseWriter, r *http.Request){}
func buy(w http.ResponseWriter, r *http.Request){}
func returnItems(w http.ResponseWriter, r *http.Request){}
func balance(w http.ResponseWriter, r *http.Request){}
func getUserItems(w http.ResponseWriter, r *http.Request){}

func main() {
	fmt.Println("Test Golang - Rest API")
	handleRequests()
}
