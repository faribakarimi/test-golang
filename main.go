package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/faribakarimi/test-golang/api/auth"
	"github.com/faribakarimi/test-golang/api/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
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

var connector *gorm.DB

func connect(connectionString string) error {
	var err error
	connector, err = gorm.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	log.Println("Connection was successful!")
	return nil
}

func migrateUser(table *models.User) {
	connector.AutoMigrate(&table)
	log.Println("Table Users migrated.")
}

func register(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user models.User
	json.Unmarshal(reqBody, &user)
	connector.Create(&user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user.ID)
}

func login(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	user := models.User{}
	json.Unmarshal(reqBody, &user)
	var err error
	user1 := models.User{}
	err = connector.Debug().Model(models.User{}).Where("username = ?", user.Username).Take(&user1).Error
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	err = models.VerifyPassword(user1.Password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	token, err := auth.CreateToken(user1.ID)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	type Response struct {
		UserId int
		Token string
	}
	json.NewEncoder(w).Encode(Response{user1.ID, token})
}

func profile(w http.ResponseWriter, r *http.Request)      {}
func editProfile(w http.ResponseWriter, r *http.Request)  {}
func buy(w http.ResponseWriter, r *http.Request)          {}
func returnItems(w http.ResponseWriter, r *http.Request)  {}
func balance(w http.ResponseWriter, r *http.Request)      {}
func getUserItems(w http.ResponseWriter, r *http.Request) {}

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

func main() {

	fmt.Println("Test Golang - Rest API")

	config := Config{
		ServerName: "localhost:3306",
		User:       "root",
		Password:   "root",
		DB:         "test_golang_api",
	}

	connectionString := getConnectionString(config)
	err := connect(connectionString)
	if err != nil {
		panic(err.Error())
	}

	migrateUser(&models.User{})

	// var item models.Item
	// item.ID = 1
	// item.Name = "item number 1"
	// item.Price = 500
	// connector.Create(item)

	handleRequests()
}
