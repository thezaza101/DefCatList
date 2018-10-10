package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	. "./config"
	. "./dao"
	"gopkg.in/mgo.v2/bson"

	. "./models"

	"github.com/gorilla/mux"
)

var config = Config{}
var dao = DefListsDAO{}

//Get all public lists
func AllListsEndPoint(w http.ResponseWriter, r *http.Request) {
	lists, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, lists)
}

//Find a list by ID
func FindListEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	lists, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid List ID")
		return
	}
	respondWithJson(w, http.StatusOK, lists)
}

//Create a list
func CreateListEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var list DefList
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	list.ID = bson.NewObjectId()
	if err := dao.Insert(list); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, list)
}

//Update a list
func UpdateListEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var list DefList
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(list); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

//Delete a list
func DeleteListEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var list DefList
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(list); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

//Respond with an error msg
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

//respond with a json string
func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Set the config info and connect the the database
func init() {
	config.Read()
	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/lists", AllListsEndPoint).Methods("GET")
	r.HandleFunc("/lists", CreateListEndPoint).Methods("POST")
	r.HandleFunc("/lists", UpdateListEndPoint).Methods("PUT")
	r.HandleFunc("/lists", DeleteListEndPoint).Methods("DELETE")
	r.HandleFunc("/lists/{id}", FindListEndpoint).Methods("GET")
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		log.Fatal(err)
	}
}
