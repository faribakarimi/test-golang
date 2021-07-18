package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Config struct {
	ServerName string
	User       string
	Password   string
	DB         string
}

var getConnectionString = func(config Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true", config.User, config.Password, config.ServerName, config.DB)
}

var Connector *gorm.DB

func connect(connectionString string) error {
	var err error
	Connector, err = gorm.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	log.Println("Connection was successful!")
	return nil
}

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

	config := Config{
		ServerName: "localhost:3306",
		User:       "root",
		Password:   "root",
		DB:         "test-golang",
	}

	connectionString := getConnectionString(config)
	err := connect(connectionString)
	if err != nil {
		panic(err.Error())
	}

	handleRequests()
}
