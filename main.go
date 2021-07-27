package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/faribakarimi/test-golang/api/auth"
	"github.com/faribakarimi/test-golang/api/database"
	"github.com/faribakarimi/test-golang/api/models"
	"github.com/faribakarimi/test-golang/api/seed"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

func register(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user models.User
	json.Unmarshal(reqBody, &user)
	database.Connector.Create(&user)
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
	err = database.Connector.Debug().Model(models.User{}).Where("username = ?", user.Username).Take(&user1).Error
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
		Token  string
	}
	json.NewEncoder(w).Encode(Response{user1.ID, token})
}

func profile(w http.ResponseWriter, r *http.Request) {
	err := auth.TokenValid(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	var user models.User
	database.Connector.First(&user, uid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Username string
		Name string
		Family string
		Gender string
		Age uint16
	}{
		user.Username,
		user.Name,
		user.Family,
		user.Gender,
		user.Age,
	})
}

func editProfile(w http.ResponseWriter, r *http.Request) {
	err := auth.TokenValid(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var user models.User
	database.Connector.Where("id = ?", uid).Find(&user)
	json.Unmarshal(reqBody, &user)
	result := database.Connector.Save(&user)

	w.Header().Set("Content-Type", "application/json")
	if result.RowsAffected == 1 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(true)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(false)
	}
}
func returnItems(w http.ResponseWriter, r *http.Request)  {
	var items []models.Item
	database.Connector.Find(&items)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func buy(w http.ResponseWriter, r *http.Request) {
	err := auth.TokenValid(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	var user models.User
	database.Connector.Where("id = ?", uid).Find(&user)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("---")
	}
	var item models.Item
	err = json.Unmarshal(reqBody, &item)
	if err != nil {
		fmt.Println("...")
		return
	}
	err = database.Connector.Where("id = ?", item.ID).Find(&item).Error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(struct{
			Error string `json:"string"`
		}{
			Error: "Item Not Found",
		})
		return
	}
	// TODO:: Duplicte item_id
	if user.Balance <= uint64(item.Price) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(false)
		return
	}
	err = database.Connector.Create(&models.UserItems{UserID: uid,ItemID: item.ID,}).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

func balance(w http.ResponseWriter, r *http.Request)      {}
func getUserItems(w http.ResponseWriter, r *http.Request) {}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/register", register)
	myRouter.HandleFunc("/login", login)
	myRouter.HandleFunc("/profile", profile).Methods("GET")
	myRouter.HandleFunc("/profile", editProfile).Methods("PUT")
	myRouter.HandleFunc("/items", returnItems)
	myRouter.HandleFunc("/buy", buy).Methods("POST")
	myRouter.HandleFunc("balance", balance)
	myRouter.HandleFunc("userItems", getUserItems)

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func initDB() {
	config := database.Config{
		ServerName: "172.26.0.2:3306",
		User:       "root",
		Password:   "root",
		DB:         "test_golang_api",
	}
	connectionString := database.GetConnectionString(config)
	err := database.Connect(connectionString)
	if err != nil {
		panic(err.Error())
	}
	database.MigrateUser(&models.User{})
	database.MigrateItem(&models.Item{})
	database.MigrateUserItems(&models.UserItems{})
}

func main() {
	initDB()
	fmt.Println("Test Golang - Rest API")
	seed.Load(database.Connector)
	handleRequests()
}
